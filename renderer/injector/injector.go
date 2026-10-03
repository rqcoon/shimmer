package injector

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/rqcoon/shimmer/renderer"
)

type RenderTarget string

const (
	TargetPage    RenderTarget = "default.html"
	TargetContent RenderTarget = "content.html"
)

type Injector struct {
	template *template.Template
	root     string
}

type TemplateData struct {
	Title       string
	Description string
	Date        string
	Tags        []string
	Content     template.HTML
}

func New(root string) (*Injector, error) {
	templates := template.New("shimmer")

	pageFiles, err := filepath.Glob(
		filepath.Join(root, "pages", "*.html"),
	)
	if err != nil {
		return nil, fmt.Errorf("finding page templates: %w", err)
	}

	if len(pageFiles) == 0 {
		return nil, fmt.Errorf(
			"no page templates found in %s",
			filepath.Join(root, "pages"),
		)
	}

	templates, err = templates.ParseFiles(pageFiles...)
	if err != nil {
		return nil, fmt.Errorf("parsing page templates: %w", err)
	}

	fragmentFiles, err := filepath.Glob(
		filepath.Join(root, "fragments", "*.html"),
	)
	if err != nil {
		return nil, fmt.Errorf("finding template fragments: %w", err)
	}

	if len(fragmentFiles) > 0 {
		templates, err = templates.ParseFiles(fragmentFiles...)
		if err != nil {
			return nil, fmt.Errorf("parsing template fragments: %w", err)
		}
	}

	return &Injector{
		template: templates,
		root:     root,
	}, nil
}

func (i *Injector) Root() string {
	return i.root
}

// stub
func (i *Injector) Inject(page renderer.Page, target RenderTarget) (string, error) {
	data := TemplateData{
		Title:       page.FM.Title,
		Description: page.FM.Description,
		Date:        page.FM.Date,
		Tags:        page.FM.Tags,
		Content:     template.HTML(page.HTML),
	}

	var output bytes.Buffer

	if err := i.template.ExecuteTemplate(&output, string(target), data); err != nil {
		return "", fmt.Errorf("rendering template %q: %w", target, err)
	}

	return output.String(), nil
}
