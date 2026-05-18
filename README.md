# portfolio-a-la-ssh

An SSH accessible terminal portfolio for [whoisxander.dev](https://whoisxander.dev).  
Built with Go, Bubble Tea, Lip Gloss, and Wish.

## Overview

portfolio-a-la-ssh runs a small SSH server that opens directly into a polished keyboard driven TUI portfolio. It is not a fake command shell, visitors can tab through sections, browse projects, open details, and find contact links from inside their terminal.

```bash
ssh term.whoisxander.dev -p 2323
```

## Features

- Wish powered SSH server on port `2323`
- Bubble Tea TUI using the alternate screen
- Responsive layout for different terminal sizes
- Keyboard navigation with tabs, arrows, enter, escape, and quit controls
- Project browser with detail views and scroll indicators
- OSC 8 terminal hyperlinks where supported
- Custom terminal styling inspired by the web portfolio
- Animated sidebar mascot

