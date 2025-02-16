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

	Line = lipgloss.AdaptiveColor{Light: light.Surface1().Hex, Dark: dark.Surface1().Hex}

	Info = lipgloss.AdaptiveColor{Light: light.Rosewater().Hex, Dark: dark.Rosewater().Hex}
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

var lineStyle = lipgloss.NewStyle().
	Foreground(Line)

var infoStyle = lipgloss.NewStyle().
	Foreground(Info)
