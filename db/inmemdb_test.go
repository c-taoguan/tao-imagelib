package db

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	images := SearchImagesByTitle("devrev")
	fmt.Println("search by Title")
	for _, url := range images {
		fmt.Println(url)
	}

	tags := []string{"birthday"}
	images = SearchImagesByTags(tags)
	fmt.Println("search by Tags")
	for _, url := range images {
		fmt.Println(url)
	}
}
