package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	if !m.introDone {
		return m.renderIntro()
	}

	if m.width < 42 || m.height < 12 {
		return appStyle.Width(m.width).Render("Xander\nresize terminal")
	}

	body := m.renderLayout()
	return m.paintViewport(body)
}

func (m Model) renderLayout() string {
	header := m.renderHeader()
	tabs := m.renderTabs()

	used := lipgloss.Height(header) + lipgloss.Height(tabs)
	bodyHeight := m.height - used - 2
	if bodyHeight < 6 {
		bodyHeight = 6
	}

	if m.width < 82 {
		return m.renderNarrowLayout(header, tabs, bodyHeight)
	}
	return m.renderWideLayout(header, tabs, bodyHeight)
}

func (m Model) renderWideLayout(header, tabs string, bodyHeight int) string {
	leftWidth := 22
	gap := 2
	contentOuterWidth := m.width - leftWidth - gap - 5
	if contentOuterWidth < 42 {
		contentOuterWidth = 42
	}

	innerWidth := contentOuterWidth - 6
	innerLines := bodyHeight - 4
	if innerLines < 1 {
		innerLines = 1
	}

	content := applyScroll(m.renderContent(innerWidth), m.scrollOffset, innerLines)
	panel := m.renderPanel(content, innerWidth, innerLines, bodyHeight)
	decor := m.renderDecor(bodyHeight, leftWidth)
	body := lipgloss.JoinHorizontal(lipgloss.Top, decor, bgSpaces(gap, bodyHeight), panel)

	layout := lipgloss.JoinVertical(lipgloss.Left, header, tabs, body)
	return lipgloss.NewStyle().PaddingLeft(2).Render(layout)
}

func (m Model) renderNarrowLayout(header, tabs string, bodyHeight int) string {
	innerWidth := m.width - 10
	if innerWidth < 24 {
		innerWidth = 24
	}

	innerLines := bodyHeight - 4
	if innerLines < 1 {
		innerLines = 1
	}

	content := applyScroll(m.renderContent(innerWidth), m.scrollOffset, innerLines)
	panel := m.renderPanel(content, innerWidth, innerLines, bodyHeight)
	layout := lipgloss.JoinVertical(lipgloss.Left, header, tabs, panel)
	return lipgloss.NewStyle().PaddingLeft(2).Render(layout)
}

func (m Model) paintViewport(body string) string {
	lines := strings.Split(body, "\n")
	bodyHeight := m.height - 1
	if bodyHeight < 1 {
		bodyHeight = 1
	}
	if len(lines) > bodyHeight {
		lines = lines[:bodyHeight]
	}
	for len(lines) < bodyHeight {
		lines = append(lines, "")
	}

	for i, line := range lines {
		lines[i] = paintLine(line, m.width)
	}

	lines = append(lines, paintLine("  "+m.renderFooter(), m.width))
	return strings.Join(lines, "\n") + "\x1b[0m"
}
