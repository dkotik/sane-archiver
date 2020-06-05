package main

import (
	"archiver"
	"bufio"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"golang.org/x/crypto/ssh/terminal"
)

type packTask struct {
	Key        string   `kong:"flag,help='Public base64-encoded key.',env='SaneArchiverPublicKey'"`
	Target     []string `kong:"arg,required,help='File or directory to pack.',type='path',sep=' '"`
	Output     string   `kong:"flag,name='output',short='o',type='path',help='Output to this file or path.'"`
	Force      bool     `kong:"flag,name='force',short='f',help='Overwrite any files that already exist.'"`
	Upload     string   `kong:"flag,name='upload',short='u',help='Upload finished archive to the cloud URI endpoint.'"`
	Warn       uint8    `kong:"flag,name='warn',short='w',help='Warn if the disk is running low on space. Issues a warning if there is less gigabytes left than the specified amount.',default='2'"`
	Leave      uint8    `kong:"flag,name='leave',short:'l',help='Delete older output-matching files, if more than the specified number.',default='12'"`
	MasterOnly bool     `kong:"flag,name='master-only',short='m',help='Archive only master branches of git repositories.'"`
	DryRun     bool     `kong:"flag,name='dry-run',short='n',help='Display operations without writing.'"`
}

func (t *packTask) outputDirFile() (string, string, error) {
	p := filepath.Clean(t.Output)
	dir := filepath.Dir(p)
	info, err := os.Stat(p)
	if err != nil {
		info, err = os.Stat(dir)
		if err == nil && info.IsDir() {
			return dir, filepath.Base(p), nil
		}
		return ``, ``, fmt.Errorf("directory %s does not exist", dir)
	}
	if info.IsDir() {
		return dir, `{year}-{month}-{day}-{md5}.sane1`, nil
	}
	return dir, filepath.Base(p), nil
}

func (t *packTask) Run(ctx *kong.Context) error {
	outputDir, outputFile, err := t.outputDirFile()
	if err != nil {
		return err
	}
	if t.Key == "" {
		fmt.Println(`Please enter public key (-k) to create encrypted archive:`)
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		t.Key = strings.TrimSpace(string(bytePassword))
	}
	// Check if there are additional sources in STDIN.
	info, err := os.Stdin.Stat()
	if err != nil {
		return err
	}
	if (info.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		args := t.Target
		for scanner.Scan() {
			args = append(args, scanner.Text())
		}
	}

	// display some warnings // TODO: make this a separate suite
	var stat syscall.Statfs_t
	if err := syscall.Statfs(outputDir, &stat); err != nil {
		log.Printf("<WARNING> Storage device at path <%s> cannot be accessed!", outputDir)
	} else if (stat.Bavail * uint64(stat.Bsize)) < uint64(t.Warn)*1024*1024*1024 {
		log.Printf("<WARNING> Storage device at path <%s> has less than %dGB of space left!",
			outputDir, t.Warn)
	}

	tmpfile, err := ioutil.TempFile(outputDir, ".sane-archiver-*.tmp")
	if err != nil {
		return err
	}
	w := &archiver.SaneWriter{PublicKey: t.Key, Writer: tmpfile}
	defer func() {
		w.Close()
		tmpfile.Close()
	}()

	for _, arg := range t.Target {
		a := &archiver.SaneDirectoryWalker{
			Target: arg,
			Master: t.MasterOnly,
			Dryrun: t.DryRun,
		}
		err = a.Walk(w)
		if err != nil {
			return fmt.Errorf(`could not pack %s: %w`, arg, err)
		}
	}
	n := time.Now()
	p := strings.Replace(outputFile, "{day}", fmt.Sprintf("%d", n.Day()), 1)
	p = strings.Replace(p, "{month}", fmt.Sprintf("%d", n.Month()), 1)
	p = strings.Replace(p, "{year}", fmt.Sprintf("%d", n.Year()), 1)
	p = strings.Replace(p, "{md5}", hex.EncodeToString(w.Hash.Sum(nil)), 1)
	t.Output = filepath.Join(outputDir, p)
	ConfirmOverwrite(t.Output)
	err = os.Rename(tmpfile.Name(), t.Output)
	if err != nil {
		return fmt.Errorf("cannot move file %s: %w", tmpfile.Name(), err)
	}
	// TODO: replace this with progress bar
	log.Printf("Wrote %.2fGB to <%s>.", float64(w.Size)/(1024*1024*1024), t.Output)
	os.Stdout.WriteString(t.Output + "\n")

	if t.Upload != "" {
		log.Println(`Attemping to upload result...`)
		err := archiver.UploadS3(t.Output, t.Upload)
		if err != nil {
			return fmt.Errorf(`uploading %s: %w`, p, err)
		}
	}
	if t.Leave > 0 {
		err = eliminateAllExcept(t.Output, int(t.Leave))
	}
	return err
}
