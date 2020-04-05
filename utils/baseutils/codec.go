package baseutils

import (
	"encoding/base64"
)

// Base64DecodeByte2Str returns plain text as string from the encrypted text as byte array
func Base64DecodeByte2Str(enc []byte) string {
	encStr := string(enc)
	decStr, err := base64.StdEncoding.DecodeString(encStr)
	if err != nil {
		return ""
	}
	return string(decStr)
}
