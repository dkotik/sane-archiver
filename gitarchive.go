package archiver

// Implement git integration for smart archiving.

import (
    // "bytes"
    "os"
    "os/exec"
    "strings"
    // "log"
    "io"
    "fmt"
)

const gitCleanBranchList = `git branch | awk -F ' +' '! /\(no branch\)/ {print $2}'`
const gitArchiveBranch = `git archive "%s"`

// GitBranchList lists all local git branches found in current directory.
func GitBranchList(path string) ([]string, error) {
    result := make([]string, 0)
    cmd := exec.Command("/bin/bash", "-c", gitCleanBranchList)
    cmd.Dir = path
    out, err := cmd.Output()
    if err == nil {
        for _, branch := range strings.Split(string(out), "\n") {
            if len(branch) > 0 {
                result = append(result, branch)
                // log.Println(`adding`, branch, `|`)
            }
        }
    }
    // log.Println(`command output:`, out)

    // p := exec.Command("/bin/bash", "-c", gitCleanBranchList)
    // p.Dir = path
    // // p.Stdin = r
    // p.Stdout = buffer
    // p.Stderr = os.Stderr
    // p.Run()
    if len(result) == 0 {
        return result, fmt.Errorf(`path <%s> does not appear to be a git folder`, path)
    }
    return result, nil
}

// GitArchiveReader adds a git archive as a stream.
func GitArchiveReader(path string, branch string) (*exec.Cmd, io.Reader, *io.PipeWriter) {
    r, w := io.Pipe()
    p := exec.Command("/bin/bash", "-c", fmt.Sprintf(gitArchiveBranch, branch))
    p.Dir = path
    // p.Stdin = r
    p.Stdout = w
    p.Stderr = os.Stderr
    return p, r, w
}
