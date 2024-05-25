package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	MongoDBAddr           string        `env:"MongoAddr" envDefault:"mongodb://127.0.0.1:27017"`
	MysqlAddr             string        `env:"MysqlAddr" envDefault:"127.0.0.1:3306"`
	MysqlName             string        `env:"MysqlName" envDefault:"root"`
	MysqlPassword         string        `env:"MysqlPasswd" envDefault:"123456"`
	EtcdAddr              []string      `env:"EtcdAddr" envDefault:"127.0.0.1:2379" envSeparator:","`
	KafkaAddr             []string      `env:"KafkaAddr" envDefault:"127.0.0.1:9092" envSeparator:","`
	Environment           string        `env:"Environment" envDefault:"Production"`
	FillErrorInfoFlag     bool          `env:"FillErrorInfoFlag" envDefault:"true"`
	DomainIp              string        `env:"DomainIp" envDefault:"127.0.0.1"`
	TrackDBAddr           string        `env:"TraceAddr" envDefault:"127.0.0.1:6831"`
	RegisterServicePrefix string        `env:"RegisterServicePrefix" envDefault:"farm"`
	GameRedisAddr         string        `env:"GameRedisAddr" envDefault:"127.0.0.1:6020"`
	UserRedisAddr         string        `env:"UserRedisAddr" envDefault:"127.0.0.1:6021"`
	BackstageRedisAddr    string        `env:"BackstageRedisAddr" envDefault:"127.0.0.1:6022"`
	TableRedisAddr        string        `env:"TableRedisAddr" envDefault:"127.0.0.1:6023"`
	ApiAddr               string        `env:"ApiAddr" envDefault:"127.0.0.1:8050/farm"`
	PackageName           string        `env:"PackageName" envDefault:"com.lt.fs.google"`
	LanguageId            uint32        `env:"LanguageId" envDefault:"102"`
	PVPConsole            bool          `env:"PVPConsole" envDefault:"false"`
	Reentry               uint32        `env:"Reentry" envDefault:"3"`
	RegisterInterval      time.Duration `env:"RegisterInterval" envDefault:"10S"`
	RegisterTTL           time.Duration `env:"RegisterTTL" envDefault:"30S"`
	DDUrl                 string        `env:"DDUrl" envDefault:""`
	DDSecret              string        `env:"DDSecret" envDefault:""`
	ServerName            string        `env:"ServerName"`
	AccessPort            string        `env:"AccessPort" envDefault:":7001"`
	ApiPort               string        `env:"ApiPort" envDefault:":8050"`
	ChatPort              string        `env:"ChatPort" envDefault:":7011"`
	BackstagePort         string        `env:"BackstagePort" envDefault:":8060"`
	CHAddr                []string      `env:"CHAddr" envDefault:"127.0.0.1:9000" envSeparator:","`
	CHDB                  string        `env:"CHDB" envDefault:"farmstory"`
	CHUser                string        `env:"CHUser" envDefault:"default"`
	CHPassword            string        `env:"CHPassword"`
}

var C Config

func init() {
	if err := env.Parse(&C); err != nil {
		panic(err)
	}
}

var (
	SessionName string
	exeName     string

	Session string
	Index   int
)

func GetSessionName() string {
	return SessionName
}

func init() {
	exeName = filepath.Base(os.Args[0])
	ext := filepath.Ext(exeName)
	if len(ext) > 0 {
		exeName = strings.TrimSuffix(exeName, ext)
	}

	name, ok := os.LookupEnv("SESSION_NAME")
	if ok {
		// name = strings.TrimPrefix(name, "svc-")
		SessionName = name
	} else {
		SessionName = fmt.Sprintf("%s-%d", exeName, os.Getpid())
	}
	Session, Index = SplitSession(SessionName)
	initBuildInfo()
}

func GetEnv(name string, def string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return def
}

func ExeName() string {
	return exeName
}

func SplitSession(s string) (string, int) {
	i := strings.LastIndex(s, "-")
	if i == -1 {
		return "", 0
	}
	prefix := s[:i]
	suffix := s[i+1:]
	si, err := strconv.ParseInt(suffix, 10, 32)
	if err != nil {
		return "", 0
	}
	return prefix, int(si)
}

const developSerEnvironment = "Develop"

func IsDevelopSerEnvironment() bool {
	return C.Environment == developSerEnvironment
}
