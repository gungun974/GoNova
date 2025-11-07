package form

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gungun974/gonova/internal/logger"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"golang.org/x/term"
)

type choiceSearchModel[T any] struct {
	searchInput           textinput.Model
	output                *T
	cursor                int
	choices               []Choice[T]
	filteredChoices       []Choice[T]
	header                string
	exit                  *bool
	alwaysShowFirstChoice bool
}

func AskChoiceSearch[T any](question string, choices []Choice[T]) T {
	var output T
	hasExit := false

	tprogram := tea.NewProgram(
		initialChoiceSearchModel(&output, question, choices, &hasExit, false),
	)
	if _, err := tprogram.Run(); err != nil {
		logger.MainLogger.Fatal(err)
	}

	handleExit(tprogram, hasExit)

	return output
}

func AskChoiceSearchWithAlwaysShowFirstChoice[T any](question string, choices []Choice[T]) T {
	var output T
	hasExit := false

	tprogram := tea.NewProgram(
		initialChoiceSearchModel(&output, question, choices, &hasExit, true),
	)
	if _, err := tprogram.Run(); err != nil {
		logger.MainLogger.Fatal(err)
	}

	handleExit(tprogram, hasExit)

	return output
}

func initialChoiceSearchModel[T any](
	output *T,
	header string,
	choices []Choice[T],
	hasExit *bool,
	alwaysShowFirstChoice bool,
) choiceSearchModel[T] {
	ti := textinput.New()
	ti.Prompt = "/ "
	ti.Focus()
	ti.TextStyle = ti.TextStyle.Bold(true)

	return choiceSearchModel[T]{
		searchInput:           ti,
		output:                output,
		choices:               choices,
		filteredChoices:       choices,
		header:                titleStyle.Render(header),
		exit:                  hasExit,
		alwaysShowFirstChoice: alwaysShowFirstChoice,
	}
}

func (m choiceSearchModel[T]) Init() tea.Cmd {
	return textinput.Blink
}

func (m choiceSearchModel[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			*m.exit = true
			return m, tea.Quit
		case "enter":
			if len(m.filteredChoices) != 0 {
				*m.output = m.filteredChoices[m.cursor].Value
				return m, tea.Quit
			}
		}
	}
	m.searchInput, cmd = m.searchInput.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		search := m.searchInput.Value()

		if search == "" {
			m.filteredChoices = m.choices
		} else {
			words := make([]string, len(m.choices))

			for i, choice := range m.choices {
				words[i] = strings.ToLower(choice.Name)
			}

			ranks := fuzzy.RankFindNormalized(strings.ToLower(search), words)

			sort.Sort(ranks)

			m.filteredChoices = make([]Choice[T], len(ranks))

			containFirst := false

			for i, rank := range ranks {
				if rank.OriginalIndex == 0 {
					containFirst = true
				}
				m.filteredChoices[i] = m.choices[rank.OriginalIndex]
			}

			if m.alwaysShowFirstChoice && !containFirst && len(m.choices) != 0 {
				m.filteredChoices = append(m.filteredChoices, m.choices[0])
			}
		}

		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.filteredChoices) - 1
			}
		case "down":
			if m.cursor < len(m.filteredChoices)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		}
	}

	if m.cursor >= len(m.filteredChoices) {
		m.cursor = len(m.filteredChoices) - 1
	}
	if len(m.filteredChoices) == 0 {
		m.cursor = 0
	}
	return m, cmd
}

func (m choiceSearchModel[T]) View() string {
	s := m.header + "\n"

	s += m.searchInput.View() + "\n"

	width, _, err := term.GetSize(0)
	if err == nil {
		text := fmt.Sprintf("  %d/%d ", len(m.filteredChoices), len(m.choices))

		var sb strings.Builder

		for i := 0; i < width-len(text)-1; i++ {
			sb.WriteString("─")
		}
		s += infoStyle.Render(text) + lineStyle.Render(sb.String()) + "\n"
	}

	for i, choice := range m.filteredChoices {
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
