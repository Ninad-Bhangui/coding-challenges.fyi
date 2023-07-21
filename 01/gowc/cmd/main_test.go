package main

import (
	"bufio"
	"bytes"
	"fmt"
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
		if got != want {
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
		want := fmt.Sprintf("  7137  58159 341836 %s\n", fileName)
		if got != want {
			t.Errorf("cli output should be %s, got %s.", want, got)
		}
	})

}
func TestByteCountFile(t *testing.T) {
	t.Run("testing bytecount of sample file should not return error", func(t *testing.T) {
		fileNames := []string{"../samples/gutenberg.org_cache_epub_132_pg132.txt"}
		got, err := getByteCountOfFile(fileNames[0])
		want := 341836
		if err != nil {
			t.Errorf("error should be nil, got: %s", err)
		}
		if got != want {
			t.Errorf("count of bytes should be %d, got %d", want, got)
		}
	})

}

func TestLineCountFile(t *testing.T) {
	t.Run("testing linecount of sample file should return expected value", func(t *testing.T) {
		fileNames := []string{"../samples/gutenberg.org_cache_epub_132_pg132.txt"}
		got, err := getLineCountOfFile(fileNames[0])
		want := 7137
		if err != nil {
			t.Errorf("error should be nil, got: %s", err)
		}
		if got != want {
			t.Errorf("count of lines should be %d, got %d", want, got)
		}
	})
}
func TestWordCountFile(t *testing.T) {
	t.Run("testing wordcount of sample file should return expected value", func(t *testing.T) {
		fileNames := []string{"../samples/gutenberg.org_cache_epub_132_pg132.txt"}
		got, err := getWordCountOfFile(fileNames[0])
		want := 58159
		if err != nil {
			t.Errorf("error should be nil, got: %s", err)
		}
		if got != want {
			t.Errorf("count of words should be %d, got %d", want, got)
		}
	})
}
func TestCharCountFile(t *testing.T) {
	t.Run("testing char count of sample file should return expected value", func(t *testing.T) {
		fileNames := []string{"../samples/gutenberg.org_cache_epub_132_pg132.txt"}
		got, err := getCharCountOfFile(fileNames[0])
		want := 339120
		if err != nil {
			t.Errorf("error should be nil, got: %s", err)
		}
		if got != want {
			t.Errorf("count of characters should be %d, got %d", want, got)
		}
	})
}
