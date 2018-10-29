// Sane Archiver is a simple command line utility for making encrypted archives.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// Flags encapsulates all the command line arguments.
var Flags = ParseCLIInput()

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

// WriteHandle returns write handle and manages overwrites gracefully.
func WriteHandle(target string) *os.File {
	ConfirmOverwrite(target)
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal("Could not write the file to disk.")
	}
	return out
}

func main() {
	log.SetOutput(os.Stderr)
	if Flags.Keygen {
		private, public := GenerateKeyPair()
		fmt.Printf("Private key: %s\nPublic key: %s\n", private, public)
	} else if Flags.Key != "" && Flags.Decrypt != "" {
		err := Decode(Flags.Decrypt, Flags.Key)
		if err != nil {
			log.Printf("Could not decrypt file <%s>. Reason: <%s>.", Flags.Decrypt, err.Error())
		}
		Flags.Warnings()
	} else if Flags.Key != "" && len(Flags.Sources) > 0 {
		w := NewWriter(Flags.Output, Flags.Key)
		defer w.Close()
		for _, arg := range Flags.Sources {
			w.Walk(arg)
		}
		Flags.Warnings()
	} else {
		Flags.Instructions()
	}
}
