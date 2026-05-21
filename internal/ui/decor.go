package ui

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderDecor(height, width int) string {
	rows := height
	cols := width
	grid := make([][]rune, rows)
	for y := range grid {
		grid[y] = make([]rune, cols)
		for x := range grid[y] {
			grid[y][x] = ' '
		}
	}

	drawStars(grid, m.frame)
	drawCat(grid, m.frame)

	lines := make([]string, rows)
	for y, row := range grid {
		lines[y] = string(row)
	}

	return decorStyle.
		Width(cols).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(strings.Join(lines, "\n"))
}

func drawStars(grid [][]rune, frame int) {
	rows := len(grid)
	if rows == 0 {
		return
	}
	cols := len(grid[0])
	phase := float64(frame) * 0.16
	for i := 0; i < 28; i++ {
		x := int(math.Mod(float64(i*7)+phase*3, float64(cols)))
		yFloat := math.Mod(float64(i*i*5)+phase*2+math.Sin(float64(i)+phase)*2, float64(rows))
		y := int(yFloat)
		if y >= 0 && y < rows && x >= 0 && x < cols {
			if i%7 == 0 {
				grid[y][x] = '*'
			} else {
				grid[y][x] = '.'
			}
		}
	}
}

func drawCat(grid [][]rune, frame int) {
	frames := [][]string{
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

	rows := len(grid)
	if rows == 0 {
		return
	}
	cols := len(grid[0])
	cat := frames[(frame/5)%len(frames)]
	startY := rows/2 - len(cat)/2
	if startY < 2 {
		startY = 2
	}
	startX := cols/2 - lipgloss.Width(cat[0])/2
	if startX < 0 {
		startX = 0
	}

	for i, line := range cat {
		drawText(grid, startY+i, strings.Repeat(" ", startX)+line)
	}
}

func drawText(grid [][]rune, y int, text string) {
	if y < 0 || y >= len(grid) {
		return
	}
	for x, ch := range []rune(text) {
		if x >= len(grid[y]) {
			return
		}
		if ch != ' ' {
			grid[y][x] = ch
		}
	}
}
