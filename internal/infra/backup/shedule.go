package backup

import "context"

type Backuper interface {
	Store(ctx context.Context) error
}

type Schedule struct {
	Backend Backuper
}

func NewSchedule(b Backuper) *Schedule {
	return &Schedule{
		Backend: b,
	}
}
