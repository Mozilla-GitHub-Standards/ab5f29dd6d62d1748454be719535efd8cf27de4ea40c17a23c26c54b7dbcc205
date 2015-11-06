package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// variable for swapping in testing
var listObjects = func(svc *s3.S3, bucket, prefix string) (objects []*s3.Object, prefixes []*s3.CommonPrefix, err error) {
	listParams := &s3.ListObjectsInput{
		Bucket:    aws.String(bucket),
		Delimiter: aws.String("/"),
		Prefix:    aws.String(prefix),
	}
	err = svc.ListObjectsPages(listParams, func(res *s3.ListObjectsOutput, lastpage bool) bool {
		prefixes = append(prefixes, res.CommonPrefixes...)
		objects = append(objects, res.Contents...)
		return true
	})
	if err != nil {
		return nil, nil, fmt.Errorf("listing %s/%s err: %s", bucket, prefix, err)
	}
	return
}

func setExpiresIn(d time.Duration, w http.ResponseWriter) {
	expiresAt := time.Now().Add(d)

	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%.0f", d.Seconds()))
	w.Header().Set("Expires", expiresAt.Format(http.TimeFormat))
}
