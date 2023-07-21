package main

import (
	"testing"
)

func TestCli(t *testing.T) {
	// assertCorrectMessage := func(t testing.TB, got, want string) {
	// 	t.Helper()
	// 	if got != want {
	// 		t.Errorf("got %q want %q", got, want)
	// 	}
	// }
	t.Run("testing bytecount of sample file should not return error", func(t *testing.T) {
		shouldGetByteCount := true
		fileNames := []string{"../samples/gutenberg.org_cache_epub_132_pg132.txt"}
		got := cliEntryPoint(&shouldGetByteCount, fileNames)
		if got != nil {
			t.Errorf("error should be nil, got: %s", got)
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
