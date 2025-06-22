package huffman

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateFreq(t *testing.T) {
	input := `aaaee1`
	
	output := CalculateFreq(strings.NewReader(input))
	assert.Equal(t, 2, output['e'])
	assert.Equal(t, 3, output['a'])
	assert.Equal(t, 1, output['1'])
}

func TestCreateTreeAndBuildEncodingMap(t *testing.T) {
	input := `aaaee1`
	
	freqTable := CalculateFreq(strings.NewReader(input))
	tree := CreateTree(freqTable)
	encodedMap := tree.BuildEncodingMap()
	
	assert.Equal(t, 3, len(encodedMap))
	
	// Verify that all characters have encodings
	assert.Contains(t, encodedMap, 'a')
	assert.Contains(t, encodedMap, 'e')
	assert.Contains(t, encodedMap, '1')
	
	// Verify encoding lengths make sense (more frequent = shorter)
	assert.True(t, len(encodedMap['a']) <= len(encodedMap['e']))
	assert.True(t, len(encodedMap['a']) <= len(encodedMap['1']))
}