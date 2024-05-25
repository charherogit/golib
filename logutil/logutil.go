package logutil

import (
	"github.com/caarlos0/env/v8"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	mlogger "go-micro.dev/v4/logger"
	"golib/config"
	"io"
	"os"
	"path/filepath"
)

type LogConfig struct {
	Level       logrus.Level `env:"LogLevel" envDefault:"trace"`
	MaxSize     int          `env:"LogMaxSize" envDefault:"50"`
	MaxAge      int          `env:"LogMaxAge" envDefault:"90"`
	MaxBackups  int          `env:"LogMaxBackups" envDefault:"500"`
	LogFilePath string       `env:"LogFilePath"`
}

var logger = &logrus.Logger{
	Out:       os.Stdout,
	Formatter: logFormatter,
	Level:     logrus.TraceLevel,
	ExitFunc:  os.Exit,
} // 默认的日志对象

var lc LogConfig

var logFormatter = &logrus.TextFormatter{
	FullTimestamp:   true,
	DisableColors:   true,
	TimestampFormat: "2006-01-02 15:04:05.000",
	SortingFunc:     sortOutput,
}

func newLogger(fn string) *logrus.Logger {
	return &logrus.Logger{
		Out: &lumberjack.Logger{
			Filename:   fn,            // 日志文件路径
			MaxSize:    lc.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
			MaxAge:     lc.MaxAge,     // 文件最多保存多少天
			MaxBackups: lc.MaxBackups, // 日志文件最多保存多少个备份
			LocalTime:  true,          // 文件名本地时间
			Compress:   false,         // 是否压缩
		},
		Hooks:     make(logrus.LevelHooks),
		Formatter: logFormatter,
		Level:     lc.Level,
		ExitFunc:  os.Exit,
	}
}

// 日志初始化(使用手动初始化的话会严重影响依赖加载顺序，还是先用init初始化吧)
func init() {
	if err := env.Parse(&lc); err != nil {
		panic(err)
	}
	if lc.LogFilePath != "" {
		logger = newLogger(filepath.Join(lc.LogFilePath, config.GetSessionName()+".log"))
	} else {
		logger = newLogger(filepath.Join("logs", config.GetSessionName()+".log"))
	}
	addHook(logger)
	if err := mlogger.Init(mlogger.WithLevel(mlogger.DebugLevel)); err != nil {
		panic(err)
	}
}

func SetLevel(l string) error {
	ll, err := logrus.ParseLevel(l)
	if err != nil {
		return err
	}
	logger.SetLevel(ll)
	return nil
}

func GetLevel() string {
	return logger.GetLevel().String()
}

func Discard() logrus.FieldLogger {
	return &logrus.Logger{
		Out: io.Discard,
	}
}

func addHook(l *logrus.Logger) {
	l.AddHook(new(h1))
}

type h1 struct{}

func (h *h1) Levels() []logrus.Level {
	return []logrus.Level{}
}

func (h *h1) Fire(_ *logrus.Entry) error {
	return nil
}

func sortOutput(ks []string) {
	// sort.Strings(ks)
	for i := 0; i < len(ks); i++ {
		if ks[i] == logrus.FieldKeyMsg {
			ks[i], ks[len(ks)-1] = ks[len(ks)-1], ks[i]
			break
		}
	}
}
