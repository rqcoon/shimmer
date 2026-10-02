package cache

import (
	"sync"
)

type Page struct {
	HTML     string
	Metadata Metadata
}

type Metadata struct {
	Title       string
	Description string
	Date        string
	Tags        []string
	Template    string
}

type Cache struct {
	mu    sync.RWMutex
	pages map[string]Page
}

func New() *Cache {
	return &Cache{
		pages: make(map[string]Page),
	}
}

func (c *Cache) Get(slug string) (Page, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	page, ok := c.pages[slug]
	return page, ok
}

func (c *Cache) Set(slug string, page Page) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.pages[slug] = page
}

func (c *Cache) Delete(slug string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.pages, slug)
}
