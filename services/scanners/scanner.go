package scanners

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/google"
	"oasisTracker/common/log"
	"oasisTracker/conf"
	"oasisTracker/dao"
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
	"time"

	"go.uber.org/zap"
)

const repeatPause = time.Second * 5

type (
	Scanner struct {
		cfg         conf.Scanner
		task        dmodels.Task
		executor    *smodels.Executor
		dao         dao.DAO
		ctx         context.Context
		stopFunc    context.CancelFunc
		credentials credentials.TransportCredentials
		blocksCh    chan uint64
		resultCh    chan error
	}
	ExecutorProvider interface {
		GetTaskExecutor(taskTitle string) (executor *smodels.Executor, err error)
	}
)

func NewScanner(cfg conf.Scanner, task dmodels.Task, d dao.DAO, ctx context.Context) (s *Scanner, err error) {
	scCtx, stop := context.WithCancel(ctx)

	s = &Scanner{
		cfg:         cfg,
		task:        task,
		dao:         d,
		ctx:         scCtx,
		stopFunc:    stop,
		credentials: google.NewDefaultCredentials().TransportCredentials(),
		blocksCh:    make(chan uint64, task.Batch),
		resultCh:    make(chan error, task.Batch),
	}

	var p ExecutorProvider
	dao, err := d.GetParserDAO()
	if err != nil {
		return nil, fmt.Errorf("GetParserDAO: %s", err.Error())
	}

	p, err = NewParser(ctx, cfg, dao)
	if err != nil {
		return nil, fmt.Errorf("Create NewParser: %s", err.Error())
	}

	s.executor, err = p.GetTaskExecutor(task.Title)
	if err != nil {
		return nil, fmt.Errorf("p.GetTaskExecutor: %s", err.Error())
	}

	return s, nil
}

func (s *Scanner) Run() {
	s.runWorkers()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			//Truncate previous executor after each round
			s.executor.Truncate()
		}

		log.Debug("Start scanner ", zap.String("task_name", s.task.Title), zap.Uint64("task_current", s.task.CurrentHeight))
		var err error
		lastHeight := s.task.EndHeight

		if lastHeight <= s.task.CurrentHeight {
			s.stopFunc()
			continue
		}

		batch := s.task.Batch
		if lastHeight-s.task.CurrentHeight < s.task.Batch {
			batch = lastHeight - s.task.CurrentHeight
		}

		currentHeight := s.task.CurrentHeight

		tm := time.Now()
		for i := currentHeight; i < currentHeight+batch; i++ {
			s.blocksCh <- i
		}

		isFail := false
		for i := currentHeight; i < currentHeight+batch; i++ {
			err = <-s.resultCh
			if err != nil {
				log.Error("Scanner Result", zap.Error(err), zap.String("task", s.task.Title))
				isFail = true
			}
		}

		log.Debug("Execution time", zap.String("task_name", s.task.Title), zap.Int64("Milliseconds", time.Since(tm).Milliseconds()))

		if isFail {
			<-time.After(repeatPause)
			continue
		}

		tm = time.Now()
		err = s.executor.Save()
		if err != nil {
			log.Error("Scanner Save", zap.Error(err), zap.String("task", s.task.Title), zap.Uint64("task_height", s.task.CurrentHeight))
			<-time.After(repeatPause)
			continue
		}

		log.Debug("Save time", zap.String("task_name", s.task.Title), zap.Int64("Milliseconds", time.Since(tm).Milliseconds()))

		s.task.CurrentHeight += batch
		if s.task.CurrentHeight == s.task.EndHeight {
			s.task.IsActive = false
		}

		for {
			err = s.dao.UpdateTask(s.task)
			if err == nil {
				break
			}
			log.Error("Scanner UpdateTask", zap.Error(err), zap.String("task", s.task.Title))
			<-time.After(repeatPause)
			continue
		}
	}
}

func (s *Scanner) runWorkers() {
	for i := uint64(0); i < s.cfg.NodeRPS; i++ {
		go func() {

			//grpcConn, err := grpc.Dial(s.cfg.NodeConfig, grpcCommon.WithTransportCredentials(s.credentials))
			//if err != nil {
			//	log.Error("grpc.Dial", zap.Error(err))
			//	return
			//}
			//defer grpcConn.Close()

			for {
				select {
				case <-s.ctx.Done():
					return
				case blockID := <-s.blocksCh:
					err := s.executor.ExecHeight(blockID)
					if err != nil {
						if s.task.Title == parserBaseTask {
							err = fmt.Errorf("block %d : %s", blockID, err.Error())
						} else {
							err = fmt.Errorf("epoch %d : %s", blockID, err.Error())
						}
					}

					s.resultCh <- err
				}
			}
		}()
	}
}
