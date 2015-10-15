package upload

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	timeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"
)

func sign(msg, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func AliUploadWithMtype(accessKey, secret, baseUrl, bucket, objectName, dir, mType string, fileBuffer io.Reader) (err error) {
	client := &http.Client{}
	var canonicalizedResource, canonicalizedOSSHeaders, verb string
	canonicalizedResource += "/" + bucket + "/" + dir + "/" + objectName
	canonicalizedOSSHeaders += "x-oss-object-acl:public-read"
	verb = "PUT"
	str := verb + "\n" + "\n" + mType + "\n" + time.Now().UTC().Format(timeFormat) + "\n" + canonicalizedOSSHeaders + "\n" + canonicalizedResource
	sign := sign(str, secret)
	authorization := "OSS " + accessKey + ":" + sign
	req, err := http.NewRequest("PUT", baseUrl+canonicalizedResource, fileBuffer)
	req.Header.Add("x-oss-object-acl", "public-read")
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", mType)
	req.Header.Add("Date", time.Now().UTC().Format(timeFormat))
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode >= 300 {
		return fmt.Errorf("status err %d", resp.StatusCode)
	}
	return nil
}
