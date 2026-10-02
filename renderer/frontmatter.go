package renderer

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

type Frontmatter struct {
	Title       string   `yaml:"title"`
	Date        string   `yaml:"date"`
	Description string   `yaml:"description"`
	Template    string   `yaml:"template"`
	Tags        []string `yaml:"tags"`
}

func ParseFrontmatter(source string) (Frontmatter, string, error) {
	const delimiter = "---"

	if !strings.HasPrefix(source, delimiter) {
		return Frontmatter{}, source, nil
	}

	parts := strings.SplitN(source, delimiter, 3)

	if len(parts) != 3 {
		return Frontmatter{}, "", fmt.Errorf("invalid frontmatter")
	}

	var meta Frontmatter

	if err := yaml.Unmarshal([]byte(parts[1]), &meta); err != nil {
		return Frontmatter{}, "", fmt.Errorf("parsing frontmatter: %w", err)
	}

	content := strings.TrimSpace(parts[2])

	return meta, content, nil
}
