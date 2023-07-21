package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestCli(t *testing.T) {
	// assertCorrectMessage := func(t testing.TB, got, want string) {
	// 	t.Helper()
	// 	if got != want {
	// 		t.Errorf("got %q want %q", got, want)
	// 	}
	// }
	t.Run("testing bytecount of sample file return same output as wc -c fileName", func(t *testing.T) {
		fileName := "../samples/gutenberg.org_cache_epub_132_pg132.txt"
		fileNames := []string{fileName}
		cliOptions := CliOptions{
			shouldGetByteCount: true,
			shouldGetLineCount: false,
			shouldGetWordCount: false,
			shouldGetCharCount: false,
			fileNames:          fileNames,
		}
		var buffer bytes.Buffer
		testWriter := bufio.NewWriter(&buffer)
		error := cliEntryPoint(&cliOptions, testWriter)
		if error != nil {
			t.Errorf("error should be nil, got: %s", error)
		}
		testWriter.Flush()
		got := buffer.String()
		want := fmt.Sprintf("341836 %s\n", fileName)
		//TODO: space seperation does not match exactly with the wc command
		if strings.Trim(got, " ") != want {
			t.Errorf("cli output should be %s, got %s.", want, got)
		}
	})
	t.Run("testing linecount of sample file return same output as wc -l fileName", func(t *testing.T) {
		fileName := "../samples/gutenberg.org_cache_epub_132_pg132.txt"
		fileNames := []string{fileName}
		cliOptions := CliOptions{
			shouldGetByteCount: false,
			shouldGetLineCount: true,
			shouldGetWordCount: false,
			shouldGetCharCount: false,
			fileNames:          fileNames,
		}
		var buffer bytes.Buffer
		testWriter := bufio.NewWriter(&buffer)
		error := cliEntryPoint(&cliOptions, testWriter)
		if error != nil {
			t.Errorf("error should be nil, got: %s", error)
		}
		testWriter.Flush()
		got := buffer.String()
		want := fmt.Sprintf("7137 %s\n", fileName)
		//TODO: space seperation does not match exactly with the wc command
		if strings.Trim(got, " ") != want {
			t.Errorf("cli output should be %s, got %s.", want, got)
		}
	})
	t.Run("testing wordcount of sample file return same output as wc -w fileName", func(t *testing.T) {
		fileName := "../samples/gutenberg.org_cache_epub_132_pg132.txt"
		fileNames := []string{fileName}
		cliOptions := CliOptions{
			shouldGetByteCount: false,
			shouldGetLineCount: false,
			shouldGetWordCount: true,
			shouldGetCharCount: false,
			fileNames:          fileNames,
		}
		var buffer bytes.Buffer
		testWriter := bufio.NewWriter(&buffer)
		error := cliEntryPoint(&cliOptions, testWriter)
		if error != nil {
			t.Errorf("error should be nil, got: %s", error)
		}
		testWriter.Flush()
		got := buffer.String()
		want := fmt.Sprintf("58159 %s\n", fileName)
		//TODO: space seperation does not match exactly with the wc command
		if strings.Trim(got, " ") != want {
			t.Errorf("cli output should be %s, got %s.", want, got)
		}
	})
	t.Run("testing charcount of sample file return same output as wc -m fileName", func(t *testing.T) {
		fileName := "../samples/gutenberg.org_cache_epub_132_pg132.txt"
		fileNames := []string{fileName}
		cliOptions := CliOptions{
			shouldGetByteCount: false,
			shouldGetLineCount: false,
			shouldGetWordCount: false,
			shouldGetCharCount: true,
			fileNames:          fileNames,
		}
		var buffer bytes.Buffer
		testWriter := bufio.NewWriter(&buffer)
		error := cliEntryPoint(&cliOptions, testWriter)
		if error != nil {
			t.Errorf("error should be nil, got: %s", error)
		}
		testWriter.Flush()
		got := buffer.String()
		want := fmt.Sprintf("339120 %s\n", fileName)
		//TODO: space seperation does not match exactly with the wc command
		if strings.Trim(got, " ") != want {
			t.Errorf("cli output should be %s, got %s.", want, got)
		}
	})

	t.Run("testing no options of sample file should behave as -c -l and -w", func(t *testing.T) {
		fileName := "../samples/gutenberg.org_cache_epub_132_pg132.txt"
		fileNames := []string{fileName}
		cliOptions := CliOptions{
			shouldGetByteCount: false,
			shouldGetLineCount: false,
			shouldGetWordCount: false,
			shouldGetCharCount: false,
			fileNames:          fileNames,
		}
		var buffer bytes.Buffer
		testWriter := bufio.NewWriter(&buffer)
		error := cliEntryPoint(&cliOptions, testWriter)
		if error != nil {
			t.Errorf("error should be nil, got: %s", error)
		}
		testWriter.Flush()
		got := buffer.String()
		want := fmt.Sprintf("  7137 58159 341836 %s\n", fileName)
		//TODO: space seperation does not match exactly with the wc command
		if strings.Trim(got, " ") != strings.Trim(want, " ") {
			t.Errorf("cli output should be %s, got %s.", want, got)
		}
	})
	t.Run("testing bytecount of stdin return same output as wc -c fileName", func(t *testing.T) {
		fileName := "../samples/gutenberg.org_cache_epub_132_pg132.txt"
		fileNames := []string{}
		cliOptions := CliOptions{
			shouldGetByteCount: true,
			shouldGetLineCount: false,
			shouldGetWordCount: false,
			shouldGetCharCount: false,
			fileNames:          fileNames,
		}
		var buffer bytes.Buffer
		testWriter := bufio.NewWriter(&buffer)
		file, _ := os.Open(fileName)
		error := cliForSingleFile(&cliOptions, file, testWriter)
		if error != nil {
			t.Errorf("error should be nil, got: %s", error)
		}
		testWriter.Flush()
		got := buffer.String()
		want := "341836"
		//TODO: space seperation does not match exactly with the wc command
		if strings.Trim(got, " ") != strings.Trim(want, " ") {
			t.Errorf("cli output should be %s, got %s.", want, got)
		}
	})

}
