package utils

import (
	"crypto/md5"
	"fmt"
)

func MD5Sum(data, seed string) string {
	hash := md5.New()
	hash.Write([]byte(data + seed))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
