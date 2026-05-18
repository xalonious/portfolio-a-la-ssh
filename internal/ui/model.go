package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/xalonious/portfolio-a-la-ssh/internal/content"
)

type section int

const (
	aboutSection section = iota
	projectsSection
	stackSection
	contactSection
)

type Model struct {
	width        int
	height       int
	active       section
	cursor       int
	frame        int
	detailOpen   bool
	scrollOffset int
	portfolio    content.Portfolio
}

type tab struct {
	Full  string
	Short string
}

var tabs = []tab{
	{Full: "About", Short: "About"},
	{Full: "Projects", Short: "Work"},
	{Full: "Stack", Short: "Stack"},
	{Full: "Contact", Short: "Contact"},
}

func New(width, height int) Model {
	return Model{
		width:     width,
		height:    height,
		portfolio: content.Data,
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) listTab() bool {
	return m.active == projectsSection && !m.detailOpen
}

func (m Model) maxCursor() int {
	if m.active == projectsSection {
		return len(m.portfolio.Projects) - 1
	}
	return 0
}

func (m *Model) resetForSection() {
	m.cursor = 0
	m.detailOpen = false
	m.scrollOffset = 0
}

func (m *Model) nextSection() {
	m.active = (m.active + 1) % section(len(tabs))
	m.resetForSection()
}

func (m *Model) previousSection() {
	m.active = (m.active - 1 + section(len(tabs))) % section(len(tabs))
	m.resetForSection()
}
