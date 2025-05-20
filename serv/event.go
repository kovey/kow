package serv

import "github.com/kovey/cli-go/app"

type EventInterface interface {
	OnFlag(app.AppInterface) error
	OnBefore(app.AppInterface) error
	OnAfter(app.AppInterface) error
	OnRun() error
	OnShutdown()
	CreateConfig(app.AppInterface) error
	Usage() error
}
