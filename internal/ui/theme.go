package ui

import "github.com/charmbracelet/lipgloss"

var theme = struct {
	background lipgloss.Color
	foreground lipgloss.Color
	card       lipgloss.Color
	muted      lipgloss.Color
	dim        lipgloss.Color
	primary    lipgloss.Color
	border     lipgloss.Color
	green      lipgloss.Color
	blue       lipgloss.Color
}{
	background: lipgloss.Color("#0B1020"),
	foreground: lipgloss.Color("#DCE7F7"),
	card:       lipgloss.Color("#0B1020"),
	muted:      lipgloss.Color("#11182A"),
	dim:        lipgloss.Color("#758196"),
	primary:    lipgloss.Color("#8FC7E8"),
	border:     lipgloss.Color("#2A3A52"),
	green:      lipgloss.Color("#A8C97F"),
	blue:       lipgloss.Color("#B2D7F0"),
}

var (
	appStyle = lipgloss.NewStyle().
			Background(theme.background).
			Foreground(theme.foreground)

	titleStyle = lipgloss.NewStyle().
			Foreground(theme.foreground).
			Bold(true)

	eyebrowStyle = lipgloss.NewStyle().
			Foreground(theme.primary).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(theme.dim)

	labelStyle = lipgloss.NewStyle().
			Foreground(theme.foreground).
			Bold(true)

	accentStyle = lipgloss.NewStyle().
			Foreground(theme.primary).
			Bold(true)

	contentStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.border).
			Background(theme.background).
			Foreground(theme.foreground).
			Padding(1, 2)

	activeTabStyle = lipgloss.NewStyle().
			Background(theme.primary).
			Foreground(theme.background).
			Bold(true).
			Padding(0, 2).
			MarginRight(1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(theme.dim).
				Background(theme.background).
				Padding(0, 1).
				MarginRight(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(theme.primary).
			Bold(true)

	tagStyle = lipgloss.NewStyle().
			Foreground(theme.foreground).
			Background(theme.muted).
			Padding(0, 1)

	linkStyle = lipgloss.NewStyle().
			Foreground(theme.primary).
			Underline(true)

	footerStyle = lipgloss.NewStyle().
			Foreground(theme.dim)

	decorStyle = lipgloss.NewStyle().
			Foreground(theme.primary).
			Background(theme.background)
)
