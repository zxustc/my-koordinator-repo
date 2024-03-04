package frontend

import (
	"context"
	"time"
)

var _ StatusFetcher = &statusFetcher{}

type statusFetcher struct {
	ctx context.Context
}

type StatusFetcher interface {
	Run()
	Started() bool
	UpdateStatus()
}

func InitStatusFetcher(ctx context.Context) *statusFetcher {
	return &statusFetcher{
		ctx: ctx,
	}
}

func getNextUpdate(now time.Time) time.Time {
	next := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	if next.Before(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}

func (s *statusFetcher) Run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(time.Until(getNextUpdate(time.Now()))):
			s.UpdateStatus()
		}
	}
}

func (s *statusFetcher) Started() bool {
	return false
}

func (s *statusFetcher) UpdateStatus() {

}

//used to update crd status
