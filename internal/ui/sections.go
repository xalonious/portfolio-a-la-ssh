package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderContent(width int) string {
	switch m.active {
	case aboutSection:
		return m.renderAbout(width)
	case projectsSection:
		if m.detailOpen {
			return m.renderProjectDetail(width)
		}
		return m.renderProjects(width)
	case stackSection:
		return m.renderStack(width)
	case contactSection:
		return m.renderContact(width)
	default:
		return ""
	}
}

func (m Model) renderAbout(width int) string {
	var b strings.Builder
	b.WriteString(sectionTitle("The story so far"))
	b.WriteString("\n\n")
	for _, paragraph := range m.portfolio.Story {
		b.WriteString(wrap(paragraph, width))
		b.WriteString("\n\n")
	}

	b.WriteString(sectionTitle("What I do"))
	b.WriteString("\n")
	for i, item := range m.portfolio.Focus {
		b.WriteString(fmt.Sprintf("  %s %s\n", accentStyle.Render(fmt.Sprintf("%02d", i+1)), item))
	}

	b.WriteString("\n")
	b.WriteString(sectionTitle("Primary stack"))
	b.WriteString("\n")
	b.WriteString(renderTags(m.portfolio.TechGroups[0].Items[:5], width))
	b.WriteString("\n")
	b.WriteString(renderTags([]string{"Prisma", "MySQL"}, width))
	return b.String()
}

func (m Model) renderProjects(width int) string {
	var b strings.Builder
	b.WriteString(sectionTitle("Selected work"))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("Everything I've built over the years. Each one taught me something new."))
	b.WriteString("\n\n")

	for i, project := range m.portfolio.Projects {
		selected := i == m.cursor
		prefix := "  "
		name := labelStyle.Render(project.Title)
		if selected {
			prefix = selectedStyle.Render("> ")
			name = selectedStyle.Render(project.Title)
		}

		b.WriteString(prefix + name + "\n")
		b.WriteString("  " + dimStyle.Render(truncate(project.Description, width-4)))
		b.WriteString("\n")
		b.WriteString("  " + dimStyle.Render(truncate(strings.Join(project.Tech, " / "), width-4)))
		b.WriteString("\n\n")
	}
	return b.String()
}

func (m *Model) ensureCursorVisible() {
	if !m.listTab() {
		return
	}

	width, visibleLines := m.contentMetrics()
	selectedStart, selectedEnd := m.projectLineRange(width, m.cursor)
	if selectedStart < m.scrollOffset {
		m.scrollOffset = selectedStart
		return
	}
	if selectedEnd >= m.scrollOffset+visibleLines {
		m.scrollOffset = selectedEnd - visibleLines + 1
	}
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
	m.clampScrollOffset()
}

func (m *Model) clampScrollOffset() {
	if m.width <= 0 || m.height <= 0 {
		m.scrollOffset = 0
		return
	}

	width, visibleLines := m.contentMetrics()
	totalLines := lipgloss.Height(m.renderContent(width))
	maxOffset := totalLines - visibleLines
	if maxOffset < 0 {
		maxOffset = 0
	}

	if m.scrollOffset < 0 {
		m.scrollOffset = 0
		return
	}
	if m.scrollOffset > maxOffset {
		m.scrollOffset = maxOffset
	}
}

func (m Model) contentMetrics() (int, int) {
	header := m.renderHeader()
	tabs := m.renderTabs()

	used := lipgloss.Height(header) + lipgloss.Height(tabs)
	bodyHeight := m.height - used - 2
	if bodyHeight < 6 {
		bodyHeight = 6
	}

	contentOuterWidth := m.width - 10
	if m.width >= 82 {
		contentOuterWidth = m.width - 22 - 2 - 5
		if contentOuterWidth < 42 {
			contentOuterWidth = 42
		}
	}

	innerWidth := contentOuterWidth - 6
	if innerWidth < 24 {
		innerWidth = 24
	}

	visibleLines := bodyHeight - 4
	if visibleLines < 1 {
		visibleLines = 1
	}

	return innerWidth, visibleLines
}

func (m Model) projectLineRange(width, index int) (int, int) {
	line := 3
	for i := range m.portfolio.Projects {
		blockLines := 4
		if i == index {
			return line, line + blockLines - 1
		}
		line += blockLines
	}
	return line, line
}

func (m Model) renderProjectDetail(width int) string {
	projects := m.portfolio.Projects
	if m.cursor < 0 || m.cursor >= len(projects) {
		return ""
	}

	project := projects[m.cursor]
	var b strings.Builder
	b.WriteString(dimStyle.Render("esc back to projects"))
	b.WriteString("\n\n")
	b.WriteString(sectionTitle(project.Title))
	b.WriteString("\n\n")
	b.WriteString(wrap(project.Description, width))
	b.WriteString("\n\n")
	b.WriteString(labelStyle.Render("Stack"))
	b.WriteString("\n")
	b.WriteString(renderTags(project.Tech, width))
	b.WriteString("\n\n")
	if project.Repo != "" {
		b.WriteString(labelStyle.Render("Repo"))
		b.WriteString("\n")
		b.WriteString("  " + terminalLink(project.Repo, project.Repo))
		b.WriteString("\n")
		b.WriteString(dimStyle.Render("  ctrl+click in terminals with OSC 8 support"))
	}
	return b.String()
}

func (m Model) renderStack(width int) string {
	var b strings.Builder
	b.WriteString(sectionTitle("Tools of the trade"))
	b.WriteString("\n\n")
	for _, group := range m.portfolio.TechGroups {
		b.WriteString(labelStyle.Render(group.Name))
		b.WriteString("\n")
		b.WriteString(renderTags(group.Items, width))
		b.WriteString("\n\n")
	}
	return b.String()
}

func (m Model) renderContact(width int) string {
	var b strings.Builder
	b.WriteString(sectionTitle("Let's work together"))
	b.WriteString("\n\n")
	b.WriteString(wrap("Have a project in mind, a question, or just want to say hello? These are the same contact paths from the website.", width))
	b.WriteString("\n\n")

	for _, link := range m.portfolio.Contact {
		b.WriteString(fmt.Sprintf("  %-8s %s\n\n", labelStyle.Render(link.Label+":"), terminalLink(link.URL, link.URL)))
	}
	return b.String()
}
