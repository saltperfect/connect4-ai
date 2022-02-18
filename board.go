package main

var window_length int = 4

type board [][]byte

type rowcol struct {
	rowdiff []int
	coldiff []int
}

var listarr = []rowcol{
	{[]int{1, 2, 3}, []int{0, 0, 0}},
	{[]int{0, 0, 0}, []int{1, 2, 3}},
	{[]int{0, 0, 0}, []int{-1, -2, -3}},
	{[]int{1, 2, 3}, []int{-1, -2, -3}},
	{[]int{1, 2, 3}, []int{1, 2, 3}},
	{[]int{-1, -2, -3}, []int{-1, -2, -3}},
	{[]int{-1, -2, -3}, []int{1, 2, 3}},
	{[]int{1, -1, -2}, []int{1, -1, -2}},
	{[]int{1, 2, -1}, []int{1, 2, -1}},
	{[]int{1, -1, -2}, []int{-1, 1, 2}},
	{[]int{1, 2, -1}, []int{-1, -2, 1}},
	{[]int{0, 0, 0}, []int{-1, 1, 2}},
	{[]int{0, 0, 0}, []int{-1, -2, 1}},
}

func (b board) isValidLocation(col int) bool {
	return b[0][col] == 0
}

func (b board) getOpenRow(c int) int {
	var r int
	for r = len(b) - 1; r >= 0; r-- {
		if b[r][c] == 0 {
			break
		}
	}
	return r
}

func (b board) getOpenColList() []int {
	var list []int
	for i := 0; i < len(b[0]); i++ {
		if b.isValidLocation(i) {
			list = append(list, i)
		}
	}
	return list
}

func (b board) is4Connected(r, c int) bool {
	for _, l := range listarr {
		totalTokenFound := 0
		for i := range l.rowdiff {
			if r+l.rowdiff[i] > 0 && r+l.rowdiff[i] < len(b) && c+l.coldiff[i] > 0 && c+l.coldiff[i] < len(b[0]) && b[r+l.rowdiff[i]][c+l.coldiff[i]] == b[r][c] {
				totalTokenFound++
			}
		}
		if totalTokenFound == 3 {
			return true
		}
	}
	return false
}

func (b board) dropCoin(c int, coin byte) {
	r := b.getOpenRow(c)
	b[r][c] = coin
}

func (b board) copy() board {
	bf := make([][]byte, len(b))
	for i := range bf {
		bf[i] = copyS(b[i])
	}
	return bf
}

func copyS(b []byte) []byte {
	t := make([]byte, len(b))
	copy(t, b)
	return t
}

type window []byte

func (b board) scorePosition(coin byte) int {
	score := 0

	mid := len(b[0]) / 2

	centreCol := window(b[mid])

	centre_count := centreCol.count(coin)

	score += centre_count * 3

	for _, r := range b {
		for c := 0; c < len(b[0])-3; c++ {
			w := window(r[c : c+window_length])
			score += w.evalute_window(coin)
		}
	}

	for c := 0; c < len(b[0]); c++ {
		for r := 0; r < len(b)-3; r++ {
			w := window(b.getColum(c)[r : r+window_length])
			score += w.evalute_window(coin)
		}
	}

	for r := 0; r < len(b)-3; r++ {
		for c :=0; c < len(b[0])-3;c++ {
			w := make([]byte, 4)
			for i := 0; i < 4; i++ {
				w[i] = b[r+i][c+i]
			}
			score += window(w).evalute_window(coin)
		}
	}

	for r := 0; r < len(b)-3; r++ {
		for c := 0; c < len(b[0])-3; c++ {
			w := make([]byte, 4)
			for i := 0; i < 4; i++ {
				w[i] = b[r+3-i][c+i]
			}
			score += window(w).evalute_window(coin)
		}
	}
	return score
}

func (b board) isTerminalNode() bool {
	return b.winingMove(ai_coin) || b.winingMove(player_coin) || len(b.getOpenColList()) == 0
}

func (b board) getColum(col int) []byte {
	cols := make([]byte, 0)

	for i := 0; i < len(b); i++ {
		cols = append(cols, b[i][col])
	}
	return cols
}

func (w window) evalute_window(coin byte) int {
	score := 0
	opp_coin := player_coin
	if coin == player_coin {
		opp_coin = ai_coin
	}

	if w.count(coin) == 4 {
		score += 100
	} else if w.count(coin) == 3 && w.count(0) == 1 {
		score += 5
	} else if w.count(coin) == 2 && w.count(0) == 2 {
		score += 2
	} else if w.count(opp_coin) == 3 && w.count(0) == 1 {
		score -= 4
	}
	return score
}

func (b board) winingMove(coin byte) bool {
	for c := 0; c < len(b[0])-3; c++ {
		for r := 0; r < len(b); r++ {
			if b[r][c] == coin && b[r][c+1] == coin && b[r][c+2] == coin && b[r][c+3] == coin {
				return true
			}
		}
	}
	for c := 0; c < len(b[0]); c++ {
		for r := 0; r < len(b)-3; r++ {
			if b[r][c] == coin && b[r+1][c] == coin && b[r+2][c] == coin && b[r+3][c] == coin {
				return true
			}
		}
	}
	for c := 0; c < len(b[0])-3; c++ {
		for r := 0; r < len(b)-3; r++ {
			if b[r][c] == coin && b[r+1][c+1] == coin && b[r+2][c+2] == coin && b[r+3][c+3] == coin {
				return true
			}
		}
	}
	for c := 0; c < len(b[0])-3; c++ {
		for r := 3; r < len(b); r++ {
			if b[r][c] == coin && b[r-1][c+1] == coin && b[r-2][c+2] == coin && b[r-3][c+3] == coin {
				return true
			}
		}
	}
	return false
}

func (w window) count(coin byte) int {
	cnt := 0
	for _, x := range w {
		if x == coin {
			cnt++
		}
	}
	return cnt
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
