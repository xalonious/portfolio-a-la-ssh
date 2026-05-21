package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderHeader() string {
	if m.width < 62 {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Xander")+" "+dimStyle.Render("@"+m.portfolio.Handle),
			eyebrowStyle.Render(m.portfolio.Role),
		)
	}

	banner := []string{
		"__  __                 _",
		"\\ \\/ /__ _ _ __   __| | ___ _ __",
		" \\  // _` | '_ \\ / _` |/ _ \\ '__|",
		" /  \\ (_| | | | | (_| |  __/ |",
		"/_/\\_\\__,_|_| |_|\\__,_|\\___|_|",
	}

	left := lipgloss.JoinVertical(
		lipgloss.Left,
		accentStyle.Render(strings.Join(banner, "\n")),
	)
	meta := eyebrowStyle.Render(m.portfolio.Role) + "  " + dimStyle.Render("@"+m.portfolio.Handle+" / "+m.portfolio.Domain)

	if m.hasPresence && m.width >= 96 {
		presenceWidth := m.width - 4 - lipgloss.Width(meta) - 4
		if presenceWidth < 18 {
			presenceWidth = 18
		}
		presence := m.renderPresenceLine(presenceWidth)
		if !m.presenceHasDetail() {
			meta = meta + "  " + presence
		} else {
			gap := m.width - 4 - lipgloss.Width(meta) - lipgloss.Width(presence)
			if gap < 2 {
				gap = 2
			}
			meta = meta + strings.Repeat(" ", gap) + presence
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, left, meta)
}

func (m Model) renderPresenceLine(width int) string {
	if !m.hasPresence {
		return ""
	}

	status := "offline"
	dot := dimStyle.Render("●")
	game := m.presence.Game()
	if game != nil {
		status = "busy"
		dot = lipgloss.NewStyle().Foreground(theme.green).Render("●")
	} else {
		switch m.presence.Status {
		case "online", "dnd":
			status = "available"
			dot = lipgloss.NewStyle().Foreground(theme.green).Render("●")
		case "idle":
			status = "inactive"
			dot = accentStyle.Render("●")
		}
	}

	var detail string
	if m.presence.ListeningToSpotify && m.presence.Spotify != nil {
		detail = "Listening to " + m.presence.Spotify.Song + " - " + m.presence.Spotify.Artist
	} else if game != nil {
		detail = "Playing " + game.Name
	} else if activity := m.presence.FirstVisibleActivity(); activity != nil {
		if activity.Type == 0 {
			detail = "Playing " + activity.Name
		} else {
			detail = activity.Name
		}
	}

	line := dot + " " + labelStyle.Render(status)
	if detail != "" {
		available := width - lipgloss.Width(stripANSI(line)) - 3
		if available < 8 {
			available = 8
		}
		line += dimStyle.Render(" / ") + dimStyle.Render(truncate(detail, available))
	}
	return line
}

func (m Model) presenceHasDetail() bool {
	if !m.hasPresence {
		return false
	}
	if m.presence.ListeningToSpotify && m.presence.Spotify != nil {
		return true
	}
	return m.presence.FirstVisibleActivity() != nil
}

func (m Model) renderTabs() string {
	names := make([]string, len(tabs))
	for i, tab := range tabs {
		name := tab.Full
		if m.width < 64 {
			name = tab.Short
		}
		if m.width < 46 {
			name = strings.ToUpper(name[:1])
		}
		if section(i) == m.active {
			names[i] = activeTabStyle.Render(name)
		} else {
			names[i] = inactiveTabStyle.Render(name)
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, names...)
}
