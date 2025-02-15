package form

import (
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

var (
	light = catppuccin.Latte
	dark  = catppuccin.Mocha

	Background = lipgloss.AdaptiveColor{Light: light.Base().Hex, Dark: dark.Base().Hex}
	Text       = lipgloss.AdaptiveColor{Light: light.Text().Hex, Dark: dark.Text().Hex}
	Title      = lipgloss.AdaptiveColor{Light: light.Mauve().Hex, Dark: dark.Mauve().Hex}

	Selected = lipgloss.AdaptiveColor{Light: light.Pink().Hex, Dark: dark.Pink().Hex}
)

var titleStyle = lipgloss.NewStyle().
	Background(Background).
	Foreground(Title).
	Bold(true).
	Padding(0, 1, 0)

var cursorStyle = lipgloss.NewStyle().
	Foreground(Title).
	Bold(true)

var selectedStyle = lipgloss.NewStyle().
	Foreground(Selected).
	Bold(true)
