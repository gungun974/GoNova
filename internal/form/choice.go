package form

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gungun974/gonova/internal/logger"
)

type Choice struct {
	Name  string
	Value string
}

type choiceModel struct {
	output  *string
	cursor  int
	choices []Choice
	header  string
	exit    *bool
}

func AskChoice(question string, choices []Choice) string {
	output := ""

	hasExit := false

	tprogram := tea.NewProgram(
		initialChoiceModel(&output, question, choices, &hasExit),
	)
	if _, err := tprogram.Run(); err != nil {
		logger.MainLogger.Fatal(err)
	}

	handleExit(tprogram, hasExit)

	return output
}

func initialChoiceModel(
	output *string,
	header string,
	choices []Choice,
	hasExit *bool,
) choiceModel {
	return choiceModel{
		output:  output,
		choices: choices,
		header:  titleStyle.Render(header),
		exit:    hasExit,
	}
}

func (m choiceModel) Init() tea.Cmd {
	return nil
}

func (m choiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			*m.exit = true
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.choices) - 1
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "enter":
			*m.output = m.choices[m.cursor].Value
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m choiceModel) View() string {
	s := m.header + "\n"

	for i, choice := range m.choices {
		cursor := " "
		title := choice.Name

		if m.cursor == i {
			cursor = cursorStyle.Render(">")
			title = selectedStyle.Render(choice.Name)
		}

		s += fmt.Sprintf("%s %s\n", cursor, title)
	}

	s += "\n"
	return s
}
