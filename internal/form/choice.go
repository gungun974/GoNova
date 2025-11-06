package form

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gungun974/gonova/internal/logger"
)

type Choice[T any] struct {
	Name  string
	Value T
}

type choiceModel[T any] struct {
	output  *T
	cursor  int
	choices []Choice[T]
	header  string
	exit    *bool
}

func AskChoice[T any](question string, choices []Choice[T]) T {
	var output T
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

func initialChoiceModel[T any](
	output *T,
	header string,
	choices []Choice[T],
	hasExit *bool,
) choiceModel[T] {
	return choiceModel[T]{
		output:  output,
		choices: choices,
		header:  titleStyle.Render(header),
		exit:    hasExit,
	}
}

func (m choiceModel[T]) Init() tea.Cmd {
	return nil
}

func (m choiceModel[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m choiceModel[T]) View() string {
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
