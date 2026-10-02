package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rqcoon/shimmer/cache"
	"github.com/rqcoon/shimmer/renderer"
	"github.com/rqcoon/shimmer/renderer/injector"
	"github.com/rqcoon/shimmer/watcher"
)

type Server struct {
	mux      *http.ServeMux
	cache    *cache.Cache
	renderer *renderer.Renderer
	injector *injector.Injector
	watcher  *watcher.Watcher
}

const WATCHRATE time.Duration = 500

func NewServer(
	cache *cache.Cache,
	renderer *renderer.Renderer,
	injector *injector.Injector,
) *Server {
	s := &Server{
		mux:      http.NewServeMux(),
		cache:    cache,
		renderer: renderer,
		injector: injector,
		watcher: watcher.New(
			renderer.GetRoot(),
			WATCHRATE*time.Millisecond,
		),
	}

	s.registerRoutes()

	return s
}

func (s *Server) Start() error {
	fmt.Println("Compiling pages...")
	if err := s.compilePages(); err != nil {
		return err
	}

	fmt.Println("Starting file watcher")

	events := make(chan watcher.Event)
	go s.watcher.Watch(events)
	go s.watchPages(events)

	return nil
}

func (s *Server) Stop() {
	fmt.Println("Stopping file watcher")
	s.watcher.Stop()
}

func (s *Server) compilePages() error {
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

func (s *Server) compilePage(slug string) error {
	rendered, err := s.renderer.RenderSlug(slug)
	if err != nil {
		return err
	}

	html, err := s.injector.FMInjector(rendered)
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

func (s *Server) slugFromPath(path string) string {
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

func (s *Server) watchPages(events chan watcher.Event) {
	go s.watcher.Watch(events)

	for event := range events {
		s.handleFileEvent(event)
	}
}

func (s *Server) handleFileEvent(event watcher.Event) {
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

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/health", s.healthHandler)
	s.mux.HandleFunc("GET /page/{slug}", s.handlePage)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handlePage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	fmt.Printf("[GET] requested page: %v/%v\n", s.renderer.GetRoot(), slug)

	page, ok := s.cache.Get(slug)

	if ok {
		fmt.Println("Found page: ", slug)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(page.HTML))
		return
	}

	fmt.Printf("%v/%v: Missing page, attempting to compile on the fly\n", s.renderer.GetRoot(), slug)
	err := s.compilePage(slug)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		fmt.Printf("%v/%v: Page not found, aborting\n", s.renderer.GetRoot(), slug)
		return
	}

	fmt.Printf("%v/%v: Page found and compiled\n", s.renderer.GetRoot(), slug)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(page.HTML))
}
