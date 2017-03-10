package golibs

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func GenerateSign(params map[string]string, secret string, taobaonick string, time int64) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		if 64 == key[0] {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	stringToBeSigned := secret
	for _, k := range keys {
		stringToBeSigned += k + params[k]
	}
	stringToBeSigned += taobaonick
	stringToBeSigned += fmt.Sprintf("%d", time)
	return Md5(stringToBeSigned)
}

func GenerateSignNoUser(params map[string]string, secret string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	stringToBeSigned := secret
	for _, k := range keys {
		stringToBeSigned += k + params[k]
	}
	return Md5(stringToBeSigned)
}

func TopSign(params url.Values, secret string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		if 64 == key[0] {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	stringToBeSigned := secret
	for _, k := range keys {
		stringToBeSigned += k + params[k][0]
	}
	return strings.ToUpper(Md5(stringToBeSigned))
}
