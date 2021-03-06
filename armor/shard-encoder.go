package armor

import (
	"encoding/binary"
	"hash"
	"io"
)

const (
	// ShardLimit is constrained by klauspost/reedsolomon limit.
	ShardLimit = 256
)

type ShardEncoder struct {
	w   io.Writer
	c   hash.Hash32
	tag [shardTagSize]byte
}

func NewShardEncoder(w io.Writer, checkSum hash.Hash32, prefill *ShardTag) (s *ShardEncoder) {
	s = &ShardEncoder{
		w: w,
		c: checkSum,
	}

	if prefill == nil {
		prefill = &ShardTag{}
	}
	if prefill.Version == 0 {
		prefill.Version = Version
	}
	prefill.Write(s.tag[:])
	return
}

func (s *ShardEncoder) Write(b []byte) (n int, err error) {
	n, err = s.w.Write(b)
	if n > 0 {
		s.c.Sum(b[:n])
	}
	return
}

// Seal writes tag and checksum.
func (s *ShardEncoder) Seal() (err error) {
	_, err = s.Write(s.tag[:])
	if err != nil {
		return
	}

	var checkSumBytes [4]byte
	binary.BigEndian.PutUint32(checkSumBytes[:], s.c.Sum32())
	s.c.Reset()
	_, err = s.w.Write(checkSumBytes[:])
	return
}

// SetBlockSequence updates tag with a new block sequence.
func (s *ShardEncoder) SetBlockSequence(n uint64) {
	binary.BigEndian.PutUint64(
		s.tag[shardTagBlockSequencePosition:shardTagShardSequencePosition], n)
}

// SetShardSequence updates tag with a new shard sequence.
func (s *ShardEncoder) SetShardSequence(n uint8) {
	s.tag[shardTagShardSequencePosition] = byte(n)
}
