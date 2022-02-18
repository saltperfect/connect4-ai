package main

import (
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
)

func main() {
	p := tea.NewProgram(initGameType())
	out, err := p.StartReturningModel()
	if err != nil {
		log.Fatalf("unable to take game type input: %v", err)
	}

	log.Infof("returned model: %v", out)

	m := out.(*gameTypeModel)

	p = tea.NewProgram(initUserInput(m.gameType))

	out, err = p.StartReturningModel()
	if err != nil {
		log.Fatalf("unable to take user input names: %v", err)
	}
	log.Infof("returned model: %v", out)
	textInputModel := out.(*userInputModel)

	p = tea.NewProgram(initialGameState(textInputModel.inputs))
	out, err = p.StartReturningModel()
	if err != nil {
		log.Fatalf("game ended abruptly: %v", err)
	}

	// log.Infof("returned model: %v", out)

}
