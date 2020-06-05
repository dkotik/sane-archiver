package archiver

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// SaneDirectoryWalker feeds files and directories into writer.
type SaneDirectoryWalker struct {
	Target         string
	Dryrun, Master bool
}

func (d *SaneDirectoryWalker) signal(message string, args ...interface{}) {
	message = fmt.Sprintf(` @%s %s`, d.Target, message)
	if d.Dryrun {
		message = `DRYRUN >` + message
	}
	log.Printf(message, args...)
}

func (d *SaneDirectoryWalker) getGitBranches(p string) ([]string, error) {
	l, err := GitBranchList(p)
	if err != nil {
		return l, err
	}
	if d.Master {
		for _, branch := range l {
			if branch == `master` {
				return []string{`master`}, nil
			}
		}
		d.signal(`Skipping directory <%s> because it does not have a master branch.`, p, l)
		return []string{}, nil
	}
	return l, err
}

func (d *SaneDirectoryWalker) processGitDirectory(w *SaneWriter, path string) error {
	l, err := d.getGitBranches(path)
	if err != nil {
		return err
	}
	d.signal(`Detected git repository at <%s> with branches %q.`, path, l)
	if !d.Dryrun {
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

// Walk feeds discovered objects into writer.
func (d *SaneDirectoryWalker) Walk(w *SaneWriter) (err error) {
	err = filepath.Walk(d.Target,
		func(file string, info os.FileInfo, err error) error {
			if err != nil {
				d.signal("Path <%s> could not be accessed.", file)
				return err
			} else if info.IsDir() {
				return d.processGitDirectory(w, file)
			} else {
				d.signal("Adding <%s> to the archive.", file)
				if d.Dryrun {
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
		d.signal(`[DONE] "%s" was fully archived!`, d.Target)
	} else {
		d.signal(`"%s" could not be fully processed. Reason: %s.`, d.Target, err)
	}
	return err
}
