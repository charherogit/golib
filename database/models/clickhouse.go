package models

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"golib/async"
	"golib/config"
	log "golib/logutil"
	"reflect"
	"time"
)

func WithAddr(addr ...string) func(opt *clickhouse.Options) {
	return func(opt *clickhouse.Options) {
		opt.Addr = addr
	}
}

func WithAuth(db, name, pwd string) func(opt *clickhouse.Options) {
	return func(opt *clickhouse.Options) {
		opt.Auth = clickhouse.Auth{
			Database: db, Username: name, Password: pwd,
		}
	}
}

func NewClickhouse() (clickhouse.Conn, error) {
	return Clickhouse(func(opt *clickhouse.Options) {
		opt.Addr = config.C.CHAddr
		opt.Auth = clickhouse.Auth{
			Database: config.C.CHDB,
			Username: config.C.CHUser,
			Password: config.C.CHPassword,
		}
	})
}

func Clickhouse(opts ...func(opt *clickhouse.Options)) (clickhouse.Conn, error) {
	opt := &clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "default", Username: "default", Password: "poiuybnm",
		},
		Debug:        false,
		Debugf:       nil,
		Settings:     clickhouse.Settings{},
		Compression:  &clickhouse.Compression{Method: clickhouse.CompressionLZ4},
		DialTimeout:  30 * time.Second,
		MaxOpenConns: 5, MaxIdleConns: 5,
		ConnMaxLifetime:      10 * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
	}
	// clickhouse.OpenDB()
	for _, v := range opts {
		v(opt)
	}
	conn, err := clickhouse.Open(opt)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}

type NoStruct [0]int

type BatchWriter[T any] struct {
	conn  clickhouse.Conn
	ctx   context.Context
	query string

	bt      *async.BatchTask[T]
	lastErr error
	logger  log.ILogger
}

// XXX 暂时支持只支持*struct, []any, [n]any， 能用*struct尽量用，因为切片没有验证列长度，容易出现插入失败
func NewBatchWriter[T any](lg log.ILogger, c clickhouse.Conn) *BatchWriter[T] {
	if lg == nil {
		lg = log.StandardLogger()
	}
	return &BatchWriter[T]{
		logger: lg,
		conn:   c,
	}
}

func (b *BatchWriter[T]) Close() error {
	b.bt.Stop()
	b.logger.Infof("%s %s", b.query, b.bt.Metric())
	return b.Error()
}

func (b *BatchWriter[T]) Error() error {
	return b.lastErr
}

func (b *BatchWriter[T]) send(data []T) {
	batch, err := b.conn.PrepareBatch(b.ctx, b.query, driver.WithReleaseConnection())
	if err != nil {
		b.logger.Error(err)
		b.lastErr = err
		return
	}
	for _, v := range data {
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Array:
			slice := make([]any, val.Len())
			for i := 0; i < val.Len(); i++ {
				slice[i] = val.Index(i).Interface()
			}
			if err := batch.Append(slice...); err != nil {
				b.logger.Error(err)
				b.lastErr = err
			}
		case reflect.Slice:
			if err := batch.Append(any(v).([]any)...); err != nil {
				b.logger.Error(err)
				b.lastErr = err
			}
		case reflect.Ptr:
			if err := batch.AppendStruct(v); err != nil {
				b.logger.Error(err)
				b.lastErr = err
			}
		default:
			b.logger.Errorf("can not handle %T", v)
		}
	}
	if err := batch.Send(); err != nil {
		b.logger.Error(err)
		b.lastErr = err
	}
}

func (b *BatchWriter[T]) Init(ctx context.Context, query string) *BatchWriter[T] {
	b.ctx = ctx
	b.query = query
	b.bt = async.NewBatchTask(3000, b.send)
	b.bt.Watch()
	return b
}

func (b *BatchWriter[T]) T(ctx context.Context, table string) *BatchWriter[T] {
	return b.Init(ctx, fmt.Sprintf("INSERT INTO %s VALUES ()", table))
}

func (b *BatchWriter[T]) Append(v T) {
	b.bt.Add(v)
}
