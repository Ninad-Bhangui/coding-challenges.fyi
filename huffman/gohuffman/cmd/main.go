package main

import (
	"bufio"
	"container/heap"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type FreqTable map[rune]int

type HuffBaseNode interface {
	IsLeaf() bool
	Priority() int
}

type HuffLeafNode struct {
	element  rune
	priority int
}

func (n *HuffLeafNode) IsLeaf() bool {
	return true
}

func (n *HuffLeafNode) Priority() int {
	return n.priority
}

type HuffInternalNode struct {
	left  HuffBaseNode
	right HuffBaseNode
}

func (n *HuffInternalNode) IsLeaf() bool {
	return false
}

func (n *HuffInternalNode) Priority() int {
	return n.left.Priority() + n.right.Priority()
}

type HuffTree struct {
	root HuffBaseNode
}

func (t HuffTree) BuildEncodingMap() map[rune]string {
	encoded_map := make(map[rune]string)
	t.walk(t.root, "", encoded_map)

	return encoded_map
}

func (t HuffTree) walk(node HuffBaseNode, path string, m map[rune]string) {
	if node.IsLeaf() {
		leaf := node.(*HuffLeafNode)
		m[leaf.element] = path
		return
	}
	internal := node.(*HuffInternalNode)
	t.walk(internal.left, path+"0", m)
	t.walk(internal.right, path+"1", m)
}

// Generic priority queue
type PriorityItem interface {
	Priority() int
}

type PriorityQueue[T PriorityItem] struct {
	items []T
}

func NewPriorityQueue[T PriorityItem]() *PriorityQueue[T] {
	return &PriorityQueue[T]{items: []T{}}
}

func (pq PriorityQueue[T]) Len() int {
	return len(pq.items)
}

func (pq PriorityQueue[T]) Less(i, j int) bool {
	return pq.items[i].Priority() < pq.items[j].Priority()
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

func (pq *PriorityQueue[T]) Push(x any) {
	pq.items = append(pq.items, x.(T))
}

func (pq *PriorityQueue[T]) Pop() any {
	n := len(pq.items)
	item := pq.items[n-1]
	pq.items = pq.items[:n-1]
	return item
}

func (pq *PriorityQueue[T]) Enqueue(item T) {
	heap.Push(pq, item)
}

func (pq *PriorityQueue[T]) Dequeue() T {
	return heap.Pop(pq).(T)
}

func (pq *PriorityQueue[T]) Peek() T {
	return pq.items[0]
}

func (pq *PriorityQueue[T]) Size() int {
	return len(pq.items)
}

func (pq *PriorityQueue[T]) Init() {
	heap.Init(pq)
}

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

	table := calculateFreq(file)
	tree := createHuffmanTree(table)
	encoded_map := tree.BuildEncodingMap()
	fmt.Println(encoded_map)
	outputFile, err := os.Create(outputpath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	writeHeader(outputFile, table)
	file.Seek(0, 0)
	writeData(file, outputFile, encoded_map)
}

func writeData(reader io.Reader, writer io.Writer, encoded_map map[rune]string) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		text := scanner.Text()
		rune := []rune(text)[0]
		bitstring := encoded_map[rune]
		for _, bit := range bitstring {
			if bit == '0' {
				writer.Write([]byte{0})
			} else {
				writer.Write([]byte{1})
			}
		}
	}
}

func calculateFreq(reader io.Reader) FreqTable {
	freq_table := FreqTable{}
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		text := scanner.Text()
		rune := []rune(text)[0]
		freq_table[rune]++
	}

	return freq_table
}

func writeHeader(w io.Writer, freqTable map[rune]int) error {
	entryCount := len(freqTable)
	err := binary.Write(w, binary.LittleEndian, int32(entryCount))
	if err != nil {
		return err
	}
	for r, freq := range freqTable {
		err := binary.Write(w, binary.LittleEndian, int32(r))
		if err != nil {
			return err
		}
		err = binary.Write(w, binary.LittleEndian, int32(freq))
		if err != nil {
			return err
		}
	}
	return nil
}

func readHeader(r io.Reader) (map[rune]int, error) {
	entryCount := int32(0)
	err := binary.Read(r, binary.LittleEndian, &entryCount)
	if err != nil {
		return nil, err
	}

	freqTable := make(map[rune]int, entryCount)
	for i := 0; i < int(entryCount); i++ {
		key := int32(0)
		err := binary.Read(r, binary.LittleEndian, &key)
		if err != nil {
			return nil, err
		}

		freq := int32(0)
		err = binary.Read(r, binary.LittleEndian, &freq)
		if err != nil {
			return nil, err
		}

		freqTable[rune(key)] = int(freq)
	}

	return freqTable, nil
}

func createHuffmanTree(freqTable FreqTable) HuffTree {
	pq := NewPriorityQueue[HuffBaseNode]()
	for key, value := range freqTable {
		node := &HuffLeafNode{key, value}
		pq.Enqueue(node)
	}

	for pq.Size() > 1 {
		left := pq.Dequeue()
		right := pq.Dequeue()
		parent := &HuffInternalNode{left, right}
		pq.Enqueue(parent)
	}

	return HuffTree{root: pq.Dequeue()}
}
