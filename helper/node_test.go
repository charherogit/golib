package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortSlice(t *testing.T) {
	type args struct {
		s []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "SortSlice",
			args: args{
				s: []int{1, 31, 222, 52, 222, 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, SortSlice(tt.args.s), "SortSlice(%v)", tt.args.s)
		})
	}
}
