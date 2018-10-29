package main

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var unrootPath = regexp.MustCompile(`^\.*\/+`)

// SaneWriter is a wrapped writer.
type SaneWriter struct {
	output        string
	size          uint64
	fileHandle    *os.File
	cipherHandle  *cipher.StreamWriter
	archiveHandle *zip.Writer
	hash          hash.Hash
	nonce         []byte
}

// NewWriter returns a cryptographic writer after setting it up.
func NewWriter(target string, base64PublicKey string) *SaneWriter {
	w := SaneWriter{output: target}
	w.nonce = make([]byte, aes.BlockSize)
	rand.Read(w.nonce)
	target = fmt.Sprintf("%s%c%s~%x.tmp",
		path.Dir(target), os.PathSeparator, time.Now().Format("2006-01-31"), w.nonce)
	w.fileHandle = WriteHandle(target)

	w.hash = md5.New()
	key := make([]byte, aes.BlockSize)
	rand.Read(key)
	fork := io.MultiWriter(w.hash, w.fileHandle)
	nA, err := fork.Write(w.nonce)
	if err != nil {
		log.Fatalf("Unable to write to <%s>.", target)
	}
	nB, err := fork.Write(Encrypt(base64PublicKey, key))
	if err != nil {
		log.Fatalf("Unable to write to <%s>.", target)
	}
	w.size = uint64(nA + nB)
	// TODO: cipher.NewOFB was used before, but that may cause problems with bit-rot.
	w.cipherHandle = &cipher.StreamWriter{
		S: cipher.NewCTR(SetupSymmetricCipherBlock(key), w.nonce), W: fork}
	w.archiveHandle = zip.NewWriter(w.cipherHandle)
	return &w
}

// Walk adds files listed in target path to archive.
func (w *SaneWriter) Walk(target string) {
	err := filepath.Walk(target,
		func(file string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Path <%s> could not be accessed.", target)
				return err
			} else if !info.IsDir() {
				if err := w.addFile(file); err != nil {
					log.Printf("File <%s> could not be accessed.", target)
					return err
				}
			}
			return nil
		})
	if err != nil {
		log.Printf("Directory <%s> could not be accessed.", target)
	}
}

func (w *SaneWriter) addFile(target string) error {
	in, err := os.OpenFile(target, os.O_RDONLY, 0755)
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
	header.Method = zip.Deflate
	f, err := w.archiveHandle.CreateHeader(header)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, in)
	if err != nil {
		return err
	}
	w.size += uint64(n)
	log.Printf("File <%s> was added.", target)
	return nil
}

// Close function closes the active IO handles and relocates the file.
func (w *SaneWriter) Close() {
	oldpath := w.fileHandle.Name()
	log.Printf("Wrote %.2fGB to <%s>.", float64(w.size)/(1024*1024*1024), oldpath)
	w.archiveHandle.Close()
	w.cipherHandle.Close()
	w.fileHandle.Close()
	newpath := w.output
	t := time.Now()
	newpath = strings.Replace(newpath, "{day}", fmt.Sprintf("%d", t.Day()), 1)
	newpath = strings.Replace(newpath, "{month}", fmt.Sprintf("%d", t.Month()), 1)
	newpath = strings.Replace(newpath, "{year}", fmt.Sprintf("%d", t.Year()), 1)
	newpath = strings.Replace(newpath, "{md5}", hex.EncodeToString(w.hash.Sum(nil)), 1)
	ConfirmOverwrite(newpath)
	err := os.Rename(oldpath, newpath)
	if err == nil {
		log.Printf("Wrote %.2fGB to <%s>.", float64(w.size)/(1024*1024*1024), newpath)
	} else {
		log.Printf("<ERROR> Unable to rename file <%s> to <%s>.", oldpath, newpath)
	}
	os.Stdout.WriteString(newpath)
}
