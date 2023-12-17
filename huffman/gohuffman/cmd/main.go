package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type FreqTable map[string]int

func main() {
	filearg := flag.String("filepath", "", "filepath")
	flag.Parse()

	if *filearg == "" {
		flag.Usage()
	}
	filepath := *filearg
	fmt.Println("Got filepath: ", filepath)
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	table := calculateFreq(file)
	fmt.Println(table)
}

func calculateFreq(reader io.Reader) FreqTable {
	freq_table := FreqTable{}
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		freq_table[scanner.Text()]++
	}
	return freq_table
}
