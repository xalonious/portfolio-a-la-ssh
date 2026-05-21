package ui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

func (m Model) renderPanel(content string, width, lines, totalHeight int) string {
	if totalHeight < 5 {
		totalHeight = 5
	}
	innerWidth := width
	innerLines := totalHeight - 4
	if innerLines < 1 {
		innerLines = 1
	}

	contentLines := strings.Split(fitLines(content, innerLines), "\n")
	top := borderLine("╭", "╮", innerWidth)
	bottom := m.panelBottomBorder(innerWidth, innerLines)

	rows := make([]string, 0, totalHeight)
	rows = append(rows, top)
	rows = append(rows, panelRow("", innerWidth))
	for _, line := range contentLines {
		rows = append(rows, panelRow(line, innerWidth))
	}
	rows = append(rows, panelRow("", innerWidth))
	rows = append(rows, bottom)

	if len(rows) > totalHeight {
		rows = rows[:totalHeight-1]
		rows = append(rows, bottom)
	}
	for len(rows) < totalHeight {
		rows = append(rows[:len(rows)-1], panelRow("", innerWidth), bottom)
	}

	return strings.Join(rows, "\n")
}

func (m Model) panelBottomBorder(width, visibleLines int) string {
	totalLines := lipgloss.Height(m.renderContent(width))
	if totalLines <= visibleLines {
		return borderLine("╰", "╯", width)
	}

	maxOffset := totalLines - visibleLines
	if maxOffset < 1 {
		maxOffset = 1
	}
	position := fmt.Sprintf(" %s %d/%d %s ",
		scrollArrow(m.scrollOffset > 0, "↑"),
		min(m.scrollOffset+visibleLines, totalLines),
		totalLines,
		scrollArrow(m.scrollOffset < maxOffset, "↓"),
	)

	plainWidth := lipgloss.Width(position)
	leftFill := (width + 4 - plainWidth) / 2
	if leftFill < 0 {
		leftFill = 0
	}
	rightFill := width + 4 - plainWidth - leftFill
	if rightFill < 0 {
		rightFill = 0
	}

	border := lipgloss.NewStyle().Foreground(theme.border)
	return border.Render("╰"+strings.Repeat("─", leftFill)) +
		dimStyle.Render(position) +
		border.Render(strings.Repeat("─", rightFill)+"╯")
}

func scrollArrow(active bool, arrow string) string {
	if active {
		return arrow
	}
	return " "
}

func borderLine(left, right string, width int) string {
	return lipgloss.NewStyle().Foreground(theme.border).Render(left + strings.Repeat("─", width+4) + right)
}

func panelRow(line string, width int) string {
	line = clampLine(line, width)
	padding := width - lipgloss.Width(line)
	if padding < 0 {
		padding = 0
	}
	content := "  " + line + strings.Repeat(" ", padding) + "  "
	border := lipgloss.NewStyle().Foreground(theme.border)
	return border.Render("│") + content + border.Render("│")
}

func clampLine(line string, width int) string {
	if lipgloss.Width(line) <= width {
		return line
	}
	return truncate(stripANSI(line), width)
}

func stripANSI(text string) string {
	return ansiPattern.ReplaceAllString(text, "")
}

func paintLine(line string, width int) string {
	lineWidth := lipgloss.Width(line)
	if lineWidth > width {
		return keepBaseBackground(appStyle.Render(line))
	}
	return keepBaseBackground(appStyle.Render(line + strings.Repeat(" ", width-lineWidth)))
}

func bgSpaces(width, height int) string {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	lines := make([]string, height)
	fill := strings.Repeat(" ", width)
	for i := range lines {
		lines[i] = appStyle.Render(fill)
	}
	return strings.Join(lines, "\n")
}

func keepBaseBackground(rendered string) string {
	const base = "\x1b[48;2;11;16;32m\x1b[38;2;220;231;247m"
	return strings.ReplaceAll(rendered, "\x1b[0m", "\x1b[0m"+base)
}
