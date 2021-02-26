package scanners

import (
	"context"
	"fmt"
	"oasisTracker/conf"
	"oasisTracker/dao"
	"sync"
	"time"
)

type (
	Manager struct {
		cfg   conf.Config
		dao   dao.DAO
		tasks map[uint64]bool

		wg   *sync.WaitGroup
		ctx  context.Context
		stop context.CancelFunc
	}
)

func NewManager(cfg conf.Config, d dao.DAO) *Manager {
	ctx, stop := context.WithCancel(context.Background())

	return &Manager{
		cfg:   cfg,
		dao:   d,
		tasks: make(map[uint64]bool),

		wg:   &sync.WaitGroup{},
		ctx:  ctx,
		stop: stop,
	}
}

func (m *Manager) Title() string {
	return "Scanners Manager"
}

func (m *Manager) Stop() error {
	m.stop()
	m.wg.Wait()
	return nil
}

func (m *Manager) Run() error {

	for {
		select {
		case <-m.ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			//Get active tasks
			tasks, err := m.dao.GetTasks(true)
			if err != nil {
				return fmt.Errorf("dao.GetTasks: %s", err.Error())
			}

			for _, task := range tasks {
				if !task.IsActive {
					continue
				}

				//Already run
				if ok := m.tasks[task.ID]; ok {
					continue
				}

				m.wg.Add(1)
				m.tasks[task.ID] = true

				scanner, err := NewScanner(m.cfg.Scanner, task, m.dao, m.ctx)
				if err != nil {
					return fmt.Errorf("NewScanner (%s): %s", task.Title, err.Error())
				}

				go func() {
					scanner.Run()
					m.wg.Done()
				}()
			}
		}
	}

	m.wg.Wait()
	return nil
}
