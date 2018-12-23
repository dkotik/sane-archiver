package archiver

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"log"
	"os"
)

// Decode decrypts stored file.
func Decode(output string, target string, base64PrivateKey string) error {
	in, err := os.OpenFile(target, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer in.Close()

	nonce := make([]byte, aes.BlockSize)
	_, err = in.Read(nonce)
	if err != nil {
		return err
	}
	key := make([]byte, KeyBytes)
	_, err = in.Read(key)
	if err != nil {
		return err
	}
	cipherHandle := &cipher.StreamReader{
		// TODO: cipher.NewOFB was used before, but that may cause problems with bit-rot.
		S: cipher.NewCTR(SetupSymmetricCipherBlock(Decrypt(base64PrivateKey, key)), nonce), R: in}
	out := WriteHandle(output)
	defer out.Close()
	_, err = io.Copy(out, cipherHandle)
	if err != nil {
		return err
	}
	log.Printf("Archive successfully recovered and stored as %s.", out.Name())
	return nil
}
