package archiver

import (
	"fmt"
	"testing"
)

var testURLExample = `s3://<credentialID>:<credentialSecret>@<awsRegion>/<bucket>/<path>`
var testURL = ``

func TestS3(t *testing.T) {
	if testURL != `` {
		fmt.Print(UploadS3(`README.md`, testURL))
	}
}
