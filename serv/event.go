package serv

import (
	"github.com/kovey/cli-go/app"
)

type EventInterface interface {
	OnFlag(app.AppInterface) error
	OnBefore(app.AppInterface) error
	OnAfter(app.AppInterface) error
	OnRun() error
	OnShutdown()
	CreateConfig(app.AppInterface) error
	Usage() bool
}

type EventBase struct {
}

func (s *EventBase) OnBefore(app.AppInterface) error {
	return nil
}

func (s *EventBase) OnAfter(app.AppInterface) error {
	return nil
}

func (s *EventBase) OnRun() error {
	return nil
}

func (s *EventBase) OnShutdown() {
}

func (s *EventBase) CreateConfig(app.AppInterface) error {
	return nil
}

func (s *EventBase) Usage() bool {
	return false
}
