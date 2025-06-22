package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/Ninad-Bhangui/gohuffman/huffman"
)

func TestCalculateFreq(t *testing.T) {
	input := `aaaee1`

	output := huffman.CalculateFreq(strings.NewReader(input))
	assert.Equal(t, 2, output['e'])
	assert.Equal(t, 3, output['a'])
	assert.Equal(t, 1, output['1'])
}

func TestGetEncodingMap(t *testing.T) {
	input := `aaaee1`

	output := huffman.CalculateFreq(strings.NewReader(input))
	tree := huffman.CreateTree(output)
	encodedMap := tree.BuildEncodingMap()
	assert.Equal(t, 3, len(encodedMap))
	assert.Equal(t, "0", encodedMap['a'])
	assert.Equal(t, "11", encodedMap['e'])
	assert.Equal(t, "10", encodedMap['1'])
}

func TestCalculateWithFile(t *testing.T) {
	fmt.Println(os.Getwd())
	file, err := os.Open("../samples/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	output := huffman.CalculateFreq(file)
	assert.Equal(t, 333, output['X'])
	assert.Equal(t, 223000, output['t'])
}
