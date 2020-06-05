package archiver

import (
	"archive/zip"
	"crypto/cipher"
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"time"
)

var unrootPath = regexp.MustCompile(`^\.*\/+`)

// SaneWriter is a wrapped writer.
type SaneWriter struct {
	Writer    io.Writer
	PublicKey string // Base64-encoded public key used for encryption.
	Hash      hash.Hash
	Size      uint64

	headerReady   bool
	cipherHandle  *cipher.StreamWriter
	archiveHandle *zip.Writer
}

// writeHeader prepares everything neccessary for writing the encrypted file.
func (w *SaneWriter) writeHeader() (err error) {
	if !w.headerReady {
		if len(w.PublicKey) == 0 {
			return fmt.Errorf(`cannot perform operations using an empty key`)
		}
		nonce, key, secret := MakeNonceKeySecret(w.PublicKey)
		w.Hash = md5.New()
		fork := io.MultiWriter(w.Hash, w.Writer)
		nA, err := fork.Write(nonce)
		if err != nil {
			return err
		}
		nB, err := fork.Write(secret)
		if err != nil {
			return err
		}

		w.Size = uint64(nA + nB)
		// TODO: cipher.NewOFB was used before, but that may cause problems with bit-rot.
		w.cipherHandle = &cipher.StreamWriter{
			S: cipher.NewCTR(SetupSymmetricCipherBlock(key), nonce), W: fork}
		w.archiveHandle = zip.NewWriter(w.cipherHandle)
		w.headerReady = true
	}
	return err
}

// AddFile writes target file into the encrypted archive.
func (w *SaneWriter) AddFile(target string) (err error) {
	err = w.writeHeader()
	if err != nil {
		return err
	}
	in, err := os.Open(target)
	if err != nil {
		return err
	}
	defer in.Close()
	info, err := in.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = unrootPath.ReplaceAllString(path.Clean(target), "")
	header.Comment = `Created by sane-archiver.`
	header.Method = zip.Deflate
	f, err := w.archiveHandle.CreateHeader(header)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, in)
	if err != nil {
		return err
	}
	w.Size += uint64(n)
	log.Printf("File <%s> was added.", target)
	return nil
}

// AddReader writes contents of provided io.Reader into the archive.
func (w *SaneWriter) AddReader(name string, target *io.Reader) (err error) {
	err = w.writeHeader()
	if err != nil {
		return err
	}
	header := &zip.FileHeader{
		Name:     name,
		Comment:  `Created by sane-archiver.`,
		Modified: time.Now(),
		NonUTF8:  false,
		Method:   zip.Deflate,
	}
	// f, err := w.archiveHandle.Create(name)
	f, err := w.archiveHandle.CreateHeader(header)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, *target)
	if err != nil {
		return err
	}
	w.Size += uint64(n)
	log.Printf("File <%s> was added.", name)
	return nil
}

// Close function closes the active IO handles and relocates the file.
func (w *SaneWriter) Close() error {
	if !w.headerReady {
		log.Fatalf(`No files were added to the archive!`)
	}
	w.archiveHandle.Close()
	return w.cipherHandle.Close()
}
