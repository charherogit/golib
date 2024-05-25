package models

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"golib/database/entity"
	log "golib/logutil"
	"sync"
	"testing"
	"time"
)

const createTestTableSql = `
CREATE TABLE IF NOT EXISTS hello 
(
    user_id UInt32,
    updated_at DateTime DEFAULT now(),
    updated_at_date Date DEFAULT toDate(updated_at)
) 
ENGINE = MergeTree
ORDER BY user_id;
`
const dropTestTableSql = `DROP TABLE IF EXISTS hello`

func createTestTable(t *testing.T) clickhouse.Conn {
	ctx := context.Background()
	conn, err := Clickhouse(WithAddr("192.168.196.18:9000"), WithAuth("test", "default", ""))
	if err != nil {
		t.Fatal(err)
	}
	if err := conn.Exec(ctx, createTestTableSql); err != nil {
		t.Fatal(err)
	}
	if err := conn.Exec(ctx, "TRUNCATE TABLE hello"); err != nil {
		t.Fatal(err)
	}
	return conn
}

type object struct {
	UserID int `ch:"user_id"`
}

func TestAppendStruct(t *testing.T) {
	conn := createTestTable(t)

	bw := NewBatchWriter[*object](nil, conn)
	bw.Init(context.Background(), "INSERT INTO hello(user_id) VALUES()")
	defer bw.Close()

	bw.Append(&object{
		UserID: 11111,
	})
	if err := bw.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestAppendSlice(t *testing.T) {
	conn := createTestTable(t)

	bw := NewBatchWriter[[]any](nil, conn)
	bw.Init(context.Background(), "INSERT INTO hello(user_id) VALUES()")
	defer bw.Close()

	bw.Append([]any{22222})
	if err := bw.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestAppendArr(t *testing.T) {
	conn := createTestTable(t)

	bw := NewBatchWriter[[1]any](nil, conn)
	bw.Init(context.Background(), "INSERT INTO hello(user_id) VALUES()")
	defer bw.Close()

	bw.Append([1]any{33333})
	if err := bw.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestASend(t *testing.T) {
	conn := createTestTable(t)

	bw := NewBatchWriter[[]any](nil, conn)
	bw.Init(context.Background(), "INSERT INTO hello (user_id) VALUES()")
	defer bw.Close()

	wg := sync.WaitGroup{}
	n := 5
	wg.Add(n)
	for j := 0; j < n; j++ {
		go func() {
			for i := 0; i < 100000; i++ {
				bw.Append([]any{i})
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

const createMaterialTableSql = `
CREATE TABLE IF NOT EXISTS material_change
(
  user_id UInt32,
  created_at DATETIME DEFAULT now(),
  level UInt32 COMMENT '玩家等级',
  variable UInt32 COMMENT '变更物资',
  affect_type UInt32 COMMENT '变更原因',
  change_count Int32 COMMENT '变更数量',
  total Int64 COMMENT '最总值'
)
ENGINE MergeTree
ORDER BY toYYYYMMDD(created_at)
COMMENT '物资变更数据';
`
const dropMaterialTableSql = `DROP TABLE IF EXISTS material_change`

func createMaterialTable(b *testing.B) clickhouse.Conn {
	ctx := context.Background()
	conn, err := Clickhouse(WithAddr("192.168.196.18:9000"), WithAuth("test", "default", ""))
	if err != nil {
		b.Fatal(err)
	}
	if err := conn.Exec(ctx, createMaterialTableSql); err != nil {
		b.Fatal(err)
	}
	if err := conn.Exec(ctx, `TRUNCATE TABLE material_change`); err != nil {
		b.Fatal(err)
	}
	return conn
}

func BenchmarkAppendMaterialByObject(b *testing.B) {
	conn := createMaterialTable(b)

	bw := NewBatchWriter[*entity.MaterialChange](log.StandardLogger(), conn)
	bw.Init(context.Background(), "INSERT INTO material_change VALUES ()")
	defer bw.Close()

	for i := 0; i < b.N; i++ {
		bw.Append(&entity.MaterialChange{
			UserId: uint32(i), CreatedAt: time.Now(), Level: 1, Variable: 1, AffectType: 1, ChangeCount: 1, Total: 1,
		})
	}
}

func BenchmarkAppendMaterialBySlice(b *testing.B) {
	conn := createMaterialTable(b)

	bw := NewBatchWriter[[]any](log.StandardLogger(), conn)
	bw.Init(context.Background(), "INSERT INTO material_change VALUES ()")
	defer bw.Close()

	for i := 0; i < b.N; i++ {
		bw.Append([]any{i, time.Now().Unix(), 2, 2, 2, 2, 2})
	}
}

func BenchmarkAppendMaterialByArr(b *testing.B) {
	conn := createMaterialTable(b)

	bw := NewBatchWriter[[7]any](log.StandardLogger(), conn)
	bw.Init(context.Background(), "INSERT INTO material_change VALUES ()")
	defer bw.Close()

	for i := 0; i < b.N; i++ {
		bw.Append([7]any{i, time.Now().Unix(), 3, 3, 3, 3, 3})
	}
}

func TestAppendMaterial(t *testing.T) {
	results := testing.Benchmark(BenchmarkAppendMaterialByObject)
	fmt.Println("BenchmarkAppendMaterialByObject:", results)

	results = testing.Benchmark(BenchmarkAppendMaterialBySlice)
	fmt.Println("BenchmarkAppendMaterialBySlice:", results)

	results = testing.Benchmark(BenchmarkAppendMaterialByArr)
	fmt.Println("BenchmarkAppendMaterialByArr:", results)
}
