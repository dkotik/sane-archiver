package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

// CLI holds the full configuration for the command line interface.
var CLI struct {
	Pack    packTask         `kong:"cmd,help='Pack files or folders into an encrypted archive.'"`
	Unpack  unpackTask       `kong:"cmd,help='Unpack all provided files.'"`
	Keygen  keygenTask       `kong:"cmd,help='Generate a base64-encoded keypair.'"`
	Version kong.VersionFlag `kong:"hidden,short='v',help='Display version information.'"`
}

// ConfirmOverwrite makes sure user agrees with file overwrite operation.
func ConfirmOverwrite(target string) {
	// TODO: this will not work for writer path?
	if _, err := os.Stat(target); err == nil {
		fmt.Printf("File <%s> already exists.\nOverwrite? (y/N): ", target)
		line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
		answer := strings.ToLower(string(line))
		if answer != `y` && answer != `yes` {
			log.Fatal("<CANCELLED> Operation aborted.")
		}
	}
}

func main() {
	log.SetOutput(os.Stderr)
	err := func() error {
		c, err := kong.New(&CLI,
			kong.Description(`A simple command line utility for making encrypted archives.`),
			kong.Vars{"version": "0.1.2"},
			kong.ConfigureHelp(kong.HelpOptions{
				Compact: true,
				Summary: true,
			}))
		if err != nil {
			return err
		}
		if len(os.Args) <= 1 {
			ctx, err := c.Parse([]string{`--help`})
			if err != nil {
				return err
			}
			return ctx.Run()
		}
		ctx, err := c.Parse(os.Args[1:])
		if err != nil {
			return err
		}
		return ctx.Run()
	}()
	if err != nil {
		log.Fatalf(`Archiver failed: %s.`, err.Error())
	}
}
