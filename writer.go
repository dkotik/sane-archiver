package archiver

import (
	"archive/zip"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

var unrootPath = regexp.MustCompile(`^\.*\/+`)

// SaneWriter is a wrapped writer.
type SaneWriter struct {
	Output string // Destination file.
	Result string // Final resting place of the file.

	headerReady   bool
	size          uint64
	fileHandle    *os.File
	cipherHandle  *cipher.StreamWriter
	archiveHandle *zip.Writer
	hash          hash.Hash
}

// writeHeader prepares everything neccessary for writing the encrypted file.
func (w *SaneWriter) writeHeader() error {
	if !w.headerReady {
		nonce, key, secret := MakeNonceKeySecret(ConfirmKey())
		target := fmt.Sprintf("%s%c%s~%x.tmp",
			path.Dir(w.Output), os.PathSeparator,
			time.Now().Format("2006-01-31"),
			nonce)
		w.fileHandle = WriteHandle(target)

		w.hash = md5.New()
		fork := io.MultiWriter(w.hash, w.fileHandle)
		nA, err := fork.Write(nonce)
		if err != nil {
			log.Fatalf("Unable to write to <%s>.", target)
		}
		nB, err := fork.Write(secret)
		if err != nil {
			log.Fatalf("Unable to write to <%s>.", target)
		}
		w.size = uint64(nA + nB)
		// TODO: cipher.NewOFB was used before, but that may cause problems with bit-rot.
		w.cipherHandle = &cipher.StreamWriter{
			S: cipher.NewCTR(SetupSymmetricCipherBlock(key), nonce), W: fork}
		w.archiveHandle = zip.NewWriter(w.cipherHandle)
		w.headerReady = true
	}
	return nil
}

// AddFile writes target file into the encrypted archive.
func (w *SaneWriter) AddFile(target string) error {
	w.writeHeader()
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
	w.size += uint64(n)
	log.Printf("File <%s> was added.", target)
	return nil
}

// AddReader writes contents of provided io.Reader into the archive.
func (w *SaneWriter) AddReader(name string, target *io.Reader) error {
	w.writeHeader()
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
	w.size += uint64(n)
	log.Printf("File <%s> was added.", name)
	return nil
}

// Close function closes the active IO handles and relocates the file.
func (w *SaneWriter) Close() error {
	if !w.headerReady {
		log.Fatalf(`No files were added to the archive!`)
	}
	oldpath := w.fileHandle.Name()
	log.Printf("Wrote %.2fGB to <%s>.", float64(w.size)/(1024*1024*1024), oldpath)
	w.archiveHandle.Close()
	w.cipherHandle.Close()
	w.fileHandle.Close()
	newpath := w.Output
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
	w.Result = newpath
	return nil
}
