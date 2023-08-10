package utils

import (
	"database/sql"
	"unsafe"

	"github.com/gookit/goutil/strutil"
)

// String convert[]byte to string
func String(data []byte) string {
	return unsafe.String(&data[0], len(data))
}

// Bytes convert string to []byte
func Bytes(data string) []byte {
	return unsafe.Slice(unsafe.StringData(data), len(data))
}

// AnyBlank has any empty value return true otherwise return false
func AnyBlank(args ...string) bool {
	for _, v := range args {
		if strutil.IsBlank(v) {
			return true
		}
	}
	return false
}

// NonBlanks has any empty string returns false otherwise return true
func NonBlanks(args ...string) bool {
	for _, v := range args {
		if strutil.IsBlank(v) {
			return false
		}
	}
	return true
}

// Blanks all value empty return true otherwise return false
func Blanks(args ...string) bool {
	for _, v := range args {
		if strutil.IsNotBlank(v) {
			return false
		}
	}
	return true
}

// SQLString returns a valid sql.NullString
func SQLString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func IsEmail(email string) bool {
	return false
}
