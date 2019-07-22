package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"

	archiver "sane-archiver"

	"github.com/spf13/cobra"
)

const basicHelp = `Sane Archiver 0.0.9 Alpha

Basic:
  sane-archiver --keygen
  sane-archiver --key [PUBLICKEY] [FILE|DIRECTORY]... [OPTION]...
  sane-archiver --key [PRIVATEKEY] --decrypt [SANEFILE]...
  sane-archiver --help
`

const fullHelp = `Sane Archiver 0.0.9 Alpha
A simple command line utility for making encrypted archives.

Usage:
  sane-archiver --keygen
  sane-archiver --key [PUBLICKEY] [FILE|DIRECTORY]... [OPTION]...
  sane-archiver --key [PRIVATEKEY] --decrypt [SANEFILE]...
  sane-archiver --help

Options:
 -k, --key <KEY>    Set private or public base64-encoded key.
     --keygen       Generate a base64-encoded keypair.
 -o, --output       Output to this file or path.
 -d, --decrypt      Decrypt all provided files.
 -u, --upload       Upload finished archive to the cloud.
 -f, --force        Overwrite any files that already exist.
 -w, --warn <GB>    Warn if the disk is running low on space.
 -l, --leave 12		Delete older output-matching files, if more than 12.
 -m, --master-only  Archive only master branches of git repositories.
 -n, --dry-run      Display operations without writing.
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
	Version: `0.0.9 Alpha`,
	Short:   `Sane Archiver: A simple command line utility for making encrypted archives.`,
	// Long: `Longer description..
	//         feel free to use a few lines here.
	//         `,
	Run: func(cmd *cobra.Command, args []string) {
		stat, _ := os.Stdin.Stat() // Check if there are additional sources in STDIN.
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				args = append(args, scanner.Text())
			}
		}
		if flags.Keygen {
			private, public := archiver.GenerateKeyPair()
			fmt.Printf("Private key: %s\nPublic key: %s\n", private, public)
		} else if len(args) == 0 {
			fmt.Print(basicHelp)
			return
		} else if flags.Decrypt {
			if flags.Output == `{year}-{month}-{day}-{md5}.sane1` {
				flags.Output = `.`
			} else {
				info, err := os.Stat(flags.Output)
				if err != nil || (err == nil && !info.IsDir()) {
					log.Fatalf(`Output location used for decryption <%s> must be a writable directory.`, flags.Output)
				}
			}
			for _, arg := range args {
				err := archiver.Decode(path.Join(flags.Output, arg+`.zip`), arg, archiver.ConfirmKey())
				if err != nil {
					log.Printf("Could not decrypt file <%s>. Reason: <%s>.", arg, err.Error())
				}
			}
			flags.Warnings()
		} else {
			w := &archiver.SaneWriter{Output: flags.Output}
			defer func() {
				w.Close()
				flags.Warnings()
				if flags.Upload != `` {
					log.Println(`Attemping to upload result...`)
					err := archiver.UploadS3(w.Result, flags.Upload)
					if err != nil {
						log.Printf(`Error uploading %s: %s.`, w.Result, err.Error())
					}
				}
				if v, _ := cmd.PersistentFlags().GetInt(`leave`); v > 0 {
					err := eliminateAllExcept(flags.Output, v)
					if err != nil {
						log.Fatalf(`Error cleaning up: %s.`, err)
					}
				}
			}()
			for _, arg := range args {
				walker := archiver.SaneDirectoryWalker(arg)
				walker.Walk(w)
			}
		}
	},
}

func main() {
	log.SetOutput(os.Stderr)
	// https://godoc.org/github.com/spf13/pflag#FlagSet
	CLI.PersistentFlags().StringVarP(&flags.Key, `key`, `k`, os.Getenv("SaneArchiverPublicKey"),
		`Set private or public base64-encoded key ($ENV[SaneArchiverPublicKey] is default).`)
	CLI.PersistentFlags().BoolVar(&flags.Keygen, `keygen`, false, `Generate a base64-encoded keypair.`)
	CLI.PersistentFlags().BoolVarP(&flags.Decrypt, `decrypt`, `d`, false, `Decrypt this file using the key.`)
	CLI.PersistentFlags().BoolVarP(&flags.Dryrun, `dry-run`, `n`, false, `Display operations without writing.`)
	CLI.PersistentFlags().BoolVarP(&flags.Force, `force`, `f`, false, `Overwrite any files that already exist.`)
	CLI.PersistentFlags().StringVarP(&flags.Upload, `upload`, `u`,
		``, `Upload finished archive to the cloud.`)
	CLI.PersistentFlags().StringVarP(&flags.Output, `output`, `o`,
		`{year}-{month}-{day}-{md5}.sane1`, `Output to target file or path ({year}-{month}-{day}-{md5}.sane1 is default).`)
	CLI.PersistentFlags().IntVarP(&flags.Warn, `warn`, `w`, 2,
		`Warn if the disk is running low on space (default is 2, issuing a warning under 2GB of free space).`)
	CLI.PersistentFlags().IntP(`leave`, `l`, 0, `Delete older output-matching files, if more than 12.`)
	CLI.PersistentFlags().BoolVarP(&flags.Master, `master-only`, `m`, false, `Archive only master branches of git repositories.`)
	// CLI.MarkFlagRequired(`key`)
	CLI.SetHelpTemplate(fullHelp)
	CLI.Execute()
}
