package armor

import (
	"encoding/binary"
)

const (
	// ReedSolomonMetaBinaryLength represents the memory required to encode the Meta.
	ReedSolomonMetaBinaryLength = 1 + 1 + 1 + 2

	shardLimit = 256 // tied to Block.order and klauspost/reedsolomon limit
)

// ReedSolomonMeta holds the all the neccessary hints to perform data reconstruction.
type ReedSolomonMeta struct {
	SequenceNumber  uint8
	RequiredShards  uint8
	RedundantShards uint8
	PaddingLength   uint16 // truncate recovered data by this much
}

// Encode ReedSolomonMeta to a binary array.
func (r *ReedSolomonMeta) Encode() (b [ReedSolomonMetaBinaryLength]byte) {
	b[0] = byte(r.SequenceNumber)
	b[1] = byte(r.RequiredShards)
	b[2] = byte(r.RedundantShards)
	binary.BigEndian.PutUint16(b[3:5], r.PaddingLength)
	return
}

// ReadMeta recovers the meta from the slice.
func (s *Slice) ReadMeta() ReedSolomonMeta {
	return ReedSolomonMeta{
		SequenceNumber:  uint8(s.Body[blockSize]),
		RequiredShards:  uint8(s.Body[blockSize+1]),
		RedundantShards: uint8(s.Body[blockSize+2]),
		PaddingLength:   binary.BigEndian.Uint16(s.Body[blockSize+3 : blockSize+5]),
	}
}
