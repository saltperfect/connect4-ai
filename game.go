package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type phase int

var player_coin byte
var ai_coin byte = 'C'

const (
	preGame phase = iota
	inGame
	postGame
	errorPhase
)

type player struct {
	name  string
	coin  byte
	isBot bool
	cmd   func(matrix board, opponent byte) tea.Cmd
}

// cursor represents the moving cursor index
// nextChance represents next coin 'X' or 'O'
// matrix is board
// game ended represent if the game is ended
type gameState struct {
	matrix        board
	currentPlayer int
	cursor        int
	phase         phase
	gameType      gameType
	players       [2]*player
}

func toggle(a int) int {
	if a == 0 {
		return 1
	}
	return 0
}

func initialGameState(inputs []textinput.Model) *gameState {

	m := make([][]byte, 8)
	for i := range m {
		m[i] = make([]byte, 7)
	}

	gameType, players, err := processText(inputs)
	if err != nil {
		return &gameState{
			phase: errorPhase,
		}
	}
	for _, p := range players {
		fmt.Printf("player: %+v\n", p)
	}
	return &gameState{
		phase:         inGame,
		matrix:        m,
		currentPlayer: 0,
		cursor:        0,
		gameType:      gameType,
		players:       players,
	}
}

func processText(inputs []textinput.Model) (gt gameType, players [2]*player, err error) {
	gt = singlePlayer
	if len(inputs) > 2 {
		gt = multiPlayer
	}

	for i := 0; i < len(inputs); i += 2 {
		if inputs[i].Value() == "" || inputs[i+1].Value() == "" {
			return gt, players, fmt.Errorf("no valid input found")
		}
		players[i/2] = &player{
			name: inputs[i].Value(),
			coin: []byte(inputs[i+1].Value())[0],
		}
	}
	if gt == singlePlayer {
		player_coin = players[0].coin
		players[1] = &player{
			name:  "Computer",
			coin:  ai_coin,
			isBot: true,
			cmd:   botCmd,
		}
	}
	return
}

func botCmd(matrix board, opponent byte) tea.Cmd {

	return func() tea.Msg {
		c := alphaBeta(matrix, opponent)
		return columnMsg(c)
	}
}

func alphaBeta(matrix board, opponent byte) int {

	colandscore := minimax(matrix, 5, math.MinInt, math.MaxInt, true)

	return colandscore[0]
}

func minimax(matrix board, depth, alpha, beta int, maximizingplayer bool) [2]int {

	valid_locations := matrix.getOpenColList()
	is_terminal := matrix.isTerminalNode()
	if depth == 0 || is_terminal {
		if is_terminal {
			if matrix.winingMove(ai_coin) {
				return [2]int{-1, 1000000000000000}
			} else if matrix.winingMove(player_coin) {
				return [2]int{-1, -100000000000000}
			} else {
				return [2]int{-1, 0}
			}
		} else {
			score := matrix.scorePosition(ai_coin)
			// fmt.Printf("score: %d\n", score)
			return [2]int{-1, score}
		}
	}
	if maximizingplayer {
		val := math.MinInt
		c := rand.Int()
		c = c % len(valid_locations)
		for _, col := range valid_locations {
			m_copy := matrix.copy()
			m_copy.dropCoin(col, ai_coin)
			new_score := minimax(m_copy, depth-1, alpha, beta, false)[1]
			if new_score > val {
				val = new_score
				c = col
			}
			alpha = max(alpha, val)
			if alpha >= beta {
				break
			}
		}
		return [2]int{c, val}
	} else {
		val := math.MaxInt
		c := rand.Int()
		c = c % len(valid_locations)
		for _, col := range valid_locations {
			m_copy := matrix.copy()
			m_copy.dropCoin(col, player_coin)
			new_score := minimax(m_copy, depth-1, alpha, beta, true)[1]
			if new_score < val {
				val = new_score
				c = col
			}
			beta = min(beta, val)
			if alpha >= beta {
				break
			}
		}
		return [2]int{c, val}
	}

}

type columnMsg int

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
	if gs.phase == errorPhase {
		return s + "error occured restart game\n"
	}
	s += "Its is " + string(gs.players[gs.currentPlayer].name) + " turn\n"
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
		return s + "Player " + string(gs.players[gs.currentPlayer].name) + " Won!"
	}
	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func (gs *gameState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if gs.phase == errorPhase {
		return gs, tea.Quit
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
						gs.matrix[r][gs.cursor] = gs.players[gs.currentPlayer].coin
						break
					}
				}
				if gs.matrix.is4Connected(r, gs.cursor) {
					gs.phase = postGame
				} else {
					opponent := int(gs.players[gs.currentPlayer].coin)
					gs.currentPlayer = toggle(gs.currentPlayer)
					if gs.players[gs.currentPlayer].isBot {
						return gs, gs.players[gs.currentPlayer].cmd(gs.matrix, byte(opponent))
					}
				}
			}
		}
	case columnMsg:
		if gs.phase == inGame {
			var r int
			for r = len(gs.matrix) - 1; r >= 0 && msg != -1; r-- {
				if gs.matrix[r][msg] == 0 {
					gs.matrix[r][msg] = gs.players[gs.currentPlayer].coin
					break
				}
			}
			if msg == -1 || gs.matrix.winingMove(ai_coin) {
				gs.phase = postGame
			} else {
				gs.currentPlayer = toggle(gs.currentPlayer)
			}
		}

	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return gs, nil
}

func (s *gameState) Init() tea.Cmd {
	return nil
}
