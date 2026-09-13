// Package sqlx 提供数据库列类型转换：
// JsonColumn 用于 JSON 列与 Go 类型互转，EncryptColumn 用于落库加密。
//
// 两者均实现 driver.Valuer / sql.Scanner 接口，
// GORM 和 database/sql 原生支持，可直接作为实体字段类型。
package sqlx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JsonColumn 在数据库 JSON 列与 Go 类型 T 之间自动转换，
// 并用 Valid 标记区分 NULL 与零值。
//
// 用法：
//
//	type User struct {
//	    Ext sqlx.JsonColumn[map[string]any] // 对应数据库 JSON 列
//	}
type JsonColumn[T any] struct {
	Val   T
	Valid bool
}

// Value 实现 driver.Valuer：写库时序列化为 JSON 字符串，无效值写 NULL。
func (c JsonColumn[T]) Value() (driver.Value, error) {
	if !c.Valid {
		return nil, nil
	}
	b, err := json.Marshal(c.Val)
	if err != nil {
		return nil, fmt.Errorf("sqlx: marshal JsonColumn: %w", err)
	}
	return string(b), nil
}

// Scan 实现 sql.Scanner：读库时反序列化，支持 []byte、string 和 NULL。
func (c *JsonColumn[T]) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		var zero T
		c.Val, c.Valid = zero, false
		return nil
	case []byte:
		return c.scanJSON(v)
	case string:
		return c.scanJSON([]byte(v))
	default:
		return fmt.Errorf("sqlx: cannot scan %T into JsonColumn", src)
	}
}

func (c *JsonColumn[T]) scanJSON(b []byte) error {
	if err := json.Unmarshal(b, &c.Val); err != nil {
		return fmt.Errorf("sqlx: unmarshal JsonColumn: %w", err)
	}
	c.Valid = true
	return nil
}
