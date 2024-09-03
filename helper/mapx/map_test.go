package mapx

import (
	. "golang.org/x/exp/constraints"
	"golib/helper/numx"
	"testing"
)

func TestOr(t *testing.T) {
	v := numx.Or(0, 2)
	t.Log(v)

	res, err := numx.A2I[int](numx.I2A(123))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestEqualMap(t *testing.T) {
	type args[T comparable, U Integer] struct {
		m map[T]U
		n map[T]U
	}
	type testCase[T comparable, U Integer] struct {
		name    string
		args    args[T, U]
		wantErr bool
	}
	tests := []testCase[uint32, uint32]{
		{
			name: "len not equal",
			args: args[uint32, uint32]{
				m: map[uint32]uint32{
					1: 1,
					2: 2,
					3: 3,
				},
				n: map[uint32]uint32{
					1: 1,
					2: 2,
					3: 3,
					4: 4,
				},
			},
			wantErr: true,
		},
		{
			name: "count not equal",
			args: args[uint32, uint32]{
				m: map[uint32]uint32{
					1: 1,
					2: 2,
					3: 3,
				},
				n: map[uint32]uint32{
					1: 1,
					2: 2,
					3: 4,
				},
			},
			wantErr: true,
		},
		{
			name: "key not exist",
			args: args[uint32, uint32]{
				m: map[uint32]uint32{
					1: 1,
					2: 2,
					3: 3,
				},
				n: map[uint32]uint32{
					1: 1,
					2: 2,
					4: 4,
				},
			},
			wantErr: true,
		},
		{
			name: "equal",
			args: args[uint32, uint32]{
				m: map[uint32]uint32{
					1: 1,
					2: 2,
					3: 3,
				},
				n: map[uint32]uint32{
					1: 1,
					2: 2,
					3: 3,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EqualMap(tt.args.m, tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("EqualMap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContainMap(t *testing.T) {
	type args[T comparable, U Integer] struct {
		m map[T]U
		n map[T]U
	}
	type testCase[T comparable, U Integer] struct {
		name    string
		args    args[T, U]
		wantErr bool
	}
	tests := []testCase[uint32, uint32]{
		{
			name: "count > max",
			args: args[uint32, uint32]{
				m: map[uint32]uint32{
					1: 1,
					2: 2,
					3: 4,
				},
				n: map[uint32]uint32{
					1: 1,
					2: 2,
					3: 3,
				},
			},
			wantErr: true,
		},
		{
			name: "key not exist",
			args: args[uint32, uint32]{
				m: map[uint32]uint32{
					1: 1,
					2: 2,
					3: 3,
				},
				n: map[uint32]uint32{
					1: 1,
					2: 2,
				},
			},
			wantErr: true,
		},
		{
			name: "equal",
			args: args[uint32, uint32]{
				m: map[uint32]uint32{
					1: 1,
					2: 2,
					3: 3,
				},
				n: map[uint32]uint32{
					1: 1,
					2: 2,
					3: 3,
				},
			},
			wantErr: false,
		},
		{
			name: "less",
			args: args[uint32, uint32]{
				m: map[uint32]uint32{
					1: 1,
					2: 1,
				},
				n: map[uint32]uint32{
					1: 1,
					2: 2,
					3: 3,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ContainMap(tt.args.m, tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("ContainMap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
