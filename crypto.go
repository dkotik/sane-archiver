package archiver

// https://cyberspy.io/articles/crypto101/

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"log"
)

// KeyBytes guides the password encryption strength and determines tag length after nonce.
const KeyBytes = 128 // 1024 bits

// FromBase64 returns data as raw bytes.
func FromBase64(data string, label string) []byte {
	dec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Fatalf("Provided %s is corrupted.", label)
	}
	return dec
}

// SetupSymmetricCipherBlock returns a keyed symmetric cipher.
func SetupSymmetricCipherBlock(key []byte) cipher.Block {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("Provided key is corrupted.")
	}
	return block
}

// GenerateKeyPair returns private and public keys in base64 encoding.
func GenerateKeyPair() (string, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, KeyBytes*8)
	if err != nil {
		log.Fatalf("Could not generate a new key pair! Reason: %s.", err)
	}
	mPublic, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		log.Fatalf("Could not generate a new key pair! Reason: %s.", err)
	}
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privateKey)),
		base64.StdEncoding.EncodeToString(mPublic)
}

// Encrypt hides a message using public key.
func Encrypt(base64PublicKey string, message []byte) []byte {
	publicKey, err := x509.ParsePKIXPublicKey(FromBase64(base64PublicKey, "public key"))
	if err != nil {
		log.Fatalf("Could not load the public key! Reason: %s.", err)
	}
	key := publicKey.(*rsa.PublicKey)
	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, message, []byte("label"))
	if err != nil {
		log.Fatalf("Error encrypting: %s.", err)
	}
	return cipherText
}

// Decrypt recovers a message using private key.
func Decrypt(base64PrivateKey string, message []byte) []byte {
	privateKey, err := x509.ParsePKCS1PrivateKey(FromBase64(base64PrivateKey, "private key"))
	if err != nil {
		log.Fatalf("Provided private key is corrupted: %s.", err.Error())
	}
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, message, []byte("label"))
	if err != nil {
		log.Fatalf("Error decrypting message: %s\n", err.Error())
	}
	return plaintext
}
