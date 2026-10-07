package shimmer

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rqcoon/shimmer/api"
	"github.com/rqcoon/shimmer/cache"
	"github.com/rqcoon/shimmer/renderer"
	"github.com/rqcoon/shimmer/renderer/injector"
	"github.com/rqcoon/shimmer/watcher"
)

type Config struct {
	ContentDir  string
	TemplateDir string
	WatchRate   time.Duration
}

type Shimmer struct {
	cache    *cache.Cache
	renderer *renderer.Renderer
	injector *injector.Injector
	watcher  *watcher.Watcher
}

func New(config Config) (*Shimmer, error) {
	if config.ContentDir == "" {
		return nil, fmt.Errorf("content directory is required")
	}

	if config.TemplateDir == "" {
		return nil, fmt.Errorf("template directory is required")
	}

	if config.WatchRate == 0 {
		config.WatchRate = 500 * time.Millisecond
	}

	c := cache.New()

	r := renderer.New(config.ContentDir)

	i, err := injector.New(config.TemplateDir)
	if err != nil {
		return nil, fmt.Errorf("creating injector: %w", err)
	}

	w := watcher.New(
		r.GetRoot(),
		config.WatchRate,
	)

	return &Shimmer{
		cache:    c,
		renderer: r,
		injector: i,
		watcher:  w,
	}, nil
}

func (s *Shimmer) Start() error {
	if err := s.compilePages(); err != nil {
		return fmt.Errorf("compiling pages: %w", err)
	}

	events := make(chan watcher.Event)

	go s.watcher.Watch(events)
	go s.watchPages(events)

	return nil
}

func (s *Shimmer) Stop() {
	s.watcher.Stop()
}

func (s *Shimmer) Handler() http.Handler {
	return api.NewServer(
		s.cache,
		s.renderer,
		s.injector,
	)
}

func (s *Shimmer) compilePages() error {
	root := s.renderer.GetRoot()

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".md" {
			return nil
		}

		slug := s.slugFromPath(path)

		fmt.Printf("Compiling: %s\n", slug)

		if err := s.compilePage(slug); err != nil {
			return fmt.Errorf("compile %s: %w", slug, err)
		}

		return nil
	})
}

func (s *Shimmer) compilePage(slug string) error {
	rendered, err := s.renderer.RenderSlug(slug)
	if err != nil {
		return err
	}

	html, err := s.injector.Inject(rendered, "default")
	if err != nil {
		return err
	}

	page := cache.Page{
		HTML: html,
		Metadata: cache.Metadata{
			Title:       rendered.FM.Title,
			Description: rendered.FM.Description,
			Date:        rendered.FM.Date,
			Tags:        rendered.FM.Tags,
			Template:    rendered.FM.Template,
		},
	}

	s.cache.Set(slug, page)

	fmt.Printf("Compiled: %s\n", slug)

	return nil
}

func (s *Shimmer) slugFromPath(path string) string {
	root := s.renderer.GetRoot()

	relative, err := filepath.Rel(root, path)
	if err != nil {
		return ""
	}

	return strings.TrimSuffix(
		relative,
		filepath.Ext(relative),
	)
}

func (s *Shimmer) watchPages(events chan watcher.Event) {
	for event := range events {
		s.handleFileEvent(event)
	}
}

func (s *Shimmer) handleFileEvent(event watcher.Event) {
	slug := s.slugFromPath(event.Path)

	switch event.Type {
	case watcher.Created:
		fmt.Printf("Page created: %s\n", slug)

		if err := s.compilePage(slug); err != nil {
			fmt.Printf("Failed to compile %s: %v\n", slug, err)
		}

	case watcher.Modified:
		fmt.Printf("Page modified: %s\n", slug)

		if err := s.compilePage(slug); err != nil {
			fmt.Printf("Failed to recompile %s: %v\n", slug, err)
		}

	case watcher.Deleted:
		fmt.Printf("Page deleted: %s\n", slug)

		s.cache.Delete(slug)
	}
}
