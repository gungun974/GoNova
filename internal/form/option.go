package form

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gungun974/gonova/internal/logger"
)

type optionModel struct {
	output *bool
	header string
	left   string
	right  string
	exit   *bool
}

func AskOption(question string, defaultValue bool, left string, right string) bool {
	output := defaultValue

	hasExit := false

	tprogram := tea.NewProgram(
		initialOptionModel(&output, question, left, right, &hasExit),
	)
	if _, err := tprogram.Run(); err != nil {
		logger.MainLogger.Fatal(err)
	}

	handleExit(tprogram, hasExit)

	return output
}

func initialOptionModel(
	output *bool,
	header string,
	left string,
	right string,
	hasExit *bool,
) optionModel {
	return optionModel{
		output: output,
		left:   left,
		right:  right,
		header: titleStyle.Render(header),
		exit:   hasExit,
	}
}

func (m optionModel) Init() tea.Cmd {
	return nil
}

func (m optionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			*m.exit = true
			return m, tea.Quit
		case "left":
			*m.output = !*m.output
		case "right":
			*m.output = !*m.output
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

var selectedOptionStyle = lipgloss.
	NewStyle().
	Background(Selected).
	Foreground(Background).
	Padding(0, 2, 0)

var normalOptionStyle = lipgloss.
	NewStyle().
	Background(Background).
	Foreground(Text).
	Padding(0, 2, 0)

func (m optionModel) View() string {
	left := m.left
	right := m.right

	if *m.output {
		left = selectedOptionStyle.Render(left)
		right = normalOptionStyle.Render(right)
	} else {
		right = selectedOptionStyle.Render(right)
		left = normalOptionStyle.Render(left)
	}

	return fmt.Sprintf("%s\n\n  %s %s\n\n",
		m.header,
		left,
		right,
	)
}
