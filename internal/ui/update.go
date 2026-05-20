package ui

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.introDone {
			switch msg.String() {
			case "enter", " ":
				m.introDone = true
				return m, nil
			case "q", "ctrl+c", "esc":
				return m, tea.Quit
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab", "right", "l":
			if !m.detailOpen {
				m.nextSection()
			}

		case "shift+tab", "left", "h":
			if !m.detailOpen {
				m.previousSection()
			}

		case "up", "k":
			if m.listTab() {
				if m.cursor > 0 {
					m.cursor--
				}
				m.ensureCursorVisible()
			} else if m.scrollOffset > 0 {
				m.scrollOffset--
			}

		case "down", "j":
			if m.listTab() {
				if m.cursor < m.maxCursor() {
					m.cursor++
				}
				m.ensureCursorVisible()
			} else {
				m.scrollOffset++
			}

		case "pgup":
			m.scrollOffset -= 5
			if m.scrollOffset < 0 {
				m.scrollOffset = 0
			}

		case "pgdown":
			m.scrollOffset += 5

		case "enter":
			if m.active == projectsSection && !m.detailOpen {
				m.detailOpen = true
				m.scrollOffset = 0
			}

		case "esc", "backspace":
			if m.detailOpen {
				m.detailOpen = false
				m.scrollOffset = 0
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ensureIntroColumns()
		if m.cursor > m.maxCursor() {
			m.cursor = m.maxCursor()
		}
		m.ensureCursorVisible()

	case tickMsg:
		m.frame++
		if !m.introDone {
			m.advanceIntroRain()
		}
		return m, tick()

	case presenceTickMsg:
		return m, fetchPresence()

	case presenceMsg:
		if msg.err == nil {
			m.presence = msg.presence
			m.hasPresence = true
		}
		return m, presenceTick()
	}

	return m, nil
}
