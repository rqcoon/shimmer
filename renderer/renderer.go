package renderer

import (
	"os"
	"path/filepath"

	"github.com/rqcoon/shimmer/renderer/markdown"
)

type Page struct {
	HTML string
	FM   Frontmatter
}

type Renderer struct {
	md   *markdown.MDRenderer
	root string
}

func New(root string) *Renderer {
	return &Renderer{
		md:   markdown.NewMDRenderer(),
		root: root,
	}
}

func (r *Renderer) GetRoot() string {
	return r.root
}

func (r *Renderer) RenderSlug(slug string) (Page, error) {
	path := filepath.Join(r.root, slug+".md")

	return r.RenderFile(path)
}

func (r *Renderer) RenderFile(path string) (Page, error) {
	source, error := os.ReadFile(path)

	if error != nil {
		return Page{}, error
	}

	fm, content, err := ParseFrontmatter(string(source))
	if err != nil {
		return Page{}, err
	}

	html, err := r.md.Render(content)
	if err != nil {
		return Page{}, err
	}

	return Page{
		HTML: html,
		FM:   fm,
	}, nil
}
