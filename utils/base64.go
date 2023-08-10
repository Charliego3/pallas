package utils

import "encoding/base64"

func B64Decode(src []byte) []byte {
	dest := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	_, _ = base64.StdEncoding.Decode(dest, src)
	return dest
}
