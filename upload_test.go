package upload

import (
	"os"
	"testing"
)

var (
	key     = ""
	secret  = ""
	baseUrl = "http://oss-cn-hangzhou.beijing.com"
	bucket  = ""
)

func TestUpload(t *testing.T) {
	f, err := os.Open("test.html")
	if err != nil {
		t.Fatalf("err %v", err)
	}
	defer f.Close()
	err = AliUploadWithMtype(key, secret, baseUrl, bucket, "test.html", "test1", "text/html", f)
	if err != nil {
		t.Fatalf("err %v", err)
	}

}
