package encryption

import (
	"crypto/md5"
	"encoding/hex"
)

func MakeHash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
