package sqlx

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

// EncryptColumn 写库时用 AES-GCM 加密、读库时解密，避免敏感数据明文落库。
// 密文以 base64 字符串存储，数据库列类型建议 TEXT / VARCHAR。
//
// 使用前必须调用 SetEncryptKey 设置全局密钥。
//
// 用法：
//
//	type User struct {
//	    Phone sqlx.EncryptColumn[string] // 手机号等敏感字段
//	}
type EncryptColumn[T any] struct {
	Val   T
	Valid bool
}

var encryptKey []byte

// SetEncryptKey 设置全局加密密钥，长度必须为 16 / 24 / 32 字节
// （对应 AES-128 / 192 / 256）。应在程序初始化阶段调用一次。
func SetEncryptKey(k []byte) error {
	switch len(k) {
	case 16, 24, 32:
		encryptKey = k
		return nil
	default:
		return errors.New("sqlx: encrypt key must be 16, 24 or 32 bytes")
	}
}

// Value 实现 driver.Valuer：序列化为 JSON 后 AES-GCM 加密、base64 编码。
func (c EncryptColumn[T]) Value() (driver.Value, error) {
	if !c.Valid {
		return nil, nil
	}
	if len(encryptKey) == 0 {
		return nil, errors.New("sqlx: encrypt key not set, call SetEncryptKey first")
	}
	b, err := json.Marshal(c.Val)
	if err != nil {
		return nil, fmt.Errorf("sqlx: marshal EncryptColumn: %w", err)
	}
	return encrypt(b, encryptKey)
}

// Scan 实现 sql.Scanner：base64 解码、解密后反序列化为 T。
func (c *EncryptColumn[T]) Scan(src any) error {
	if src == nil {
		var zero T
		c.Val, c.Valid = zero, false
		return nil
	}
	if len(encryptKey) == 0 {
		return errors.New("sqlx: encrypt key not set, call SetEncryptKey first")
	}
	var b []byte
	switch v := src.(type) {
	case string:
		b = []byte(v)
	case []byte:
		b = v
	default:
		return fmt.Errorf("sqlx: cannot scan %T into EncryptColumn", src)
	}
	// Value() 存的是 base64 字符串，先解码再解密
	data, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return fmt.Errorf("sqlx: decode EncryptColumn: %w", err)
	}
	pt, err := decrypt(data, encryptKey)
	if err != nil {
		return fmt.Errorf("sqlx: decrypt EncryptColumn: %w", err)
	}
	if err := json.Unmarshal(pt, &c.Val); err != nil {
		return fmt.Errorf("sqlx: unmarshal EncryptColumn: %w", err)
	}
	c.Valid = true
	return nil
}

// encrypt 使用 AES-GCM 加密：随机 nonce 拼接在密文前，整体 base64 编码。
func encrypt(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	// Seal 的第一个参数是 dst，不能用 nonce 本身；先加密再把 nonce 拼在密文前
	out := append(nonce, gcm.Seal(nil, nonce, plaintext, nil)...)
	return base64.StdEncoding.EncodeToString(out), nil
}

// decrypt 是 encrypt 的逆操作。密文被篡改时 GCM 校验失败，返回错误。
func decrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return nil, errors.New("sqlx: ciphertext too short")
	}
	return gcm.Open(nil, data[:ns], data[ns:], nil)
}
