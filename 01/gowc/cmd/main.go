package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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
	}
	reader = bytes.NewReader(data)
	if cliOptions.shouldGetWordCount {
		wordCount, err = getSplitCount(reader, bufio.ScanWords)
		fmt.Fprintf(output, " %d", wordCount)
	}
	reader = bytes.NewReader(data)
	if cliOptions.shouldGetCharCount {
		charCount, err = getSplitCount(reader, bufio.ScanRunes)
		fmt.Fprintf(output, " %d", charCount)
	}
	reader = bytes.NewReader(data)
	if cliOptions.shouldGetByteCount {
		byteCount, err = getSplitCount(reader, bufio.ScanBytes)
		fmt.Fprintf(output, " %d", byteCount)
	}
	if err != nil {
		return err
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
