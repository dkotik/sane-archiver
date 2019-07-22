package archiver

import (
	"errors"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// ErrS3URLError is a generic message indicating that provided S3 URL is improper.
var ErrS3URLError = errors.New(`provided S3 URL does not follow the proper format: s3://<credentialID>:<credentialSecret>@<awsRegion>/<bucket>/<path>`)

// UploadS3 pushes one file to AWS S3 bucket.
func UploadS3(file string, URL string) error {
	if strings.HasSuffix(URL, `/`) {
		URL += filepath.Base(file)
	}
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}
	i := strings.Index(u.Path, "/")
	password, ok := u.User.Password()
	if u.Scheme != `s3` || i == -1 || u.User == nil || !ok {
		return ErrS3URLError
	}
	bucketName, keyName := u.Path[:i+1], u.Path[i+1:]

	// user := &u.User
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(u.Host),
		Credentials: credentials.NewStaticCredentials(u.User.Username(), password, ""),
	})
	if err != nil {
		return ErrS3URLError
	}
	handle, err := os.Open(file)
	if err != nil {
		return err
	}
	defer handle.Close()
	uploader := s3manager.NewUploader(sess)
	upParams := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &keyName,
		Body:   handle,
	}
	// Perform an upload.
	_, err = uploader.Upload(upParams)
	// spew.Dump(result)
	// Perform upload with options different than the those in the Uploader.
	// result, err := uploader.Upload(upParams, func(u *s3manager.Uploader) {
	// 	u.PartSize = 10 * 1024 * 1024 // 10MB part size
	// 	u.LeavePartsOnError = true    // Don't delete the parts if the upload fails.
	// })
	if err == nil {
		log.Printf("[S3] Successfully uploaded %s to %s.", file, u.Path)
	} else {
		log.Printf("[S3] There was an error uploading %s to %s: %s.", file, u.Path, err.Error())
	}
	return err
}
