package injector

import (
	"html/template"
	"strings"

	"github.com/rqcoon/shimmer/renderer"
)

// TODO: import page template from (custom?) file
const pageTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">

	<title>{{ .Title }}</title>

	{{ if .Description }}
	<meta name="description" content="{{ .Description }}">
	{{ end }}
</head>

<body>

	<main>
		{{ .Content }}
	</main>

</body>
</html>`

type Injector struct {
	template *template.Template
}

func New() (*Injector, error) {
	tmpl, err := template.New("Page").Parse(pageTemplate)
	if err != nil {
		return nil, err
	}

	return &Injector{
		template: tmpl,
	}, nil
}

func (i *Injector) FMInjector(page renderer.Page) (string, error) {
	data := struct {
		Title       string
		Description string
		Content     template.HTML
	}{
		Title:       page.FM.Title,
		Description: page.FM.Description,
		Content:     template.HTML(page.HTML),
	}

	var output strings.Builder

	if err := i.template.Execute(&output, data); err != nil {
		return "", err
	}

	return output.String(), nil
}
