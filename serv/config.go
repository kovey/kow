package serv

import (
	"fmt"
	"os"

	"github.com/kovey/discovery/etcd"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Etcd   etcd.Config `yaml:"etcd"`
	Listen Listen      `yaml:"listen"`
	App    App         `yaml:"app"`
}

type App struct {
	TimeZone  string `yaml:"time_zone"`
	PprofOpen string `yaml:"pprof_open"`
	EtcdOpen  string `yaml:"etcd_open"`
}

type Listen struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (l *Listen) Addr() string {
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}

func (c *Config) Load(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, c)
}
