package ui

import (
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderIntro() string {
	w := max(m.width, 42)
	h := max(m.height, 14)

	if w < 58 || h < 18 {
		return m.renderCompactIntro(w, h)
	}

	title := introTitle(m.frame)
	cat := renderMascot(m.frame)
	meta := dimStyle.Render("@xalonious / whoisxander.dev")
	prompt := introPrompt(m.frame)

	placements := make([][]introPlacement, h)
	blockHeight := lipgloss.Height(title) + lipgloss.Height(cat) + 5
	top := max(1, h/2-blockHeight/2)

	lines := make([]string, h)
	paintIntroConstellation(placements, w, h, m.frame)
	putIntroBlock(placements, w, top, title)
	putIntroBlock(placements, w, top+lipgloss.Height(title), meta)
	putIntroBlock(placements, w, top+lipgloss.Height(title)+2, cat)
	putIntroBlock(placements, w, top+lipgloss.Height(title)+lipgloss.Height(cat)+4, prompt)

	for i := range lines {
		lines[i] = paintLine(renderIntroLine(w, placements[i]), w)
	}
	return strings.Join(lines, "\n") + "\x1b[0m"
}

func (m Model) renderCompactIntro(width, height int) string {
	lines := make([]string, height)
	placements := make([][]introPlacement, height)

	title := introCompactTitle(m.frame)
	cat := renderMascot(m.frame)
	prompt := introPrompt(m.frame)

	top := max(1, height/2-4)
	putIntroBlock(placements, width, top, title)
	putIntroBlock(placements, width, top+1, cat)
	putIntroBlock(placements, width, top+8, prompt)

	for i := range lines {
		lines[i] = paintLine(renderIntroLine(width, placements[i]), width)
	}
	return strings.Join(lines, "\n") + "\x1b[0m"
}

type introPlacement struct {
	x    int
	text string
}

func introTitle(frame int) string {
	title := []string{
		"X   X  AAA  N   N DDDD  EEEEE RRRR ",
		" X X  A   A NN  N D   D E     R   R",
		"  X   AAAAA N N N D   D EEEE  RRRR ",
		" X X  A   A N  NN D   D E     R R  ",
		"X   X A   A N   N DDDD  EEEEE R  RR",
	}

	var b strings.Builder
	glowCol := frame % len(title[0])
	for y, line := range title {
		if y > 0 {
			b.WriteByte('\n')
		}
		for x, ch := range line {
			if ch == ' ' {
				b.WriteByte(' ')
				continue
			}
			style := accentStyle
			if abs(x-glowCol) <= 1 {
				style = labelStyle
			}
			b.WriteString(style.Render(string(ch)))
		}
	}
	return b.String()
}

func introCompactTitle(frame int) string {
	if (frame/8)%2 == 0 {
		return accentStyle.Render("Xander")
	}
	return labelStyle.Render("Xander")
}

func introPrompt(frame int) string {
	key := labelStyle.Render("enter")
	text := dimStyle.Render(" open portfolio")
	if (frame/7)%2 == 1 {
		key = accentStyle.Render("enter")
	}
	return key + text
}

func paintIntroConstellation(lines [][]introPlacement, width, height, frame int) {
	centerX := width / 2
	centerY := height / 2
	size := min(width/4, height/2)
	if size < 6 {
		size = 6
	}

	for i := -size; i <= size; i += 2 {
		for _, x := range []int{centerX + i, centerX - i} {
			y := centerY + i/2
			if x < 1 || x >= width-1 || y < 1 || y >= height-1 {
				continue
			}
			if introHash(i, frame/6, x, 31)%100 > 34 {
				continue
			}
			ch := "."
			style := dimStyle
			if introHash(i, frame/4, x, 47)%7 == 0 {
				ch = "*"
				style = accentStyle
			}
			lines[y] = append(lines[y], introPlacement{x: x, text: style.Render(ch)})
		}
	}

	for i := 0; i < 12; i++ {
		x := introHash(i, frame/9, 0, 11) % width
		y := introHash(i, frame/11, 0, 23) % height
		if y >= len(lines) || introHash(i, frame/5, 0, 59)%100 > 42 {
			continue
		}

		ch := "."
		if (i+frame/4)%8 == 0 {
			ch = "*"
		}
		style := dimStyle
		if (i+frame/3)%6 == 0 {
			style = accentStyle
		}
		lines[y] = append(lines[y], introPlacement{x: x, text: style.Render(ch)})
	}
}

func putIntroBlock(lines [][]introPlacement, width, y int, block string) {
	if y < 0 || y >= len(lines) {
		return
	}
	for i, line := range strings.Split(block, "\n") {
		row := y + i
		if row < 0 || row >= len(lines) {
			continue
		}
		x := max(0, (width-lipgloss.Width(line))/2)
		lines[row] = append(lines[row], introPlacement{x: x, text: line})
	}
}

func renderIntroLine(width int, placements []introPlacement) string {
	if len(placements) == 0 {
		return ""
	}

	sort.SliceStable(placements, func(i, j int) bool {
		return placements[i].x < placements[j].x
	})

	var b strings.Builder
	cursor := 0
	for _, placement := range placements {
		if placement.x < cursor {
			continue
		}
		b.WriteString(strings.Repeat(" ", placement.x-cursor))
		b.WriteString(placement.text)
		cursor = placement.x + lipgloss.Width(placement.text)
		if cursor >= width {
			break
		}
	}
	return b.String()
}
