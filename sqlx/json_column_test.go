package sqlx

import "testing"

type testUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestJsonColumnValue(t *testing.T) {
	col := JsonColumn[testUser]{Val: testUser{Name: "tom", Age: 18}, Valid: true}
	got, err := col.Value()
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"name":"tom","age":18}`; got != want {
		t.Errorf("Value() = %v, want %s", got, want)
	}

	// 无效值写 NULL
	col2 := JsonColumn[testUser]{}
	got2, err := col2.Value()
	if err != nil {
		t.Fatal(err)
	}
	if got2 != nil {
		t.Errorf("Value() = %v, want nil", got2)
	}
}

func TestJsonColumnScan(t *testing.T) {
	var col JsonColumn[testUser]
	if err := col.Scan(`{"name":"jerry","age":20}`); err != nil {
		t.Fatal(err)
	}
	if !col.Valid || col.Val.Name != "jerry" || col.Val.Age != 20 {
		t.Errorf("Scan(string) = (%+v, %v), want ({jerry 20}, true)", col.Val, col.Valid)
	}

	var col2 JsonColumn[testUser]
	if err := col2.Scan([]byte(`{"name":"spike","age":3}`)); err != nil {
		t.Fatal(err)
	}
	if !col2.Valid || col2.Val.Name != "spike" {
		t.Errorf("Scan([]byte) = (%+v, %v)", col2.Val, col2.Valid)
	}

	var col3 JsonColumn[testUser]
	if err := col3.Scan(nil); err != nil {
		t.Fatal(err)
	}
	if col3.Valid {
		t.Error("Scan(nil) Valid = true, want false")
	}

	var col4 JsonColumn[testUser]
	if err := col4.Scan(12345); err == nil {
		t.Error("Scan(int) should fail")
	}
}
