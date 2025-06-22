package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Ninad-Bhangui/gohuffman/huffman"
)

func main() {
	filearg := flag.String("filepath", "", "filepath")
	outputarg := flag.String("outputpath", "", "output path")
	flag.Parse()

	if *filearg == "" {
		flag.Usage()
	}
	filepath := *filearg
	fmt.Println("Got filepath: ", filepath)
	outputpath := *outputarg
	fmt.Println("Got outputpath: ", outputpath)
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	table := huffman.CalculateFreq(file)
	tree := huffman.CreateTree(table)
	encodedMap := tree.BuildEncodingMap()
	fmt.Println(encodedMap)
	outputFile, err := os.Create(outputpath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	huffman.WriteHeader(outputFile, table)
	file.Seek(0, 0)
	huffman.WriteData(file, outputFile, encodedMap)
}