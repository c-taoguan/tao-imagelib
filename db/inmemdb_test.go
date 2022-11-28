package db

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	images := SearchImageByTitle("devrev")
	for _, url := range images {
		fmt.Println(url)
	}
}
