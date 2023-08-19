package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type CliOptions struct {
	shouldGetByteCount bool
	shouldGetLineCount bool
	shouldGetWordCount bool
	shouldGetCharCount bool
	fileNames          []string
}

func (cliOptions *CliOptions) noOptionsSetDefault() {
	if !cliOptions.shouldGetByteCount && !cliOptions.shouldGetLineCount && !cliOptions.shouldGetWordCount && !cliOptions.shouldGetCharCount {
		cliOptions.shouldGetByteCount = true
		cliOptions.shouldGetLineCount = true
		cliOptions.shouldGetWordCount = true
	}
}
func main() {
	shouldGetByteCount := flag.Bool("c", false, "Number of bytes in file or stdin")
	shouldGetLineCount := flag.Bool("l", false, "Number of lines in file or stdin")
	shouldGetWordCount := flag.Bool("w", false, "Number of words in file or stdin")
	shouldGetCharCount := flag.Bool("m", false, "Number of characters in file or stdin")
	flag.Parse()

	nonFlagArgCount := flag.NArg()
	fileNameIndex := len(os.Args) - nonFlagArgCount
	fileNames := os.Args[fileNameIndex:]

	cliOptions := CliOptions{
		shouldGetByteCount: *shouldGetByteCount,
		shouldGetLineCount: *shouldGetLineCount,
		shouldGetWordCount: *shouldGetWordCount,
		shouldGetCharCount: *shouldGetCharCount,
		fileNames:          fileNames,
	}

	err := cliEntryPoint(&cliOptions, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func cliEntryPoint(cliOptions *CliOptions, output io.Writer) error {
	cliOptions.noOptionsSetDefault()
	if len(cliOptions.fileNames) == 0 {
		//Stdin mode
		cliForSingleFile(cliOptions, os.Stdin, output)
		fmt.Fprintf(output, " \n")
		return nil
	}
	for i := range cliOptions.fileNames {
		file, err := os.Open(cliOptions.fileNames[i])
		if err != nil {
			return nil
		}

		err = cliForSingleFile(cliOptions, file, output)
		fmt.Fprintf(output, " %s\n", cliOptions.fileNames[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func cliForSingleFile(cliOptions *CliOptions, input io.Reader, output io.Writer) error {
	var byteCount int
	var lineCount int
	var wordCount int
	var charCount int
	var err error
	var data []byte
	var reader io.Reader

	//TODO: BUFFERING IN MEMORY. OPTIMIZE. DOING THIS BECUASE input gets reset after one option is executed
	data, err = io.ReadAll(input)
	reader = bytes.NewReader(data)
	if cliOptions.shouldGetLineCount {
		lineCount, err = getSplitCount(reader, bufio.ScanLines)
		fmt.Fprintf(output, " %d", lineCount)
		reader = nil
	}
	reader = bytes.NewReader(data)
	if cliOptions.shouldGetWordCount {
		wordCount, err = getSplitCount(reader, bufio.ScanWords)
		fmt.Fprintf(output, " %d", wordCount)
		reader = nil
	}
	reader = bytes.NewReader(data)
	if cliOptions.shouldGetCharCount {
		charCount, err = getSplitCount(reader, bufio.ScanRunes)
		fmt.Fprintf(output, " %d", charCount)
		reader = nil
	}
	reader = bytes.NewReader(data)
	if cliOptions.shouldGetByteCount {
		byteCount, err = getSplitCount(reader, bufio.ScanBytes)
		fmt.Fprintf(output, " %d", byteCount)
		reader = nil
	}
	if err != nil {
		return err
	}
	return nil

}
func filterEmptyStrings(words []string) []string {
	finalWords := []string{}
	for _, word := range words {
		if word != "" {
			finalWords = append(finalWords, word)
		}
	}
	return finalWords

}
func getWordCountInLine(line string) int {
	count := 0
	splitLine := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' ' || r == '\t'
	})
	count += len(filterEmptyStrings(splitLine))
	return count
}

// TODO: Fix issues. This should be more performant than previous implementation. However it has issues with some special characters that take more bytes
func cliForSingleFileNew(cliOptions *CliOptions, input io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(input)
	lineCount, wordCount, charCount, byteCount := 0, 0, 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		bytes := []byte(line)
		byteCount += len(bytes)

		// bytes := scanner.Bytes()
		// byteCount += len(bytes)

		// including \n in bytecount
		// byteCount += 1

		wordCount += getWordCountInLine(line)

		charCount += len(line) + 1

		// including \n in charCount
		// charCount += 1
	}

	//Adding charCount and byteCount for newline character
	charCount += lineCount - 1
	byteCount += lineCount - 1

	if cliOptions.shouldGetLineCount {
		fmt.Fprintf(output, " %d", lineCount)
	}
	if cliOptions.shouldGetWordCount {
		fmt.Fprintf(output, " %d", wordCount)
	}
	if cliOptions.shouldGetByteCount {
		fmt.Fprintf(output, " %d", byteCount)
	}
	if cliOptions.shouldGetCharCount {
		fmt.Fprintf(output, " %d", charCount)
	}

	return nil
}
func getSplitCount(input io.Reader, splitFunc bufio.SplitFunc) (int, error) {
	count := 0
	scanner := bufio.NewScanner(input)
	scanner.Split(splitFunc)
	for scanner.Scan() {
		count++
	}
	return count, nil
}
