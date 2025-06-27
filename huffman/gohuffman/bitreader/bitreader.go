package bitreader

import "io"

type BitReader struct {
	reader io.Reader
	buffer byte
	count  int
}

func NewBitReader(reader io.Reader) BitReader {
	return BitReader{
		reader: reader,
		buffer: 0,
		count:  0,
	}
}

func (br *BitReader) ReadBit() (bool, error) {
	if br.count == 0 {
		buf := []byte{0}
		n, err := br.reader.Read(buf)
		if err != nil {
			return false, err
		}
		if n == 0 {
			return false, io.EOF
		}
		br.buffer = buf[0]
		br.count = 8
	}
	bitPos := br.count - 1
	bit := (br.buffer >> bitPos) & 1
	br.count--
	return bit == 1, nil
}
