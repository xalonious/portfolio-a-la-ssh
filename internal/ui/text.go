package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func sectionTitle(title string) string {
	return accentStyle.Render(strings.ToUpper(title))
}

func renderTags(tags []string, width int) string {
	var lines []string
	var current string
	for _, tag := range tags {
		rendered := tagStyle.Render(tag)
		next := rendered
		if current != "" {
			next = current + " " + rendered
		}
		if lipgloss.Width(next) > width && current != "" {
			lines = append(lines, current)
			current = rendered
			continue
		}
		current = next
	}
	if current != "" {
		lines = append(lines, current)
	}
	return strings.Join(lines, "\n")
}

func terminalLink(url, label string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, linkStyle.Render(label))
}

func wrap(text string, width int) string {
	if width < 20 {
		width = 20
	}
	return lipgloss.NewStyle().Width(width).Render(text)
}

func truncate(text string, width int) string {
	if width < 4 {
		width = 4
	}
	if lipgloss.Width(text) <= width {
		return text
	}

	runes := []rune(text)
	var b strings.Builder
	for _, r := range runes {
		next := b.String() + string(r)
		if lipgloss.Width(next+"...") > width {
			break
		}
		b.WriteRune(r)
	}
	return strings.TrimRight(b.String(), " ") + "..."
}

func applyScroll(content string, offset, maxLines int) string {
	if maxLines < 1 {
		maxLines = 1
	}

	lines := strings.Split(content, "\n")
	if len(lines) <= maxLines {
		for len(lines) < maxLines {
			lines = append(lines, "")
		}
		return strings.Join(lines, "\n")
	}

	maxOffset := len(lines) - maxLines
	if offset < 0 {
		offset = 0
	}
	if offset > maxOffset {
		offset = maxOffset
	}

	return strings.Join(lines[offset:offset+maxLines], "\n")
}

func fitLines(content string, lineCount int) string {
	if lineCount < 1 {
		lineCount = 1
	}
	lines := strings.Split(content, "\n")
	if len(lines) > lineCount {
		lines = lines[:lineCount]
	}
	for len(lines) < lineCount {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func introHash(x, y, frame, salt int) int {
	v := uint32(x*73856093) ^ uint32(y*19349663) ^ uint32(frame*83492791) ^ uint32(salt*2654435761)
	v ^= v >> 13
	v *= 1274126177
	v ^= v >> 16
	return int(v & 0x7fffffff)
}
