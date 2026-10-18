package markdown

import (
	"net/url"
	"path"
	"strings"
)

func resolveImageURL(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || raw == "" {
		return ""
	}

	if u.Scheme == "https" || u.Scheme == "http" { //todo: warn http as unsafe
		if u.Host == "" {
			return ""
		}
		return u.String()
	}

	if u.Scheme != "" || u.Host != "" {
		return ""
	}

	// Local image paths must be relative and must not traverse outside asset directory
	if strings.HasPrefix(u.Path, "/") {
		return ""
	}

	for _, segment := range strings.Split(u.Path, "/") {
		if segment == ".." {
			return ""
		}
	}

	clean := path.Clean(u.Path)
	if clean == "." || clean == "" {
		return ""
	}

	return "/assets/" + clean
}
