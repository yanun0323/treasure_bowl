package util

import (
	"context"

	"github.com/robfig/cron/v3"
)

type CronJobManager interface {
	AddFunc(spec string, fn func())
	Clean()
	ResetFunc(spec string, fn func())
	Start()
	Stop() []context.Context
}

type cronJobManager struct {
	fns   map[string][]func()
	crons map[string]*cron.Cron
}

func NewCronJobManager(capacity ...int) CronJobManager {
	cap := 0
	if len(capacity) != 0 {
		cap = capacity[0]
	}
	return &cronJobManager{
		fns:   make(map[string][]func(), cap),
		crons: make(map[string]*cron.Cron, cap),
	}
}

func (cm *cronJobManager) AddFunc(spec string, fn func()) {
	cm.fns[spec] = append(cm.fns[spec], fn)
}

func (cm *cronJobManager) ResetFunc(spec string, fn func()) {
	delete(cm.fns, spec)
}

func (cm *cronJobManager) Clean() {
	clear(cm.fns)
	cm.Stop()
}

func (cm *cronJobManager) Start() {
	for k, fns := range cm.fns {
		c := cm.requireCron(k)
		c.AddFunc(k, func() {
			for _, fn := range fns {
				fn()
			}
		})
		c.Start()
	}

}

func (cm *cronJobManager) Stop() []context.Context {
	ctxs := make([]context.Context, 0, len(cm.crons))
	for _, c := range cm.crons {
		ctxs = append(ctxs, c.Stop())
	}
	clear(cm.crons)
	return ctxs
}

func (cm *cronJobManager) requireCron(key string) *cron.Cron {
	if cm.crons[key] == nil {
		cm.crons[key] = cron.New()
	}
	return cm.crons[key]
}
