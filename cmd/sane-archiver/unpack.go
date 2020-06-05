package main

import (
	"archiver"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/alecthomas/kong"
	"golang.org/x/crypto/ssh/terminal"
)

type unpackTask struct {
	Key    string   `kong:"flag,help='Private base64-encoded key.'"`
	File   []string `kong:"arg,required,help='File to unpack.',type='existingfile',sep=' '"`
	Output string   `kong:"flag,name='output',short='o',type='path',help='Output directory.',default='.'"`
	Force  bool     `kong:"flag,name='force',short='f',help='Overwrite any files that already exist.'"`
}

func (c *unpackTask) Run(ctx *kong.Context) error {
	if c.Key == "" {
		fmt.Println(`Please enter private key (-k) to decrypt target archives:`)
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		c.Key = strings.TrimSpace(string(bytePassword))
	}
	info, err := os.Stat(c.Output)
	if err != nil || (err == nil && !info.IsDir()) {
		// TODO: should be able to take an output file!
		return fmt.Errorf(`output directory <%s> must be writable`, c.Output)
	}

	var p string
	for _, arg := range c.File {
		p = path.Join(c.Output, strings.TrimSuffix(filepath.Base(arg), `.sane1`)+`.zip`)
		if !c.Force {
			ConfirmOverwrite(p)
		}
		err := archiver.Decode(p, arg, c.Key)
		if err != nil {
			return fmt.Errorf("could not decrypt file <%s>: %w", arg, err)
		}
	}
	return nil
}
