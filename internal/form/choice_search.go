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

type choiceSearchModel struct {
	searchInput     textinput.Model
	output          *string
	cursor          int
	choices         []Choice
	filteredChoices []Choice
	header          string
	exit            *bool
}

func AskChoiceSearch(question string, choices []Choice) string {
	output := ""

	hasExit := false

	tprogram := tea.NewProgram(
		initialChoiceSearchModel(&output, question, choices, &hasExit),
	)
	if _, err := tprogram.Run(); err != nil {
		logger.MainLogger.Fatal(err)
	}

	handleExit(tprogram, hasExit)

	return output
}

func initialChoiceSearchModel(
	output *string,
	header string,
	choices []Choice,
	hasExit *bool,
) choiceSearchModel {
	ti := textinput.New()
	ti.Prompt = "/ "
	ti.Focus()
	ti.TextStyle = ti.TextStyle.Bold(true)

	return choiceSearchModel{
		searchInput:     ti,
		output:          output,
		choices:         choices,
		filteredChoices: choices,
		header:          titleStyle.Render(header),
		exit:            hasExit,
	}
}

func (m choiceSearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m choiceSearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			*m.exit = true
			return m, tea.Quit
		case "enter":
			*m.output = m.filteredChoices[m.cursor].Value
			return m, tea.Quit
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

			m.filteredChoices = make([]Choice, len(ranks))

			for i, rank := range ranks {
				m.filteredChoices[i] = m.choices[rank.OriginalIndex]
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

		if m.cursor >= len(m.filteredChoices) {
			m.cursor = len(m.filteredChoices) - 1
		}
	}
	return m, cmd
}

func (m choiceSearchModel) View() string {
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
