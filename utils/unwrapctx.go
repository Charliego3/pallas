package utils

import (
	"context"
	"github.com/Charliego93/go-i18n/v2"
	"golang.org/x/text/language"
	"google.golang.org/grpc/metadata"
)

var DefaultLanguage language.Tag

func init() {
	SetDefaultLanguage(language.Chinese)
}

func SetDefaultLanguage(lang language.Tag) {
	DefaultLanguage = lang
}

func Language(ctx context.Context) language.Tag {
	lang := DefaultLanguage
	if l, ok := MDExtract(ctx, "accept-language"); ok {
		tag := i18n.ParseFromHeader(l)
		if tag != language.Und {
			lang = tag
		}
	}
	return lang
}

func MDExtract(ctx context.Context, key string) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	val := md[key]
	if len(val) == 0 {
		return "", false
	}
	return val[0], true
}
