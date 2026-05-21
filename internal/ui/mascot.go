package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var mascotFrames = [][]string{
	{
		`    /\___/\     `,
		`   /  o o  \    `,
		`  ( == ^ == )   `,
		`   )       (    `,
		`  (         )   `,
		` (  )     (  )  `,
		`(___)-----(___) `,
		`      \___      `,
	},
	{
		`    /\___/\     `,
		`   /  - -  \    `,
		`  ( == ^ == )   `,
		`   )       (    `,
		`  (         )   `,
		` (  )     (  )  `,
		`(___)-----(___) `,
		`       ___/     `,
	},
	{
		`    /\___/\     `,
		`   /  o o  \    `,
		`  ( == ^ == )   `,
		`   )    \  (    `,
		`  (      \  )   `,
		` (  )     (  )  `,
		`(___)-----(___) `,
		`     ___/       `,
	},
	{
		`    /\___/\     `,
		`   /  ^ ^  \    `,
		`  ( == ^ == )   `,
		`   )       (    `,
		`  (         )   `,
		` (  )     (  )  `,
		`(___)-----(___) `,
		`      \___      `,
	},
}

func mascotFrame(frame int) []string {
	return mascotFrames[(frame/5)%len(mascotFrames)]
}

func renderMascot(frame int) string {
	style := lipgloss.NewStyle().Foreground(theme.primary)
	frameLines := mascotFrame(frame)
	lines := make([]string, len(frameLines))
	for i, line := range frameLines {
		lines[i] = style.Render(line)
	}
	return strings.Join(lines, "\n")
}
