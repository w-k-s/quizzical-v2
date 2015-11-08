package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
)

func Hash(strings ...string) string {

	var buffer bytes.Buffer
	for _, str := range strings {
		buffer.WriteString(str)
	}

	hasher := md5.New()
	hasher.Write([]byte(buffer.String()))
	return hex.EncodeToString(hasher.Sum(nil))
}
