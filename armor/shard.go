package armor

// Shard is a basic component of a ReedSolomon collection, from which data can be recovered.
type Shard struct {
	Body           []byte
	SequenceNumber uint8
}
