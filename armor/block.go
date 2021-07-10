package armor

import (
	"bytes"
	"hash/crc32"
	"io"
)

const (
	blockSize        = 512
	blockBufferSize  = blockSize * 2
	blockHashSize    = crc32.Size
	blockDecodeLimit = blockSize + blockHashSize
)

// Block is a minimal unit of the armored file.
type Block struct {
	contents [blockDecodeLimit]byte
	length   int
}

// NewBlock creates a valid block using the contents of the provided reader.
func NewBlock(r io.Reader) (b *Block, err error) {
	var n int
	for n, err = r.Read(b.contents[:blockSize]); ; n, err = r.Read(b.contents[b.length:blockSize]) {
		b.length += n
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
	}
	copy(b.contents[b.length:blockHashSize], Hash(b.contents[:b.length]))
	b.length += blockHashSize
	return
}

// Bytes returns the proper slice of bytes representing its contents. Must be called after the block was validated!
func (b *Block) Bytes() []byte {
	return b.contents[:b.length-blockHashSize]
}

// IsValid returns true if the hash value matches the body.
func (b *Block) IsValid() bool {
	l := b.length
	if l <= blockHashSize {
		return false
	}
	return 0 == bytes.Compare(
		b.contents[l-blockHashSize:l],
		Hash(b.contents[:l]))
}

// IsBoundary returns true if the block is small and contains many encodingBoundaryRunes. Boundary blocks are always invalid, because they do not contain a hash
func (b *Block) IsBoundary() bool {
	if b.length > encodingTelomereLength+4 {
		return false
	}
	var boundaryCount uint8
	for i := 0; i < b.length; i++ {
		if b.contents[i] == encodingBoundaryRune {
			boundaryCount++
			if boundaryCount > 6 {
				return true
			}
		}
	}
	return false
}
