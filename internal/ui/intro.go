package ui

import (
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderIntro() string {
	w := max(m.width, 42)
	h := max(m.height, 12)
	chars := []rune("XANDER<>[]{}01")
	title := []string{
		"__  __                 _",
		"\\ \\/ /__ _ _ __   __| | ___ _ __",
		" \\  // _` | '_ \\ / _` |/ _ \\ '__|",
		" /  \\ (_| | | | | (_| |  __/ |",
		"/_/\\_\\__,_|_| |_|\\__,_|\\___|_|",
	}
	subtitle := "portfolio-a-la-ssh"
	prompt := "enter to continue"
	help := "q / esc to quit"

	blockHeight := len(title) + 4
	top := max(1, h/2-blockHeight/2)
	titleWidth := 0
	for _, line := range title {
		if len(line) > titleWidth {
			titleWidth = len(line)
		}
	}
	left := max(0, (w-titleWidth)/2)
	focusLeft := max(0, left-8)
	focusRight := min(w-1, left+titleWidth+8)
	focusTop := max(0, top-3)
	focusBottom := min(h-1, top+blockHeight+3)

	lines := make([]string, h)
	for y := 0; y < h; y++ {
		var b strings.Builder
		for x := 0; x < w; x++ {
			if y >= top && y < top+len(title) && x >= left && x < left+len(title[y-top]) {
				ch := rune(title[y-top][x-left])
				if ch == ' ' {
					b.WriteByte(' ')
				} else {
					b.WriteString(accentStyle.Render(string(ch)))
				}
				continue
			}

			subtitleY := top + len(title) + 1
			subtitleX := max(0, (w-len(subtitle))/2)
			if y == subtitleY && x >= subtitleX && x < subtitleX+len(subtitle) {
				b.WriteString(dimStyle.Render(string(subtitle[x-subtitleX])))
				continue
			}

			promptY := subtitleY + 2
			promptX := max(0, (w-len(prompt))/2)
			if y == promptY && x >= promptX && x < promptX+len(prompt) {
				if (m.frame/6)%2 == 0 {
					b.WriteString(labelStyle.Render(string(prompt[x-promptX])))
				} else {
					b.WriteString(accentStyle.Render(string(prompt[x-promptX])))
				}
				continue
			}

			helpY := promptY + 2
			helpX := max(0, (w-len(help))/2)
			if y == helpY && x >= helpX && x < helpX+len(help) {
				b.WriteString(dimStyle.Render(string(help[x-helpX])))
				continue
			}

			if x >= focusLeft && x <= focusRight && y >= focusTop && y <= focusBottom {
				noise := introHash(x, y, m.frame/4, 17) % 100
				if noise < 86 {
					b.WriteByte(' ')
					continue
				}
			}

			b.WriteString(m.introCell(x, y, chars))
		}
		lines[y] = paintLine(b.String(), w)
	}
	return strings.Join(lines, "\n") + "\x1b[0m"
}

func (m Model) introCell(x, y int, chars []rune) string {
	if x >= 0 && x < len(m.introColumns) {
		col := m.introColumns[x]
		if col.active && y <= col.head && col.head-y < col.trail {
			dist := col.head - y
			noise := introHash(x, y, m.frame/2, col.glyphShift)
			ch := string(chars[noise%len(chars)])
			switch {
			case dist == 0:
				return labelStyle.Render(ch)
			case dist < 3:
				return accentStyle.Render(ch)
			case dist < 8:
				return lipgloss.NewStyle().Foreground(theme.blue).Render(ch)
			default:
				return dimStyle.Render(ch)
			}
		}
	}

	noise := introHash(x, y, m.frame/8, 29)
	if noise%1000 < 5 {
		ch := string(chars[noise%len(chars)])
		return dimStyle.Render(ch)
	}
	return " "
}

func (m *Model) ensureIntroColumns() {
	w := max(m.width, 42)
	if len(m.introColumns) == w {
		return
	}
	columns := make([]introColumn, w)
	copy(columns, m.introColumns)
	for i := len(m.introColumns); i < w; i++ {
		columns[i] = m.newIntroColumn()
		columns[i].active = m.randIntro(100) < 25
	}
	m.introColumns = columns
}

func (m *Model) advanceIntroRain() {
	m.ensureIntroColumns()
	h := max(m.height, 12)
	for i := range m.introColumns {
		col := m.introColumns[i]
		if col.active {
			col.head += col.speed
			if col.head-col.trail > h {
				col.active = false
				col.cooldown = 2 + m.randIntro(max(8, h/2))
			}
			m.introColumns[i] = col
			continue
		}

		if col.cooldown > 0 {
			col.cooldown--
			m.introColumns[i] = col
			continue
		}

		if m.randIntro(100) < 16 {
			col = m.newIntroColumn()
			col.active = true
			col.head = -m.randIntro(max(3, h/2))
		} else {
			col.cooldown = 2 + m.randIntro(12)
		}
		m.introColumns[i] = col
	}
}

func (m *Model) newIntroColumn() introColumn {
	h := max(m.height, 12)
	return introColumn{
		active:     false,
		head:       -1,
		speed:      1 + m.randIntro(2),
		trail:      8 + m.randIntro(max(8, h/2)),
		cooldown:   2 + m.randIntro(max(10, h/2)),
		glyphShift: m.randIntro(1000),
	}
}

func (m *Model) randIntro(n int) int {
	if n <= 0 {
		return 0
	}
	if m.introRand == nil {
		m.introRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return m.introRand.Intn(n)
}
