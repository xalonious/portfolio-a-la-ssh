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
	Projects: []Project{
		{
			Title:       "Serendipity Scheduling App",
			Description: "A centralized scheduling web app and API for managing staff shifts and trainings for a Roblox roleplay group.",
			Repo:        "https://github.com/xalonious/serendipity-scheduling-app",
			CaseStudy:   "https://whoisxander.dev/projects/serendipity-scheduling-app",
			Tech:        []string{"TypeScript", "React", "Node.js", "Express", "Tailwind", "Prisma", "MySQL"},
		},
		{
			Title:       "Streaming App",
			Description: "A self-hosted media streaming web app for discovering and playing movies and TV shows from user-configured sources.",
			Repo:        "https://github.com/xalonious/streaming-app",
			CaseStudy:   "https://whoisxander.dev/projects/streaming-app",
			Tech:        []string{"TypeScript", "React", "Node.js", "Express", "Tailwind"},
		},
		{
			Title:       "xanderGPT",
			Description: "A self-hosted ChatGPT-style app powered by Qwen3 through Ollama, with streamed reasoning, persistent conversations, and web-aware tool orchestration.",
			Repo:        "https://github.com/xalonious/xanderGPT",
			CaseStudy:   "https://whoisxander.dev/projects/xandergpt",
			Tech:        []string{"TypeScript", "React", "Node.js", "Express", "Tailwind", "Prisma", "MySQL"},
		},
		{
			Title:       "Bridgely",
			Description: "A Discord-to-Roblox verification bot with profile and game-based verification, automatic group role and nickname synchronization, and configurable rank, badge, and game-pass role binds.",
			Repo:        "https://github.com/xalonious/Bridgely",
			Tech:        []string{"JavaScript", "Node.js", "Discord.js", "Express", "MongoDB", "Luau"},
		},
		{
			Title:       "Barber App",
			Description: "A full-stack appointment booking app where users can schedule barber appointments through a clean, responsive interface.",
			Repo:        "https://github.com/xalonious/barber-app",
			Tech:        []string{"TypeScript", "React", "Node.js", "Bootstrap", "Express", "MySQL"},
		},
		{
			Title:       "My Portfolio Website",
			Description: "The site this SSH portfolio is based on. Built with Next.js and Tailwind CSS to showcase projects, skills, and the developer journey.",
			Repo:        "https://github.com/xalonious/portfolio",
			Tech:        []string{"TypeScript", "React", "Next.js", "Tailwind CSS", "Framer Motion"},
		},
		{
			Title:       "My SSH Portfolio",
			Description: "The portfolio you are currently on. built as a keyboard driven TUI with Go, Bubble Tea, Lip Gloss, and Wish",
			Repo:        "https://github.com/xalonious/portfolio-a-la-ssh",
			Tech:        []string{"Go"},
		},
		{
			Title:       "Robux Spent Calculator",
			Description: "An Electron desktop app that tracks Robux inflow, outflow, and current balance, with charts and spending insights.",
			Repo:        "https://github.com/xalonious/robux-spent",
			Tech:        []string{"JavaScript", "Node.js", "Electron", "HTML", "CSS"},
		},
		{
			Title:       "Statuswatch",
			Description: "A self-hosted Go service that monitors third-party status pages and fires Discord webhook alerts for new incidents, updates, and resolutions.",
			Repo:        "https://github.com/xalonious/statuswatch",
			Tech:        []string{"Go"},
		},
		{
			Title:       "Media Tool",
			Description: "A cross-platform Python CLI for converting, compressing, and precisely cutting images, audio, and video, with batch procesing and a project-local FFmpeg runtime.",
			Repo:        "https://github.com/xalonious/media_tool",
			Tech:        []string{"Python", "FFmpeg", "Pillow"},
		},
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
