package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	// log "github.com/sirupsen/logrus"
)

type phase int

const (
	preGame phase = iota
	inGame
	postGame
)

type gameType int

const (
	singlePlayer gameType = iota
	multiPlayer
)

type player struct {
	name  string
	coin  byte
	isBot bool
}

type board [][]byte

// cursor represents the moving cursor index
// nextChance represents next coin 'X' or 'O'
// matrix is board
// game ended represent if the game is ended
type gameState struct {
	matrix   board
	nextByte byte
	cursor   int
	phase    phase
	gameType gameType
	players  [2]player
}

func toggle(a byte) byte {
	if a == 'X' {
		return 'O'
	}
	return 'X'
}

func initialGameState() *gameState {
	m := make([][]byte, 8)
	for i := range m {
		m[i] = make([]byte, 8)
	}
	return &gameState{
		phase:    preGame,
		matrix:   m,
		nextByte: 'X',
		cursor:   0,
	}
}

func getCursor(col, cur int) string {
	var s string
	for j := 0; j < col; j++ {
		if j == cur {
			s += " ! "
			continue
		}
		s += "   "
	}
	return s + "\n"
}

func (gs *gameState) View() string {
	// The header
	s := "Play Connect 4!\n\n"
	if gs.phase == preGame {
		return s + " 1. Single Player\n 2. Two Player\n"
	}
	if gs.gameType == singlePlayer {
		return s + "Single player in construction\n"
	}
	s += "Its is " + string(gs.nextByte) + " turn\n"
	// Iterate over our choices
	s += getCursor(len(gs.matrix[0]), gs.cursor)
	for _, row := range gs.matrix {
		var pr string
		for _, e := range row {
			// pr += string(e)
			if e == 0 {
				pr += " _ "
			} else {
				pr += (" " + string(e) + " ")
			}
		}
		s += pr + "\n"
	}

	if gs.phase == postGame {
		return s + "Player " + string(gs.nextByte) + " Won!"
	}
	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func (gs *gameState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if gs.phase == preGame {
		switch msg := msg.(type) {
		case tea.KeyMsg:

			switch msg.String() {

			case "1":
				gs.gameType = singlePlayer
				gs.phase = inGame
			case "2":
				gs.gameType = multiPlayer
				gs.phase = inGame
			}
		}
		return gs, nil
	}
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return gs, tea.Quit

		// The "left" keys move the cursor left
		case "left":
			if gs.cursor > 0 {
				gs.cursor--
			}

		// The "right" keys move the cursor right
		case "right":
			if gs.cursor < len(gs.matrix[0])-1 {
				gs.cursor++
			}

		// to drop your coin in the given column
		case "enter", " ":
			if gs.phase == inGame {
				var r int
				for r = len(gs.matrix) - 1; r >= 0; r-- {
					if gs.matrix[r][gs.cursor] == 0 {
						gs.matrix[r][gs.cursor] = gs.nextByte
						break
					}
				}
				if gs.matrix.is4Connected(r, gs.cursor, gs.nextByte) {
					gs.phase = postGame
				} else {
					gs.nextByte = toggle(gs.nextByte)
				}
			}
		}
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return gs, nil
}

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

func (b board) is4Connected(r, c int, player byte) bool {
	for _, l := range listarr {
		totalTokenFound := 0
		for i := range l.rowdiff {
			if r+l.rowdiff[i] > 0 && r+l.rowdiff[i] < len(b) && c+l.coldiff[i] > 0 && c+l.coldiff[i] < len(b[0]) && b[r+l.rowdiff[i]][c+l.coldiff[i]] == player {
				totalTokenFound++
			}
		}
		if totalTokenFound == 3 {
			return true
		}
	}
	return false
}

func (gs *gameState) Init() tea.Cmd {
	return nil
}

func main() {
	p := tea.NewProgram(initialGameState())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
