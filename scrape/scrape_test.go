package scrape_test

import (
	"context"
	"testing"

	"github.com/jtagcat/util/scrape"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	s := scrape.InitScraper(context.Background(), &scrape.Scraper{})
	nodes, _, err := s.Get("https://www.c7.ee/", "document")
	assert.Nil(t, err)
	print(nodes)
	panic("test not implemented")
}
