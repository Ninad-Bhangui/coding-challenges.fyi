package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	shouldGetByteCount := flag.Bool("c", false, "Number of bytes in file or stdin")
	flag.Parse()
	nonFlagArgCount := flag.NArg()
	fileNameIndex := len(os.Args) - nonFlagArgCount
	fileNames := os.Args[fileNameIndex:]
	err := cliEntryPoint(shouldGetByteCount, fileNames)
	if err != nil {
		log.Fatal(err)
	}
}

func cliEntryPoint(shouldGetByteCount *bool, fileNames []string) error {
	for i := range fileNames {
		if *shouldGetByteCount {
			byteCount, err := getByteCountOfFile(fileNames[i])
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "%d %s\n", byteCount, fileNames[i])
			return nil
		}
	}
	return nil

}

func getByteCountOfFile(fileName string) (int, error) {
	bytes, err := os.ReadFile(fileName)
	// info, err := os.Stat(fileName)
	if err != nil {
		return 0, err
	}
	return len(bytes), nil

}
