package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/rqcoon/shimmer/cache"
	"github.com/rqcoon/shimmer/renderer"
	"github.com/rqcoon/shimmer/renderer/injector"
)

type MDIndexEntry struct {
	Slug string   `json:"slug"`
	Name string   `json:"name"`
	Date string   `json:"date,omitempty"`
	Tags []string `json:"tags,omitempty"`
}

type Server struct {
	mux      *http.ServeMux
	cache    *cache.Cache
	renderer *renderer.Renderer
	injector *injector.Injector
}

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
	}

	s.registerRoutes()

	return s
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

// does not check for empty metadata yay
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	pages := s.cache.List()
	index := make([]MDIndexEntry, 0, len(pages))

	for slug, page := range pages {
		index = append(index, MDIndexEntry{
			Slug: slug,
			Name: page.Metadata.Title,
			Date: page.Metadata.Date,
			Tags: page.Metadata.Tags,
		})
	}

	sort.Slice(index, func(i, j int) bool {
		return index[i].Slug < index[j].Slug
	})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(index); err != nil {
		fmt.Printf("Failed to encode page index: %v\n", err)
	}
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/health", s.healthHandler)
	s.mux.HandleFunc("/index", s.indexHandler)
	s.mux.HandleFunc("GET /page/{slug}", s.handlePage)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handlePage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	partial := r.URL.Query().Get("frag")

	if partial == "" {
		s.serveFullPage(w, slug)
		return
	}

	s.serveFragment(w, slug, partial)
}

func (s *Server) serveFullPage(w http.ResponseWriter, slug string) {
	fmt.Printf("[GET] requested page: %v/%v\n", s.renderer.GetRoot(), slug)

	page, ok := s.cache.Get(slug)

	if ok {
		fmt.Println("Found page: ", slug)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(page.HTML))
		return
	}

	fmt.Printf("%v/%v: Page found and compiled\n", s.renderer.GetRoot(), slug)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(page.HTML))
}

func (s *Server) serveFragment(w http.ResponseWriter, slug string, target string) {
	page, err := s.renderer.RenderSlug(slug)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	html, err := s.injector.Inject(page, injector.RenderTarget(target))

	if err != nil {
		http.Error(w, "Fragment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	_, _ = w.Write([]byte(html))
}
