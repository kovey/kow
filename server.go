package kow

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	_ "net/http/pprof"

	"github.com/kovey/cli-go/app"
	"github.com/kovey/cli-go/env"
	"github.com/kovey/debug-go/debug"
	"github.com/kovey/discovery/etcd"
	"github.com/kovey/kow/resolver"
	"github.com/kovey/kow/serv"
)

var engine = NewDefault()

type server struct {
	*app.ServBase
	wait  sync.WaitGroup
	e     serv.EventInterface
	pprof *http.Server
}

func newServer(e serv.EventInterface) *server {
	return &server{wait: sync.WaitGroup{}, e: e}
}

func (s *server) Flag(a app.AppInterface) error {
	if s.e != nil {
		return s.e.OnFlag(a)
	}

	return nil
}

func (s *server) Init(a app.AppInterface) error {
	location, err := time.LoadLocation(os.Getenv(APP_TIME_ZONE))
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
	if ppOpen, err := env.GetBool(APP_PPROF_OPEN); err != nil || !ppOpen {
		return
	}
	port, err := env.GetInt(SERV_PORT)
	if err != nil {
		return
	}

	s.pprof = &http.Server{Addr: fmt.Sprintf("%s:%d", os.Getenv(SERV_HOST), port+10000), Handler: http.DefaultServeMux}
	if err := s.pprof.ListenAndServe(); err != nil {
		debug.Erro("run pprof failure, error: %s", err)
	}
}

func (s *server) runOhter() {
	defer s.wait.Done()
	if s.e == nil {
		return
	}

	if err := s.e.OnRun(); err != nil {
		debug.Erro("event.OnRun failure, error: %s", err)
	}
}

func (s *server) Run(a app.AppInterface) error {
	if etcdOpen, err := env.GetBool(APP_ETCD_OPEN); err == nil && etcdOpen {
		timeout, _ := env.GetInt(ETCD_TIMEOUT)
		conf := etcd.Config{
			Endpoints:   strings.Split(os.Getenv(ETCD_ENDPOINTS), ","),
			DialTimeout: timeout,
			Username:    os.Getenv(ETCD_USERNAME),
			Password:    os.Getenv(ETCD_PASSWORD),
			Namespace:   os.Getenv(ETCD_NAMESPACE),
		}

		if err := resolver.Register(conf); err != nil {
			return err
		}
	}

	s.wait.Add(1)
	go s.runMonitor()
	s.wait.Add(1)
	go s.runOhter()

	debug.Info("app[%s] listen on [%s:%s]", a.Name(), os.Getenv(SERV_HOST), os.Getenv(SERV_PORT))
	if err := engine.Run(fmt.Sprintf("%s:%s", os.Getenv(SERV_HOST), os.Getenv(SERV_PORT))); err != nil {
		return err
	}
	return nil
}

func (s *server) Shutdown(a app.AppInterface) error {
	if s.pprof != nil {
		if err := s.pprof.Shutdown(context.Background()); err != nil {
			debug.Erro("shutdown pprof failure, error: %s", err)
		}
	}

	if err := engine.Shutdown(); err != nil {
		debug.Erro("engine shutdown failure, error: %s", err)
	}

	if s.e != nil {
		s.e.OnShutdown()
	}

	if etcdOpen, err := env.GetBool(APP_ETCD_OPEN); err == nil && etcdOpen {
		resolver.Shutdown()
	}
	s.wait.Wait()
	return nil
}
