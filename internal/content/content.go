// quick test

package content

type Link struct {
	Label string
	URL   string
}

type Project struct {
	Title       string
	Description string
	Repo        string
	Tech        []string
}

type TechGroup struct {
	Name  string
	Items []string
}

type Portfolio struct {
	Name       string
	Handle     string
	Domain     string
	Role       string
	Location   string
	Tagline    string
	Story      []string
	Focus      []string
	Projects   []Project
	TechGroups []TechGroup
	Contact    []Link
	InfraNotes []string
	CTONotes   []string
}

var Data = Portfolio{
	Name:     "Xander",
	Handle:   "xalonious",
	Domain:   "whoisxander.dev",
	Role:     "Full-Stack Developer & Designer",
	Location: "Belgium",
	Tagline:  "I build things for the web. Self-taught, full-stack, 7 years in and still hooked.",
	Story: []string{
		"I am a self-taught full-stack developer who went headfirst into coding and never really came back up for air.",
		"My background is in JavaScript, TypeScript, and Java. I build full-stack applications with a focus on clean architecture and interfaces that actually feel good to use.",
		"I care about the open-source community and I am always experimenting with something new. Currently available for the right opportunity.",
	},
	Focus: []string{
		"Full-Stack Web Development",
		"UI & UX Design",
		"API Design & Architecture",
		"Self-Hosted Infrastructure",
		"Open Source",
	},
	Projects: []Project{
		{
			Title:       "Serendipity Scheduling App",
			Description: "A centralized scheduling web app and API for managing staff shifts and trainings for a Roblox roleplay group.",
			Repo:        "https://github.com/xalonious/serendipity-scheduling-app",
			Tech:        []string{"TypeScript", "React", "Node.js", "Express", "Tailwind", "Prisma", "MySQL"},
		},
		{
			Title:       "Serendipity Assistant",
			Description: "A general purpose Discord bot featuring moderation tools, fun/community commands, and automation utilities.",
			Repo:        "https://github.com/xalonious/serendipity-assistant",
			Tech:        []string{"JavaScript", "Node.js", "MongoDB"},
		},
		{
			Title:       "Streaming App",
			Description: "A self-hosted media streaming web app for discovering and playing movies and TV shows from user-configured sources.",
			Repo:        "https://github.com/xalonious/streaming-app",
			Tech:        []string{"TypeScript", "React", "Node.js", "Express", "Tailwind"},
		},
		{
			Title:       "xanderGPT",
			Description: "A ChatGPT-style web app powered by a local LLM via Ollama, featuring real-time streaming responses, persistent conversations and web search.",
			Repo:        "https://github.com/xalonious/xanderGPT",
			Tech:        []string{"TypeScript", "React", "Node.js", "Express", "Tailwind", "Prisma", "MySQL"},
		},
		{
			Title:       "Barber App",
			Description: "A school project: a simple booking app where users schedule appointments with barbers. Built with React and Express for a smooth UX.",
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
			Title:	   	"My SSH Portfolio",
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
			Title:       "BongoClicker",
			Description: "A Windows autoclicker for bongocat. Fires every key on the keyboard in rapid succession while the cat jams along.",
			Repo:        "https://github.com/xalonious/bongoclicker",
			Tech:        []string{"Python", "Tkinter"},
		},
		{
			Title:       "Statuswatch",
			Description: "A self-hosted Go service that monitors third-party status pages and fires Discord webhook alerts for new incidents, updates, and resolutions.",
			Repo:        "https://github.com/xalonious/statuswatch",
			Tech:        []string{"Go"},
		},
		{
			Title:       "Image Tool",
			Description: "A simple CLI tool to convert and compress images with sensible defaults.",
			Repo:        "https://github.com/xalonious/image_tool",
			Tech:        []string{"Python"},
		},
		{
			Title:       "AutoClicker",
			Description: "A lightweight C# autoclicker to automate mouse clicks for repetitive tasks.",
			Repo:        "https://github.com/xalonious/autoclicker",
			Tech:        []string{"C#", ".NET"},
		},
		{
			Title:       "PassGuard",
			Description: "A simple password manager built with Electron and SQLite. Securely store and manage credentials locally.",
			Repo:        "https://github.com/xalonious/password-manager",
			Tech:        []string{"JavaScript", "Electron", "HTML", "CSS", "SQLite"},
		},
		{
			Title:       "Backup Code Encryptor",
			Description: "CLI tool to encrypt 2FA backup codes with a password and keep them safe from prying eyes.",
			Repo:        "https://github.com/xalonious/backup-code-encryptor",
			Tech:        []string{"Python"},
		},
	},
	TechGroups: []TechGroup{
		{Name: "Frontend", Items: []string{"React", "Next.js", "TypeScript", "TailwindCSS", "JavaScript", "CSS", "HTML", "Bootstrap", "Electron"}},
		{Name: "Backend", Items: []string{"Node.js", "Express", "Java", "C#", ".NET", "Python", "Prisma", "MySQL", "MongoDB", "SQLite"}},
		{Name: "Infrastructure", Items: []string{"Docker", "Nginx", "Linux", "Bash"}},
		{Name: "Tools", Items: []string{"Git", "GitHub", "Figma"}},
	},
	Contact: []Link{
		{Label: "Web", URL: "https://whoisxander.dev"},
		{Label: "Email", URL: "mailto:contact@whoisxander.dev"},
		{Label: "GitHub", URL: "https://github.com/xalonious"},
		{Label: "Discord", URL: "https://discordid.netlify.app/?id=531484240114876416"},
	},
	InfraNotes: []string{
		"Self-hosting shows up directly in the project list: Streaming App is built around user-configured media sources, Statuswatch monitors third-party status pages, and xanderGPT runs through a local LLM via Ollama.",
		"The web portfolio lists Docker, Nginx, Linux, and Bash as infrastructure tools.",
	},
	// TODO: add a more specific homelab/server inventory if the website grows a dedicated infrastructure section.
	CTONotes: []string{
		"Sealy is Xander's cat and CTO.",
		"Known responsibilities include sitting on the keyboard, reviewing PRs by walking across the trackpad, and general vibes management.",
		"The 404 page says Sealy deleted the page and was placed on a performance improvement plan. Treats suspended pending investigation.",
	},
}
