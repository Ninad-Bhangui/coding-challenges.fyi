package bitwriter

import "io"

type BitWriter struct {
	writer io.Writer
	buffer byte
	count  int
}

func NewBitWriter(writer io.Writer) BitWriter {
	return BitWriter{
		writer: writer,
		buffer: 0,
		count:  0,
	}
}

func (bw *BitWriter) WriteBitsFromString(bits string) error {
	for _, b := range bits {
		err := bw.WriteBit(b == '1')
		if err != nil {
			return err
		}
	}
	return nil
}

func (bw *BitWriter) WriteBit(bit bool) error {
	if bit {
		bw.buffer |= 1 << (7 - bw.count)
	}
	bw.count++

	if bw.count == 8 {
		err := bw.flushByte()
		if err != nil {
			return err
		}
	}
	return nil
}

func (bw *BitWriter) flushByte() error {
	_, err := bw.writer.Write([]byte{bw.buffer})
	if err != nil {
		return err
	}
	bw.buffer = 0
	bw.count = 0
	return err
}

func (bw *BitWriter) Flush() error {
	if bw.count == 0 {
		return nil
	}
	return bw.flushByte()
}
