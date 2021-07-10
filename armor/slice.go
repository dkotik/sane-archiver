package armor

import (
	"bytes"
	"fmt"
)

// SliceType provides a hint as to how the Slice should be processed.
type SliceType uint8

const (
	// SliceSizeLimit determines when the slice constructor stops retaining data. It is calculated by adding Shard size, one byte for the number of required shards, one byte for the number of redundancy shards, two bytes for the padding length (by which the recovered data is truncated), and the Hash size.
	SliceSizeLimit = blockSize + ReedSolomonMetaBinaryLength + blockHashSize
)

// Slice can contain either a shard, a shard family boundary, or a tail marker.
type Slice struct {
	Body   [SliceSizeLimit]byte // contents buffer
	Index  int                  // the cursor position where the slice began
	Length int                  // contents length
}

func (s *Slice) Write(b []byte) (n int, err error) {
	original := s.Length
	length := len(b)
	s.Length += length // track slice size any, even when discarding data
	remainingSpace := SliceSizeLimit - length

	if length > remainingSpace {
		// TODO: can be optimized by not setting boundary original+?
		copy(s.Body[original:original+remainingSpace], b[:remainingSpace])
	} else if remainingSpace > 0 {
		// TODO: can be optimized by not setting boundary original+?
		copy(s.Body[original:original+length], b[:])
	}
	return length, nil
}

// WriteChecksum adds a checksum to the body. If there is not enough space in the buffer, the Slice will overrun in length and will be marked as invalid.
func (s *Slice) WriteChecksum() {
	s.Write(Hash(s.Body[:s.Length]))
}

// IsValid returns true if the hash value matches the body.
func (s *Slice) IsValid() bool {
	length := s.Length
	if length < encodingTelomereLength+blockHashSize*2 || length > SliceSizeLimit {
		return false
	}
	return 0 == bytes.Compare(
		s.Body[length-blockHashSize:length], Hash(s.Body[:length-blockHashSize]))
}

// // Checksum returns a slice of bytes in which a checksum should be located for reading or writing.
// func (s *Slice) Checksum() []byte {
// 	if s.IsShard() {
// 		return s.Body[SliceSizeLimit-blockHashSize:] // from the very end
// 	}
// 	return s.Body[encodingTelomereLength+blockHashSize : encodingTelomereLength+blockHashSize*2]
// }
//
// // Bytes returns the proper slice of bytes representing the contents.
// func (s *Slice) Bytes() []byte {
// 	if s.IsShard() {
// 		return s.Body[:SliceSizeLimit-blockHashSize]
// 	}
// 	return s.Body[:SliceSizeLimit-blockHashSize]
// }

// IsShard returns true if the Slice likely contains a shard.
// func (s *Slice) IsShard() bool {
// 	return s.Length == SliceSizeLimit
// }

func (s *Slice) String() string {
	return fmt.Sprintf("Slice#%x@%d-%d",
		s.Body[s.Length-blockHashSize:s.Length], s.Index, s.Index+s.Length)
}
