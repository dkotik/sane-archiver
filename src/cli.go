// Provide convenient command line interface.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"syscall"
)

const usage = `Sane Archiver Alpha 0.0.4
A simple command line utility for making encrypted archives.

Usage:
  sane-archiver --keygen
  sane-archiver --key [PUBLICKEY] [FILE|DIRECTORY]... [OPTION]...
  sane-archiver --key [PRIVATEKEY] --decrypt [SANEFILE]

Options:
 -k, --key <KEY>    Set private or public base64-encoded key.
     --keygen       Generate a base64-encoded keypair.
 -o, --output       Output to this file or path.
 -d, --decrypt      Decrypt this file using the key.
 -f, --force        Overwrite any files that already exist.
 -w, --warn <GB>    Warn if the disk is running low on space.
 -h, --help         Print this message.

 Defaults:
   --output defaults to {year}-{month}-{day}-{md5}.[sane1|zip]
   --key [PUBLICKEY] defaults to $ENV[SaneArchiverPublicKey]
   --warn defaults to 2, issuing a warning under 2GB of free space
`

// CLI holds flags
type CLI struct {
	Key     string
	Keygen  bool
	Force   bool
	Output  string
	Decrypt string
	Warn    uint
	Sources []string
}

// Instructions prints CLI introduction and help.
func (c *CLI) Instructions() { fmt.Fprintf(flag.CommandLine.Output(), usage) }

// Warnings logs warnings that might prevent failed archivation.
func (c *CLI) Warnings() {
	var stat syscall.Statfs_t
	path := path.Dir(c.Output)
	if err := syscall.Statfs(path, &stat); err != nil {
		log.Printf("<WARNING> Storage device at path <%s> cannot be accessed!", path)
	} else if (stat.Bavail * uint64(stat.Bsize)) < uint64(c.Warn*1024*1024*1024) {
		log.Printf("<WARNING> Storage device at path <%s> has less than %dGB of space left!",
			path, c.Warn)
	}
}

// ParseCLIInput returns populated flag structure.
func ParseCLIInput() *CLI {
	result := CLI{}
	flag.StringVar(&result.Key, "key", *flag.String("k", "", ""), "")
	flag.BoolVar(&result.Keygen, "keygen", false, "")
	flag.BoolVar(&result.Force, "force", *flag.Bool("f", false, ""), "")
	flag.StringVar(&result.Output, "output", *flag.String("o", "", ""), "")
	flag.StringVar(&result.Decrypt, "decrypt", *flag.String("d", "", ""), "")
	flag.UintVar(&result.Warn, "warn", *flag.Uint("w", 2, ""), "")
	flag.Usage = result.Instructions
	flag.Parse()
	if result.Output == "" {
		if result.Decrypt == "" {
			result.Output = "{year}-{month}-{day}-{md5}.sane1"
		} else {
			result.Output = "{year}-{month}-{day}-{md5}.zip"
		}
	}
	if result.Decrypt != "" && !strings.HasSuffix(result.Output, ".zip") {
		result.Output += ".zip"
	}
	if result.Decrypt == "" && strings.HasSuffix(result.Output, ".zip") {
		result.Output = result.Output[:len(result.Output)-4] + ".sane1"
	}
	if result.Decrypt == "" && result.Key == "" {
		result.Key = os.Getenv("SaneArchiverPublicKey")
	}
	// Remaining arguments are sources to process.
	result.Sources = flag.Args()
	// Check if there are additional sources in STDIN.
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			result.Sources = append(result.Sources, scanner.Text())
		}
	}
	return &result
}
