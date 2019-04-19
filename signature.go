package golibs

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func GenerateSign(params map[string]string, secret string, taobao_nick string, time int64) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		if 64 == key[0] {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var string_to_signed *StringBuilder
	string_to_signed.Append(secret)
	for _, k := range keys {
		string_to_signed.Append(k).Append(params[k])
	}
	string_to_signed.Append(taobao_nick)
	string_to_signed.Append(fmt.Sprintf("%d", time))
	return Md5(string_to_signed.ToString())
}

func GenerateSignNoUser(params map[string]string, secret string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var string_to_signed *StringBuilder
	string_to_signed.Append(secret)
	for _, k := range keys {
		string_to_signed.Append(k).Append(params[k])
	}
	return Md5(string_to_signed.ToString())
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
	var string_to_signed *StringBuilder
	string_to_signed.Append(secret)
	for _, k := range keys {
		string_to_signed.Append(k).Append(params[k][0])
	}
	return strings.ToUpper(Md5(string_to_signed.ToString()))
}

// api签名，用于开放平台
func ApiSignIsValid(params url.Values, secret string) bool {
	sign := ""
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	data := NewStringBuilder()
	for _, k := range keys {
		if strings.EqualFold(k, "") || strings.EqualFold(params[k][0], "") {
			continue
		}
		if strings.EqualFold(k, "sign") {
			sign = params[k][0]
			continue
		}
		data.Append(k).Append(params[k][0])
	}
	return sign == fmt.Sprintf("%x", HmacMd5(StringToSliceByte(data.ToString()), StringToSliceByte(secret)))
}

// api签名，用于开放平台
func ApiSign(params url.Values, secret string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	data := NewStringBuilder()
	for _, k := range keys {
		if strings.EqualFold(k, "") || strings.EqualFold(params[k][0], "") {
			continue
		}
		if strings.EqualFold(k, "sign") {
			continue
		}
		data.Append(k).Append(params[k][0])
	}
	return fmt.Sprintf("%x", HmacMd5(StringToSliceByte(data.ToString()), StringToSliceByte(secret)))
}
