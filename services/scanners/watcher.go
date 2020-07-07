package scanners

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"oasisTracker/common/log"
	"oasisTracker/conf"
	"oasisTracker/dao"
	"oasisTracker/dmodels"
)

type Watcher struct {
	dao        dao.DAO
	cfg        conf.Config
	parser     *Parser
	ctx        context.Context
	cancelFunc context.CancelFunc

	ReSyncInit bool
}

func NewWatcher(cfg conf.Config, d dao.DAO) (*Watcher, error) {
	ctx, cancel := context.WithCancel(context.Background())

	dao, err := d.GetParserDAO()
	if err != nil {
		return nil, fmt.Errorf("GetParserDAO: %s", err.Error())
	}

	parser, err := NewParser(ctx, cfg.Scanner, dao)
	if err != nil {
		return nil, fmt.Errorf("NewParser: %s", err.Error())
	}

	return &Watcher{
		cfg:        cfg,
		dao:        d,
		parser:     parser,
		ctx:        ctx,
		cancelFunc: cancel,
	}, nil
}

func (m *Watcher) Title() string {
	return "Watcher module"
}

func (m *Watcher) Stop() error {
	m.cancelFunc()
	return nil
}

func (m *Watcher) Run() error {
	ch, cPub, err := m.parser.api.WatchBlocks(m.ctx)
	if err != nil {
		return fmt.Errorf("WatchBlocks error: %s", err)
	}

	for {
		select {
		case <-m.ctx.Done():
			cPub.Close()
			return nil
		case block := <-ch:
			if !m.ReSyncInit {
				err = m.addReSyncTask(block.Height)
				if err != nil {
					log.Error("AddReSyncTask error", zap.Error(err))
					continue
				}
				m.ReSyncInit = true
			}

			err = m.parser.ParseWatchBlock(block)
			if err != nil {
				log.Error("ParseBlock error", zap.Error(err))
				continue
			}

			err = m.parser.Save()
			if err != nil {
				log.Error("Save error", zap.Error(err))
				continue
			}

		}
	}
}

func (m *Watcher) addReSyncTask(currentHeight int64) error {
	task, isFound, err := m.dao.GetLastTask()
	if err != nil {
		return fmt.Errorf("GetLastTask error: %s", err)
	}

	startHeight := task.EndHeight + 1
	if !isFound {
		startHeight = 0
	}

	//Blocks sync
	err = m.dao.CreateTask(dmodels.Task{
		IsActive:      true,
		Title:         parserBaseTask,
		StartHeight:   startHeight,
		CurrentHeight: startHeight,
		EndHeight:     uint64(currentHeight),
		Batch:         m.cfg.Scanner.BatchSize,
	})
	if err != nil {
		return fmt.Errorf("CreateTask error: %s", err)
	}

	//Snaps sync
	err = m.dao.CreateTask(dmodels.Task{
		IsActive:      true,
		Title:         parserBalancesSnapshotTask,
		StartHeight:   startHeight,
		CurrentHeight: startHeight,
		EndHeight:     uint64(currentHeight),
		//1 Epoch = 600 blocks
		Batch: 20000,
	})
	if err != nil {
		return fmt.Errorf("CreateTask error: %s", err)
	}

	//Return when refactor workers
	//for startHeight <= uint64(currentHeight-1) {
	//	endHeight := startHeight + m.cfg.Scanner.NodeRPS
	//	if endHeight > uint64(currentHeight-1) {
	//		endHeight = uint64(currentHeight - 1)
	//	}
	//
	//	err = m.dao.CreateTask(dmodels.Task{
	//		IsActive:      true,
	//		Title:         parserSignaturesTask,
	//		StartHeight:   startHeight,
	//		CurrentHeight: startHeight,
	//		EndHeight:     endHeight,
	//		Batch:         200,
	//	})
	//	if err != nil {
	//		return fmt.Errorf("CreateTask error: %s", err)
	//	}
	//
	//	startHeight += m.cfg.Scanner.NodeRPS + 1
	//}

	return nil
}
