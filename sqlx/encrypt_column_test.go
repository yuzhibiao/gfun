package sqlx

import (
	"strings"
	"sync"
	"testing"
)

func TestSetEncryptKeyConcurrent(t *testing.T) {
	keys := [][]byte{
		[]byte("0123456789abcdef"),
		[]byte("0123456789abcdef01234567"),
		[]byte("0123456789abcdef0123456789abcdef"),
	}
	col := EncryptColumn[string]{Val: "secret", Valid: true}
	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = SetEncryptKey(keys[n%3])
			_, _ = col.Value()
			var c2 EncryptColumn[string]
			_ = c2.Scan("MTIz") // 无效密文，预期报错；仅为触发并发读取
		}(i)
	}
	wg.Wait()
}

func TestEncryptColumnScanNilWithoutKey(t *testing.T) {
	// 此测试必须位于 TestSetEncryptKey 之前运行（依赖全局密钥尚未设置）
	var col EncryptColumn[string]
	if err := col.Scan(nil); err != nil {
		t.Errorf("Scan(nil) without key = %v, want nil", err)
	}
	if col.Valid {
		t.Error("Scan(nil) Valid = true, want false")
	}
}

func TestSetEncryptKey(t *testing.T) {
	if err := SetEncryptKey([]byte("short")); err == nil {
		t.Error("SetEncryptKey(short) should fail")
	}
	if err := SetEncryptKey([]byte("0123456789abcdef")); err != nil {
		t.Errorf("SetEncryptKey(16 bytes) = %v, want nil", err)
	}
}

func TestEncryptColumnRoundTrip(t *testing.T) {
	if err := SetEncryptKey([]byte("0123456789abcdef0123456789abcdef")); err != nil {
		t.Fatal(err)
	}
	col := EncryptColumn[testUser]{Val: testUser{Name: "billyu", Age: 18}, Valid: true}

	got, err := col.Value()
	if err != nil {
		t.Fatal(err)
	}
	s, ok := got.(string)
	if !ok {
		t.Fatalf("Value() = %T, want string", got)
	}
	if strings.Contains(s, "billyu") {
		t.Error("ciphertext leaks plaintext")
	}

	var rt EncryptColumn[testUser]
	if err := rt.Scan(s); err != nil {
		t.Fatal(err)
	}
	if !rt.Valid || rt.Val != col.Val {
		t.Errorf("round trip = (%+v, %v), want (%+v, true)", rt.Val, rt.Valid, col.Val)
	}
}

func TestEncryptColumnNull(t *testing.T) {
	if err := SetEncryptKey([]byte("0123456789abcdef")); err != nil {
		t.Fatal(err)
	}
	col := EncryptColumn[testUser]{}
	v, err := col.Value()
	if err != nil || v != nil {
		t.Errorf("Value() = (%v, %v), want (nil, nil)", v, err)
	}

	var col2 EncryptColumn[testUser]
	if err := col2.Scan(nil); err != nil {
		t.Fatal(err)
	}
	if col2.Valid {
		t.Error("Scan(nil) Valid = true, want false")
	}
}

func TestEncryptColumnTampered(t *testing.T) {
	if err := SetEncryptKey([]byte("0123456789abcdef")); err != nil {
		t.Fatal(err)
	}
	col := EncryptColumn[string]{Val: "secret", Valid: true}
	v, err := col.Value()
	if err != nil {
		t.Fatal(err)
	}
	b := []byte(v.(string))
	if b[0] == 'A' {
		b[0] = 'B'
	} else {
		b[0] = 'A'
	}
	var rt EncryptColumn[string]
	if err := rt.Scan(string(b)); err == nil {
		t.Error("Scan(tampered) should fail (GCM auth)")
	}
}
