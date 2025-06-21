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
	input := `aaaee1`

	output := calculateFreq(strings.NewReader(input))
	assert.Equal(t, 2, output['e'])
	assert.Equal(t, 3, output['a'])
	assert.Equal(t, 1, output['1'])
}

func TestGetEncodingMap(t *testing.T) {
	input := `aaaee1`

	output := calculateFreq(strings.NewReader(input))
	tree := createHuffmanTree(output)
	encoded_map := tree.BuildEncodingMap()
	assert.Equal(t, 3, len(encoded_map))
	assert.Equal(t, "0", encoded_map['a'])
	assert.Equal(t, "11", encoded_map['e'])
	assert.Equal(t, "10", encoded_map['1'])
}

func TestCalculateWithFile(t *testing.T) {
	fmt.Println(os.Getwd())
	file, err := os.Open("../samples/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	output := calculateFreq(file)
	assert.Equal(t, 333, output['X'])
	assert.Equal(t, 223000, output['t'])
}
