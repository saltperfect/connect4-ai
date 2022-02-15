package main

import tea "github.com/charmbracelet/bubbletea"

type gameType int

const (
	singlePlayer gameType = iota
	multiPlayer
)

type gameTypeModel struct {
	gameType gameType
}

func initGameType() *gameTypeModel {
	return &gameTypeModel{}
}

func (gs *gameTypeModel) View() string {
	s := "Play Connect 4!\n\n"
	return s + " 1. Single Player\n 2. Two Player\n"
}

func (gs *gameTypeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {

		case "1":
			gs.gameType = singlePlayer
		case "2":
			gs.gameType = multiPlayer
		case "ctrl+c", "q":
			return gs, tea.Quit

		}
		return gs, tea.Quit
	}
	return gs, nil
}

func (gs *gameTypeModel) Init() tea.Cmd {
	return nil
}
