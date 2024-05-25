package models

import (
	"context"
	"errors"
	"fmt"
	"golib/config"
	pkg "golib/helper"
	log "golib/logutil"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Name    string              // Collection name
	Struct  string              // 存储的结构
	DB      string              // DB
	Desc    string              // 说明
	Indexes []*mongo.IndexModel // 索引
}

func (m *MongoConfig) String() string {
	return fmt.Sprintf("db: %s collection name: %s struct: %s desc: %s", m.DB, m.Name, m.Struct, m.Desc)
}

const (
	defaultMinMongoPoolSize uint64 = 10
	defaultMaxMongoPoolSize uint64 = 100
)

var (
	client *mongo.Client
	once   sync.Once

	defaultConf = &MongoOptionConf{
		Uri:         config.C.MongoDBAddr,
		MinPoolSize: defaultMinMongoPoolSize,
		MaxPoolSize: defaultMaxMongoPoolSize,
	}
)

type MongoOptionConf struct {
	MinPoolSize uint64
	MaxPoolSize uint64
	Uri         string
}

func (b *MongoOptionConf) GetUri() string {
	if len(b.Uri) == 0 {
		return config.C.MongoDBAddr
	}
	return b.Uri
}

func (b *MongoOptionConf) GetMinPool() uint64 {
	if b.MinPoolSize == 0 {
		return defaultMinMongoPoolSize
	}
	return b.MinPoolSize
}

func (b *MongoOptionConf) GetMaxPool() uint64 {
	if b.MaxPoolSize == 0 {
		return defaultMaxMongoPoolSize
	}
	return b.MaxPoolSize
}

func GetCol(conf MongoConfig) *mongo.Collection {
	col := GetDb(conf).Collection(conf.Name)
	if err := InitIndexes(context.Background(), col, &conf); err != nil {
		log.Fatal(err)
	}
	return col
}

func GetDb(conf MongoConfig) *mongo.Database {
	// log.Tracef("db: %s collection: %s", conf.DB, conf.Name)
	once.Do(func() {
		client = NewClient(defaultConf)
	})

	return client.Database(conf.DB)
}

func UpdateOneAndCheck(ctx context.Context, col *mongo.Collection, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) error {
	result, err := col.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 && result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		return fmt.Errorf("update not success filter: %v update: %v opts: %v", filter, update, opts)
	}
	return nil
}

func InitIndexes(ctx context.Context, c *mongo.Collection, cfg *MongoConfig) error {
	for _, v := range cfg.Indexes {
		if err := checkIndex(ctx, c, v.Options); err != nil {
			if errors.Is(err, indexNotFound) {
				if _, err := c.Indexes().CreateOne(ctx, *v); err != nil {
					return fmt.Errorf("col: %s create index: %s err: %s", cfg.Name, *v.Options.Name, err)
				}
			} else {
				return fmt.Errorf("col: %s check index: %s err: %s", cfg.Name, *v.Options.Name, err)
			}
		}
	}
	return nil
}

var indexNotFound = errors.New("mongodb index not found")

func checkIndex(ctx context.Context, c *mongo.Collection, opt *options.IndexOptions) error {
	type index struct {
		Key  map[string]int
		Name string
	}
	indexes, err := c.Indexes().List(ctx)
	if err != nil {
		return err
	}
	defer indexes.Close(ctx)

	for indexes.Next(ctx) {
		var idx index
		if err = indexes.Decode(&idx); err != nil {
			return err
		}
		if idx.Name == *opt.Name {
			return nil
		}
	}
	return indexNotFound
}

type MongoUpdateData = MongoUpdate

type MongoUpdate map[string]map[string]interface{}

func MongoSet(k string, v interface{}) MongoUpdate {
	return MongoUpdate{}.Set(k, v)
}

func MongoUpdateWithTime() MongoUpdate {
	return MongoUpdate{}.Set("update_time", time.Now().Unix())
}

func (m MongoUpdate) GetData() map[string]map[string]interface{} {
	if m == nil {
		return nil
	}
	return m
}

func (m MongoUpdate) update(act, key string, value interface{}) MongoUpdate {
	if m == nil {
		m = make(map[string]map[string]interface{})
		m[act] = map[string]interface{}{
			key: value,
		}
	} else {
		if v, ok := m[act]; ok {
			v[key] = value
		} else {
			m[act] = map[string]interface{}{
				key: value,
			}
		}
	}
	return m
}

func (m MongoUpdate) Set(k string, v interface{}) MongoUpdate {
	return m.update("$set", k, v)
}

func (m MongoUpdate) Unset(k string) MongoUpdate {
	return m.update("$unset", k, 1)
}

func (m MongoUpdate) Inc(k string, v interface{}) MongoUpdate {
	return m.update("$inc", k, v)
}

func (m MongoUpdate) Push(k string, v interface{}) MongoUpdate {
	return m.update("$push", k, v)
}

func (m MongoUpdate) Pull(k string, v interface{}) MongoUpdate {
	return m.update("$pull", k, v)
}

func GetObjectIdFilter(id string) (primitive.M, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return bson.M{"_id": objectId}, nil
}

func GetIdFilter(userId uint32) primitive.M {
	return bson.M{"_id": userId}
}

func ConvertStrToObjectId(str string) (primitive.ObjectID, error) {
	objectId, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		return [12]byte{}, fmt.Errorf("convert: %s to object id failed: %s", str, err)
	}
	return objectId, nil
}

func ConvertStrsToObjectIdList(strs []string) ([]primitive.ObjectID, error) {
	var objectIds []primitive.ObjectID
	for _, str := range strs {
		objectId, err := primitive.ObjectIDFromHex(str)
		if err != nil {
			return nil, fmt.Errorf("convert: %s to object id failed: %s", str, err)
		}
		objectIds = append(objectIds, objectId)
	}
	return objectIds, nil
}

func NoDoc(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}

func NilDocs(err error) bool {
	return errors.Is(err, mongo.ErrNilDocument)
}
func NewClient(conf *MongoOptionConf) *mongo.Client {
	opt := options.Client().ApplyURI(conf.GetUri())
	opt.SetMinPoolSize(conf.GetMinPool())
	opt.SetMaxPoolSize(conf.GetMaxPool())
	ctx, cancel := pkg.GetCtxTimeOut10()
	defer cancel()
	var err error
	client, err = mongo.Connect(ctx, opt)
	if err != nil {
		panic("connect " + conf.GetUri() + " err: " + err.Error())
	}
	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	log.Infof("Uri: %s MinPoolSize: %d MaxPoolSize: %d", conf.GetUri(), conf.GetMinPool(), conf.GetMaxPool())
	return client
}
