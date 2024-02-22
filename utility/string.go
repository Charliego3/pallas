package utility

import (
	"database/sql"
	"unicode"
	"unsafe"
)

// DString returns d if source blank else return source
func DString(source, d string) string {
	if IsBlank(source) {
		return d
	}
	return source
}

// String convert []byte to string
func String(data []byte) string {
	return unsafe.String(&data[0], len(data))
}

// Bytes convert string to []byte
func Bytes(data string) []byte {
	return unsafe.Slice(unsafe.StringData(data), len(data))
}

func IsBlank(b string) bool {
	for _, v := range b {
		if !unicode.IsSpace(v) {
			return false
		}
	}
	return true
}

func NonBlank(v string) bool {
	return !IsBlank(v)
}

// AnyBlank has any empty value return true otherwise return false
func AnyBlank(args ...string) bool {
	for _, v := range args {
		if IsBlank(v) {
			return true
		}
	}
	return false
}

// NonBlanks has any empty string returns false otherwise return true
func NonBlanks(args ...string) bool {
	for _, v := range args {
		if IsBlank(v) {
			return false
		}
	}
	return true
}

// Blanks all value empty return true otherwise return false
func Blanks(args ...string) bool {
	for _, v := range args {
		if NonBlank(v) {
			return false
		}
	}
	return true
}

// SQLString returns a valid sql.NullString
func SQLString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  len(s) > 0,
	}
}

func IsEmail(email string) bool {
	return false
}
