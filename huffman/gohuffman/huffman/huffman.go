package huffman

import (
	"bufio"
	"encoding/binary"
	"io"
	"strconv"
	"strings"

	"github.com/Ninad-Bhangui/gohuffman/priorityqueue"
)

type FreqTable map[rune]int

type BaseNode interface {
	IsLeaf() bool
	Priority() int
}

type LeafNode struct {
	Element  rune
	priority int
}

func (n *LeafNode) IsLeaf() bool {
	return true
}

func (n *LeafNode) Priority() int {
	return n.priority
}

type InternalNode struct {
	Left  BaseNode
	Right BaseNode
}

func (n *InternalNode) IsLeaf() bool {
	return false
}

func (n *InternalNode) Priority() int {
	return n.Left.Priority() + n.Right.Priority()
}

type Tree struct {
	Root BaseNode
}

func (t Tree) BuildEncodingMap() map[rune]string {
	encodedMap := make(map[rune]string)
	t.walk(t.Root, "", encodedMap)
	return encodedMap
}

func (t Tree) walk(node BaseNode, path string, m map[rune]string) {
	if node.IsLeaf() {
		leaf := node.(*LeafNode)
		m[leaf.Element] = path
		return
	}
	internal := node.(*InternalNode)
	t.walk(internal.Left, path+"0", m)
	t.walk(internal.Right, path+"1", m)
}

func CalculateFreq(reader io.Reader) FreqTable {
	freqTable := FreqTable{}
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		text := scanner.Text()
		rune := []rune(text)[0]
		freqTable[rune]++
	}
	return freqTable
}

func WriteHeader(w io.Writer, freqTable map[rune]int) error {
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

func ReadHeader(r io.Reader) (map[rune]int, error) {
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

func CreateTree(freqTable FreqTable) Tree {
	pq := priorityqueue.New[BaseNode]()
	for key, value := range freqTable {
		node := &LeafNode{key, value}
		pq.Enqueue(node)
	}

	for pq.Size() > 1 {
		left := pq.Dequeue()
		right := pq.Dequeue()
		parent := &InternalNode{left, right}
		pq.Enqueue(parent)
	}

	return Tree{Root: pq.Dequeue()}
}

func WriteData(reader io.Reader, writer io.Writer, encodedMap map[rune]string) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)

	var bitBuffer strings.Builder

	// Collect all bits into one string
	for scanner.Scan() {
		text := scanner.Text()
		rune := []rune(text)[0]
		bitstring := encodedMap[rune]
		bitBuffer.WriteString(bitstring)
	}

	allBits := bitBuffer.String()

	// Process 8 bits at a time
	for i := 0; i < len(allBits); i += 8 {
		end := i + 8
		var chunk string
		if end > len(allBits) {
			// Pad remaining bits with zeros
			chunk = allBits[i:] + strings.Repeat("0", 8-(len(allBits)-i))
		} else {
			chunk = allBits[i:end]
		}

		if val, err := strconv.ParseUint(chunk, 2, 8); err == nil {
			writer.Write([]byte{byte(val)})
		}
	}
}
