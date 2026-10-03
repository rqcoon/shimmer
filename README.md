# Shimmer

Lightweight file-based CMS written in Go. Pages are written in markdown and rendered on server-side.

---
## Usage

Clone the repository and compile using Go 1.24+. Shimmer exists as a standalone binary for now, docker integrations are planned for future releases.

Flags:
- `-dir`: markdown content directory (default `./page`)
- `-port`: specified port to listen on (default `:8080`)
- `-template`: directory containing template(s)

File names determine the slug, e.g. 
- `./page/hello.md` -> `/page/hello`
- `./page/abc/xyz.md` -> `/page/abc/xyz`

---
## Features
- Markdown compiler
- YAML frontmatter
- Page caching
- File watcher
- Backend rendering (partial hydration for better frontend integration is WIP!)
- HTML template injection
