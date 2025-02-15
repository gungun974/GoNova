package form

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gungun974/gonova/internal/logger"
)

type inputModel struct {
	textInput textinput.Model
	output    *string
	header    string
	exit      *bool
}

func AskInput(question string) string {
	output := ""

	hasExit := false

	tprogram := tea.NewProgram(
		initialInputModel(&output, question, "", &hasExit),
	)
	if _, err := tprogram.Run(); err != nil {
		logger.MainLogger.Fatal(err)
	}

	handleExit(tprogram, hasExit)

	return output
}

func AskInputWithPlaceholder(question string, placeholder string) string {
	output := ""

	hasExit := false

	tprogram := tea.NewProgram(
		initialInputModel(&output, question, placeholder, &hasExit),
	)
	if _, err := tprogram.Run(); err != nil {
		logger.MainLogger.Fatal(err)
	}

	handleExit(tprogram, hasExit)

	return output
}

func initialInputModel(
	output *string,
	header string,
	placeholder string,
	hasExit *bool,
) inputModel {
	ti := textinput.New()
	ti.Focus()
	ti.Placeholder = placeholder

	return inputModel{
		textInput: ti,
		output:    output,
		header:    titleStyle.Render(header),
		exit:      hasExit,
	}
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if len(m.textInput.Value()) >= 1 {
				*m.output = m.textInput.Value()
				return m, tea.Quit
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			*m.exit = true
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	return fmt.Sprintf("%s\n%s\n\n",
		m.header,
		m.textInput.View(),
	)
}
