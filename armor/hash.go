package armor

import (
	"hash/crc32"
)

// Official documentation says that Koopman is superior for error detection.
var tablePolynomial = crc32.MakeTable(crc32.Koopman)

// Hash returns a checksum for error correction.
func Hash(in []byte) []byte {
	// TODO: make sure this function does not cause race conditions due to using the table?
	hash := crc32.New(tablePolynomial)
	return hash.Sum(in)[:]
}
