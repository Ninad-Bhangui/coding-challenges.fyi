package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateFreq(t *testing.T) {
	input := `sample text 123#`

	output := calculateFreq(strings.NewReader(input))
	assert.Equal(t, 2, output["e"])
	assert.Equal(t, 1, output["s"])
	assert.Equal(t, 1, output["1"])
	assert.Equal(t, 1, output["#"])
}

func TestCalculateWithFile(t *testing.T) {
	fmt.Println(os.Getwd())
	file, err := os.Open("../samples/gutenberg.org_cache_epub_132_pg132.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	output := calculateFreq(file)
	assert.Equal(t, 333, output["X"])
	assert.Equal(t, 223000, output["t"])
}
