package ui

import (
	"fmt"
	"strings"
)

func (m Model) renderFooter() string {
	pairs := []string{
		footerKey("left/right") + " tabs",
	}

	switch {
	case m.detailOpen:
		pairs = append(pairs,
			footerKey("up/down")+" scroll",
			footerKey("esc")+" back",
		)
	case m.listTab():
		pairs = append(pairs,
			footerKey("up/down")+" select",
			footerKey("enter")+" open",
		)
	default:
		pairs = append(pairs,
			footerKey("up/down")+" scroll",
		)
	}

	pairs = append(pairs, footerKey("q")+" quit")
	hint := strings.Join(pairs, strings.Repeat(" ", 4))
	if m.scrollOffset > 0 {
		hint = fmt.Sprintf("%s  scrolled +%d", hint, m.scrollOffset)
	}
	return footerStyle.Render(hint)
}

func footerKey(key string) string {
	replacer := strings.NewReplacer("left/right", "<->", "up/down", "↑↓")
	return dimStyle.Render(replacer.Replace(key))
}
