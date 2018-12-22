// Package cache provides a simple single-occupancy string cache.
// It caches stuff to RAM and disk, with a priority on using the former.
package cache

import (
	"io/ioutil"
)

// Cache is a simple single-occupancy string cache. It caches stuff
// to RAM and disk, with a priority on using the former.
type Cache struct {
	file string // filepath to data cached on disk

	data string // data cached in memory, if at all
}

// New returns a new Cache object.
func New(file string) *Cache {
	return &Cache{file: file}
}

func (c *Cache) Get() (string, bool) {
	if c.data != "" {
		return c.data, true
	}

	data, err := ioutil.ReadFile(c.file)
	if err == nil {
		c.data = string(data)
		return string(data), true
	}

	return "", false
}

func (c *Cache) Set(data string) {
	_ = ioutil.WriteFile(c.file, []byte(data), 0644)

	c.data = data
}
