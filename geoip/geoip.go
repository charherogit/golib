package geoip

import (
	"context"
	"fmt"
	"github.com/maxmind/geoipupdate/v6/pkg/geoipupdate"
	"github.com/oschwald/geoip2-golang"
	"github.com/robfig/cron/v3"
	log "golib/logutil"
	"os"
	"sync"
)

type IpSearch struct {
	reader     *geoip2.Reader
	path       string
	cronConfig *CronConfig

	lock sync.RWMutex
}

type CronConfig struct {
	CronSpec    string // 定时器规则(建议为"0 23 * * 2,5")
	ConfigFile  string // 配置文件路径
	DatabaseDir string // 数据库文件路径 XXX 应配置为path所在目录
}

const defaultDatabaseFile = "./config/GeoLite2-City.mmdb"

func NewIpSearch(options ...Option) (*IpSearch, error) {
	i := &IpSearch{
		lock: sync.RWMutex{},
		path: defaultDatabaseFile,
	}
	for _, option := range options {
		if err := option(i); err != nil {
			return nil, fmt.Errorf("set option err: %s", err)
		}
	}
	return i, nil
}

func (i *IpSearch) load() error {
	newReader, err := geoip2.Open(i.path)
	if err != nil {
		return err
	}

	i.lock.Lock()
	defer i.lock.Unlock()
	i.reader = newReader
	return nil
}

func (i *IpSearch) ReLoad() error {
	return i.load()
}

func (i *IpSearch) Close() {
	i.reader.Close()
}

type Option func(f *IpSearch) error

func WithDataBasePath(path string) Option {
	return func(i *IpSearch) error {
		if path == "" {
			return fmt.Errorf("path is empty")
		}

		i.path = path
		return i.load()
	}
}

func WithCronConfig(cron *CronConfig) Option {
	return func(i *IpSearch) error {
		if cron == nil {
			return fmt.Errorf("cronConfig is nil")
		}

		// 检测配置文件是否存在
		if _, err := os.Stat(cron.ConfigFile); err != nil {
			return fmt.Errorf("cron config file err: %s", err)
		}

		i.cronConfig = cron
		return createReaderTimer(i)
	}
}

func createReaderTimer(ipSearch *IpSearch) error {
	if ipSearch.cronConfig == nil {
		return fmt.Errorf("cronConfig is nil")
	}

	cron := cron.New()
	if _, err := cron.AddFunc(ipSearch.cronConfig.CronSpec, func() {
		if err := updateDatabaseFile(ipSearch.cronConfig); err != nil {
			log.Errorf("update geoip2 database err: %s", err)
			return
		}

		if err := ipSearch.ReLoad(); err != nil {
			log.Errorf("geoip2 ReLoad err: %s", err)
			return
		}
		log.Debugf("update geoip2 database success")
	}); err != nil {
		return err
	}
	cron.Start()
	return nil
}

// XXX 定时更新只支持固定文件名称GeoLite2-City.mmdb、GeoLite2-Country.mmdb、GeoLite2-ASN.mmdb
func updateDatabaseFile(cfg *CronConfig) error {
	config, err := geoipupdate.NewConfig(
		geoipupdate.WithConfigFile(cfg.ConfigFile),
		geoipupdate.WithDatabaseDirectory(cfg.DatabaseDir),
		geoipupdate.WithVerbose(false),
		geoipupdate.WithOutput(false),
	)
	if err != nil {
		return err
	}

	client := geoipupdate.NewClient(config)
	if err = client.Run(context.Background()); err != nil {
		return err
	}
	return nil
}
