package content

type Link struct {
	Label string
	URL   string
}

type Project struct {
	Title       string
	Description string
	Repo        string
	CaseStudy   string
	Tech        []string
}

type TechGroup struct {
	Name  string
	Items []string
}

type Portfolio struct {
	Handle     string
	Domain     string
	Role       string
	Story      []string
	Focus      []string
	Projects   []Project
	TechGroups []TechGroup
	Contact    []Link
}

var Data = Portfolio{
	Handle: "xalonious",
	Domain: "whoisxander.dev",
	Role:   "Full-Stack Developer & Designer",
	Story: []string{
		"I'm a self-taught full-stack developer who builds modern web applications end to end, from interfaces that feel good to use to backends that don't fall over.",
		"TypeScript is at the center of most of my work, especially with React and Node.js. But I'm not married to one stack. I pick the language and tools that actually fit the problem, whether that's a web app, a CLI tool, or something a client didn't even know they needed.",
		"I'm also a little obsessed with automation. If something is repetitive, slow, or error-prone, odds are I've already turned it into a script, integration, or workflow that runs itself while I do something more interesting.",
	},
	Focus: []string{
		"Full-Stack Web Development",
		"UI & UX Design",
		"API Design & Architecture",
		"Self-Hosted Infrastructure",
		"Automation & Integrations",
	},
	TechGroups: []TechGroup{
		{Name: "Frontend", Items: []string{"React", "Next.js", "TypeScript", "TailwindCSS", "JavaScript", "CSS", "HTML", "Electron"}},
		{Name: "Backend", Items: []string{"Node.js", "Express", "Java", "C#", ".NET", "Python", "Go", "Prisma", "MySQL", "MongoDB"}},
		{Name: "Infrastructure", Items: []string{"Docker", "Nginx", "Linux", "Bash"}},
		{Name: "Tools", Items: []string{"Git", "GitHub"}},
	},
	Contact: []Link{
		{Label: "Web", URL: "https://whoisxander.dev"},
		{Label: "Email", URL: "mailto:contact@whoisxander.dev"},
		{Label: "GitHub", URL: "https://github.com/xalonious"},
		{Label: "Discord", URL: "https://discordid.netlify.app/?id=531484240114876416"},
	},
}
