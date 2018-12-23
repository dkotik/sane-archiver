// Package archiver is a simple command line utility for making encrypted archives.
package archiver

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// flags holds flags
type flags struct {
	Key     string
	Keygen  bool
	Force   bool
	Dryrun  bool
	Upload  string
	Output  string
	Master  bool
	Decrypt bool
	Warn    int
}

// Warnings logs warnings that might prevent failed archivation.
func (c *flags) Warnings() {
	var stat syscall.Statfs_t
	path := path.Dir(c.Output)
	if err := syscall.Statfs(path, &stat); err != nil {
		log.Printf("<WARNING> Storage device at path <%s> cannot be accessed!", path)
	} else if (stat.Bavail * uint64(stat.Bsize)) < uint64(c.Warn*1024*1024*1024) {
		log.Printf("<WARNING> Storage device at path <%s> has less than %dGB of space left!",
			path, c.Warn)
	}
}

// Flags holds all the command line values.
var Flags = flags{}

// ConfirmOverwrite makes sure user agrees with file overwrite operation.
func ConfirmOverwrite(target string) {
	if !Flags.Force {
		if _, err := os.Stat(target); err == nil {
			fmt.Printf("File <%s> already exists.\nOverwrite? (y/N): ", target)
			line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
			answer := strings.ToLower(string(line))
			if answer != `y` && answer != `yes` {
				log.Fatal("<CANCELLED> Operation aborted.")
			}
		}
	}
}

// ConfirmKey gathers key from io.Stdin if it is not provided.
func ConfirmKey() string {
	if len(Flags.Key) == 0 {

		// reader := bufio.NewReader(os.Stdin)
		if Flags.Decrypt {
			fmt.Println(`Please enter private key (-k) to decrypt target archives:`)
		} else {
			fmt.Println(`Please enter public key (-k) to create encrypted archive:`)
		}
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println(`Key received. Proceeding...`)
		// input, _ := reader.ReadString('\n')
		Flags.Key = strings.TrimSpace(string(bytePassword))
		if len(Flags.Key) == 0 {
			log.Fatalf(`Cannot perform operations using an empty key.`)
		}
	}
	return Flags.Key
}

// WriteHandle returns write handle and manages overwrites gracefully.
func WriteHandle(target string) *os.File {
	ConfirmOverwrite(target)
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Could not write the file to disk.")
	}
	return out
}
