package caller

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBriefInfo(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		skip int
	}
	tests := []struct {
		name string
		args args
		want bytes.Buffer
	}{
		{
			name: "simple",
			args: args{
				skip: 1,
			},
			want: *bytes.NewBuffer([]byte("[caller_test.go:func1:31]")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BriefInfo(tt.args.skip); !reflect.DeepEqual(got, tt.want) {
				assert.Equal(got, tt.want, "test: ["+tt.name+"] got not match want")
			}
		})
	}
}

func TestTrimFn(t *testing.T) {
	assert.Equal(t, TrimFn("a.b.c"), "c")
	assert.Equal(t, TrimFn("a.b"), "b")
	assert.Equal(t, TrimFn("a"), "a")
	assert.Equal(t, TrimFn(""), "")
}

func TestFuncName(t *testing.T) {
	const e = "not a function"
	var x = 1
	var y = "1"
	assert.Equal(t, e, FuncName(x))
	assert.Equal(t, e, FuncName(y))
	assert.Equal(t, e, FuncName(&x))
	assert.Equal(t, "func1", FuncName(func() {}))
	assert.Equal(t, "func2", FuncName(func() {}))
	assert.Equal(t, "TestFuncName", FuncName(TestFuncName))
}
