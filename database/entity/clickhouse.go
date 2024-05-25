package entity

import (
	"time"
)

// 物资变更数据
type MaterialChange struct {
	UserId      uint32    `ch:"user_id"`
	CreatedAt   time.Time `ch:"created_at"`
	Level       uint32    `ch:"level"`        // 玩家等级
	Variable    uint32    `ch:"variable"`     // 变更物资
	AffectType  uint32    `ch:"affect_type"`  // 变更原因
	ChangeCount int32     `ch:"change_count"` // 变更数量
	Total       int64     `ch:"total"`        // 最总值
}
