package helper

import (
	"testing"
)

func TestOr(t *testing.T) {
	v := Or(0, 2)
	t.Log(v)

	res, err := A2I[int](I2A(123))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestGetMaxKeyAndValue(t *testing.T) {

}
