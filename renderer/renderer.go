package renderer

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"

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

	if fm.Title == "" {
		fm.Title = titleHelper(path)
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

func titleHelper(path string) string {
	fname := filepath.Base(path)
	fname = strings.TrimSuffix(fname, filepath.Ext(fname))

	fname = strings.NewReplacer("-", " ", "_", " ").Replace(fname)
	fname = strings.Join(strings.Fields(fname), " ")

	words := strings.Fields(fname)
	for i, word := range words {
		runes := []rune(word)
		if len(runes) > 0 {
			runes[0] = unicode.ToTitle(runes[0])
			words[i] = string(runes)
		}
	}

	return strings.Join(words, " ")
}
