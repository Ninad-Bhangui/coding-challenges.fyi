package main

import (
	"bufio"
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
	for i := range cliOptions.fileNames {
		if cliOptions.shouldGetByteCount {
			byteCount, err := getByteCountOfFile(cliOptions.fileNames[i])
			if err != nil {
				return err
			}
			fmt.Fprintf(output, "%d %s\n", byteCount, cliOptions.fileNames[i])
			return nil
		} else if cliOptions.shouldGetLineCount {
			lineCount, err := getLineCountOfFile(cliOptions.fileNames[i])
			if err != nil {
				return err
			}
			fmt.Fprintf(output, "%d %s\n", lineCount, cliOptions.fileNames[i])
		} else if cliOptions.shouldGetWordCount {
			lineCount, err := getWordCountOfFile(cliOptions.fileNames[i])
			if err != nil {
				return err
			}
			fmt.Fprintf(output, "%d %s\n", lineCount, cliOptions.fileNames[i])
		} else if cliOptions.shouldGetCharCount {
			lineCount, err := getCharCountOfFile(cliOptions.fileNames[i])
			if err != nil {
				return err
			}
			fmt.Fprintf(output, "%d %s\n", lineCount, cliOptions.fileNames[i])
		}

	}
	return nil

}

func getByteCountOfFile(fileName string) (int, error) {
	return getSplitCountFile(fileName, bufio.ScanBytes)
}

func getLineCountOfFile(fileName string) (int, error) {
	return getSplitCountFile(fileName, bufio.ScanLines)
}

func getWordCountOfFile(fileName string) (int, error) {
	return getSplitCountFile(fileName, bufio.ScanWords)
}

func getCharCountOfFile(fileName string) (int, error) {
	return getSplitCountFile(fileName, bufio.ScanRunes)
}

func getSplitCountFile(fileName string, splitFunc bufio.SplitFunc) (int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, nil
	}
	return getSplitCount(file, splitFunc)
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
