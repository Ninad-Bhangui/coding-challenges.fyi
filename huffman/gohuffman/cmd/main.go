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
	actionarg := flag.String("action", "encode", "encode/decode")
	flag.Parse()

	if *filearg == "" || (*actionarg != "encode" && *actionarg != "decode") || *outputarg == "" {
		flag.Usage()
		os.Exit(1)
	}
	filepath := *filearg
	fmt.Println("Got filepath: ", filepath)
	outputpath := *outputarg
	fmt.Println("Got outputpath: ", outputpath)

	action := *actionarg
	fmt.Println("Got action: ", action)
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	outputFile, err := os.Create(outputpath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	if action == "encode" {
		table := huffman.CalculateFreq(file)
		tree := huffman.CreateTree(table)
		encodedMap := tree.BuildEncodingMap()
		fmt.Println(encodedMap)

		huffman.WriteHeader(outputFile, table)
		file.Seek(0, 0)
		huffman.WriteData(file, outputFile, encodedMap)
	} else {
		frequency_table, charCount, err := huffman.ReadHeader(file)
		if err != nil {
			log.Fatal(err)
		}
		tree := huffman.CreateTree(frequency_table)
		huffman.DecodeAndWriteData(file, outputFile, tree, charCount)
	}
}
