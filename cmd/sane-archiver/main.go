package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"sane-archiver"

	"github.com/spf13/cobra"
)

const usage = `Sane Archiver 0.0.6 Alpha
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

var flags = &archiver.Flags

// CLI cobra boiler plate.
var CLI = &cobra.Command{
	Use:     `sane-archiver`,
	Version: `0.0.6 Alpha`,
	Short:   `Sane Archiver: A simple command line utility for making encrypted archives.`,
	// Long: `Longer description..
	//         feel free to use a few lines here.
	//         `,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if there are additional sources in STDIN.
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				args = append(args, scanner.Text())
			}
		}

		if flags.Keygen {
			private, public := archiver.GenerateKeyPair()
			fmt.Printf("Private key: %s\nPublic key: %s\n", private, public)
		} else if flags.Decrypt != "" {
			if !strings.HasSuffix(flags.Output, `.zip`) {
				flags.Output = strings.TrimSuffix(flags.Output, `.sane1`) + `.zip`
			}
			err := archiver.Decode(flags.Decrypt, flags.Key)
			if err != nil {
				log.Printf("Could not decrypt file <%s>. Reason: <%s>.", flags.Decrypt, err.Error())
			}
			flags.Warnings()
		} else if len(args) > 0 {
			w := archiver.NewWriter(flags.Output, flags.Key)
			defer w.Close()
			for _, arg := range args {
				w.Walk(arg)
			}
			flags.Warnings()
		} else {
			cmd.Help()
		}
	},
}

func main() {
	log.SetOutput(os.Stderr)
	// https://godoc.org/github.com/spf13/pflag#FlagSet
	CLI.PersistentFlags().StringVarP(&flags.Key, `key`, `k`, os.Getenv("SaneArchiverPublicKey"),
		`Set private or public base64-encoded key ($ENV[SaneArchiverPublicKey] is default).`)
	CLI.PersistentFlags().BoolVar(&flags.Keygen, `keygen`, false, `Generate a base64-encoded keypair.`)
	CLI.PersistentFlags().StringVarP(&flags.Decrypt, `decrypt`, `d`, ``, `Decrypt this file using the key.`)
	CLI.PersistentFlags().BoolVarP(&flags.Force, `force`, `f`, false, `Overwrite any files that already exist.`)
	CLI.PersistentFlags().StringVarP(&flags.Output, `output`, `o`,
		`{year}-{month}-{day}-{md5}.sane1`, `Output to target file or path ({year}-{month}-{day}-{md5}.sane1 is default).`)
	CLI.PersistentFlags().IntVarP(&flags.Warn, `warn`, `w`, 2,
		`Warn if the disk is running low on space (default is 2, issuing a warning under 2GB of free space).`)
	CLI.MarkFlagRequired(`key`)
	CLI.SetHelpTemplate(usage)
	CLI.Execute()
}
