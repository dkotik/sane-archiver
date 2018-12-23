package archiver

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// SaneDirectoryWalker feeds files and directories into writer.
type SaneDirectoryWalker string

func (d SaneDirectoryWalker) signal(message string, args ...interface{}) {
	message = ` @` + string(d) + ` ` + message
	if Flags.Dryrun {
		message = `DRYRUN >` + message
	}
	log.Printf(message, args...)
}

func (d SaneDirectoryWalker) processGitDirectory(w *SaneWriter, path string) error {
	if l, err := GitBranchList(path); err == nil {
		d.signal(`Detected git repository at <%s> with branches %q.`, path, l)
		if !Flags.Dryrun {
			for _, branch := range l {
				cmd, r, gw := GitArchiveReader(path, branch)
				go func() {
					cmd.Run()
					gw.Close()
				}()
				if err := w.AddReader(strings.TrimSuffix(path, `.git`)+`-`+branch+`.tar`, &r); err != nil {
					d.signal("Git repository <%s> could not be accessed.", path)
					return err
				}
			}
		}
		// do not recurse within git directories
		return filepath.SkipDir
	}
	return nil
}

// Walk feeds discovered objects into writer.
func (d SaneDirectoryWalker) Walk(w *SaneWriter) {
	target := string(d)
	err := filepath.Walk(target,
		func(file string, info os.FileInfo, err error) error {
			if err != nil {
				d.signal("Path <%s> could not be accessed.", file)
				return err
			} else if info.IsDir() {
				return d.processGitDirectory(w, file)
			} else {
				d.signal("Adding <%s> to the archive.", file)
				if Flags.Dryrun {
					d.signal("Skipping file <%s> for dryrun.", file)
					return nil
				} else if err := w.AddFile(file); err != nil {
					d.signal("File <%s> could not be accessed.", file)
					return err
				}
				d.signal("File <%s> successfully added.", file)
			}
			return nil
		})
	if err == nil {
		d.signal("[DONE] Directory <%s> was fully archived!", target)
	} else {
		d.signal("Directory <%s> could not be fully processed. Reason: %s.", target, err)
	}
}
