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

func TestEncodeDecodeRoundtrip(t *testing.T) {
	testText := "hello world! this is a test message with various characters 123"

	freqTable := CalculateFreq(strings.NewReader(testText))

	tree := CreateTree(freqTable)
	encodedMap := tree.BuildEncodingMap()

	var encodedBuffer strings.Builder
	err := WriteHeader(&encodedBuffer, freqTable)
	assert.NoError(t, err)

	err = WriteData(strings.NewReader(testText), &encodedBuffer, encodedMap)
	assert.NoError(t, err)

	encodedReader := strings.NewReader(encodedBuffer.String())

	decodedFreqTable, charCount, err := ReadHeader(encodedReader)
	assert.NoError(t, err)
	assert.Equal(t, freqTable, decodedFreqTable)

	decodedTree := CreateTree(decodedFreqTable)

	var decodedBuffer strings.Builder
	err = DecodeAndWriteData(encodedReader, &decodedBuffer, decodedTree, charCount)
	assert.NoError(t, err)

	decodedText := decodedBuffer.String()
	assert.Equal(t, testText, decodedText)
}
