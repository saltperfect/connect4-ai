package main

import "testing"

func Test_board_winingMove(t *testing.T) {
	type args struct {
		coin byte
	}
	tests := []struct {
		name string
		b    board
		args args
		want bool
	}{
		{
			name: "f",
			b: [][]byte{
				{'A', 'A', 'A', 'A', 'A'},
				{'A', 'C', 'A', 'A', 'C'},
				{'A', 'A', 'C', 'C', 'A'},
				{'A', 'A', 'C', 'C', 'A'},
				{'A', 'C', 'C', 'C', 'A'},
			},
			args: args{
				coin: 'C',
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.winingMove(tt.args.coin); got != tt.want {
				t.Errorf("board.winingMove() = %v, want %v", got, tt.want)
			}
		})
	}
}
