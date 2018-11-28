package archiver

import (
	// "reflect"
	"testing"
    // "log"
)

func TestGitBranches(t *testing.T) {
    _, err := GitBranchList(`.`)
    if err != nil {
        t.Fatal(err)
    }
}
