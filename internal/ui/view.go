package ui

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

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

func (m Model) renderHeader() string {
	if m.width < 62 {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Xander")+" "+dimStyle.Render("@"+m.portfolio.Handle),
			eyebrowStyle.Render(m.portfolio.Role),
		)
	}

	banner := []string{
		"__  __                 _",
		"\\ \\/ /__ _ _ __   __| | ___ _ __",
		" \\  // _` | '_ \\ / _` |/ _ \\ '__|",
		" /  \\ (_| | | | | (_| |  __/ |",
		"/_/\\_\\__,_|_| |_|\\__,_|\\___|_|",
	}

	left := lipgloss.JoinVertical(
		lipgloss.Left,
		accentStyle.Render(strings.Join(banner, "\n")),
	)
	meta := eyebrowStyle.Render(m.portfolio.Role) + "  " + dimStyle.Render("@"+m.portfolio.Handle+" / "+m.portfolio.Domain)

	if m.hasPresence && m.width >= 96 {
		presenceWidth := m.width - 4 - lipgloss.Width(meta) - 4
		if presenceWidth < 18 {
			presenceWidth = 18
		}
		presence := m.renderPresenceLine(presenceWidth)
		if !m.presenceHasDetail() {
			meta = meta + "  " + presence
		} else {
			gap := m.width - 4 - lipgloss.Width(meta) - lipgloss.Width(presence)
			if gap < 2 {
				gap = 2
			}
			meta = meta + strings.Repeat(" ", gap) + presence
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, left, meta)
}

func (m Model) renderPresenceLine(width int) string {
	if !m.hasPresence {
		return ""
	}

	status := "offline"
	dot := dimStyle.Render("●")
	game := m.presence.Game()
	if game != nil {
		status = "busy"
		dot = lipgloss.NewStyle().Foreground(theme.green).Render("●")
	} else {
		switch m.presence.Status {
		case "online", "dnd":
			status = "available"
			dot = lipgloss.NewStyle().Foreground(theme.green).Render("●")
		case "idle":
			status = "inactive"
			dot = accentStyle.Render("●")
		}
	}

	var detail string
	if m.presence.ListeningToSpotify && m.presence.Spotify != nil {
		detail = "Listening to " + m.presence.Spotify.Song + " - " + m.presence.Spotify.Artist
	} else if game != nil {
		detail = "Playing " + game.Name
	} else if activity := m.presence.FirstVisibleActivity(); activity != nil {
		if activity.Type == 0 {
			detail = "Playing " + activity.Name
		} else {
			detail = activity.Name
		}
	}

	line := dot + " " + labelStyle.Render(status)
	if detail != "" {
		available := width - lipgloss.Width(stripANSI(line)) - 3
		if available < 8 {
			available = 8
		}
		line += dimStyle.Render(" / ") + dimStyle.Render(truncate(detail, available))
	}
	return line
}

func (m Model) presenceHasDetail() bool {
	if !m.hasPresence {
		return false
	}
	if m.presence.ListeningToSpotify && m.presence.Spotify != nil {
		return true
	}
	return m.presence.FirstVisibleActivity() != nil
}

func (m Model) renderTabs() string {
	names := make([]string, len(tabs))
	for i, tab := range tabs {
		name := tab.Full
		if m.width < 64 {
			name = tab.Short
		}
		if m.width < 46 {
			name = strings.ToUpper(name[:1])
		}
		if section(i) == m.active {
			names[i] = activeTabStyle.Render(name)
		} else {
			names[i] = inactiveTabStyle.Render(name)
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, names...)
}

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

func centerText(text string, width int) string {
	padding := (width - lipgloss.Width(text)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + text
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
	b.WriteString(dimStyle.Render("Everything from the site, navigable without a prompt."))
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

func indent(text, prefix string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
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
