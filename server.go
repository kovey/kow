package kow

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	_ "net/http/pprof"

	"github.com/kovey/cli-go/app"
	"github.com/kovey/debug-go/debug"
	"github.com/kovey/kow/resolver"
	"github.com/kovey/kow/serv"
)

var engine = NewDefault()

type server struct {
	*app.ServBase
	conf  *serv.Config
	wait  sync.WaitGroup
	e     serv.EventInterface
	pprof *http.Server
}

func newServer(e serv.EventInterface) *server {
	return &server{wait: sync.WaitGroup{}, e: e}
}

func (s *server) loadConf(a app.AppInterface) error {
	path, err := a.Get("c")
	if err != nil {
		return err
	}

	tmp := path.String()
	if tmp == "" {
		return fmt.Errorf("path is empty")
	}

	conf := &serv.Config{}
	if err := conf.Load(tmp); err != nil {
		return err
	}

	s.conf = conf
	return nil
}

func (s *server) Flag(a app.AppInterface) error {
	a.Flag("c", "", app.TYPE_STRING, "app config file path")
	if s.e != nil {
		return s.e.OnFlag(a)
	}

	return nil
}

func (s *server) Init(a app.AppInterface) error {
	if err := s.loadConf(a); err != nil {
		return err
	}

	location, err := time.LoadLocation(s.conf.TimeZone)
	if err != nil {
		return err
	}
	time.Local = location

	if s.e != nil {
		return s.e.OnBefore(a)
	}

	if s.e != nil {
		return s.e.OnAfter(a)
	}

	return nil
}

func (s *server) runMonitor() {
	defer s.wait.Done()
	s.pprof = &http.Server{Addr: fmt.Sprintf("%s:%d", s.conf.Listen.Host, s.conf.Listen.Port+10000), Handler: http.DefaultServeMux}
	if err := s.pprof.ListenAndServe(); err != nil {
		debug.Erro("run pprof failure, error: %s", err)
	}
}

func (s *server) runOhter() {
	defer s.wait.Done()
	if s.e == nil {
		return
	}

	s.e.OnRun()
}

func (s *server) Run(a app.AppInterface) error {
	if err := resolver.Register(s.conf.Etcd); err != nil {
		return err
	}

	s.wait.Add(1)
	go s.runMonitor()
	s.wait.Add(1)
	go s.runOhter()

	debug.Info("app[%s] listen on [%s]", a.Name(), s.conf.Listen.Addr())
	if err := engine.Run(s.conf.Listen.Addr()); err != nil {
		return err
	}
	return nil
}

func (s *server) Shutdown(a app.AppInterface) error {
	if s.pprof != nil {
		s.pprof.Shutdown(context.Background())
	}

	engine.Shutdown()

	if s.e != nil {
		s.e.OnShutdown()
	}

	s.wait.Wait()
	return nil
}
