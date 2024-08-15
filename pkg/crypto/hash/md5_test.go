package hash

import "testing"

const (
	testDataMD5Encrypt123456 = "e10adc3949ba59abbe56e057f20f883e"
)

func TestMD5(t *testing.T) {
	if MD5("123456") != testDataMD5Encrypt123456 {
		t.Fatalf("md5 encrypt failed\n")
	}
}
