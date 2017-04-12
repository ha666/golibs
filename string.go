package golibs

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

//字串截取
func SubString(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func GetFileSuffix(s string) string {
	re, _ := regexp.Compile(".(jpg|jpeg|png|gif|exe|doc|docx|ppt|pptx|xls|xlsx)")
	suffix := re.ReplaceAllString(s, "")
	return suffix
}

func RandInt64(min, max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < min {
		RandInt64(min, max)
	}
	return i.Int64()
}

func Strim(str string) string {
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	return str
}

func String2Unicode(rs string) string {
	json := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			json += string(r)
		} else {
			json += "\\u" + strconv.FormatInt(int64(rint), 16)
		}
	}
	return json
}

func HTMLEncode(rs string) string {
	html := ""
	for _, r := range rs {
		html += "&#" + strconv.Itoa(int(r)) + ";"
	}
	return html
}

func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(Base64(b))
}

func GetIPNums(ipAddress string) (ipNum uint32, err error) {
	if strings.EqualFold(ipAddress, "") {
		return ipNum, errors.New("ipAddress is null")
	}
	items := strings.Split(ipAddress, ".")
	if len(items) != 4 {
		return ipNum, errors.New("ipAddress is error")
	}
	item0, err := strconv.Atoi(items[0])
	if err != nil {
		return ipNum, errors.New("ipAddress is error0")
	}
	item1, err := strconv.Atoi(items[1])
	if err != nil {
		return ipNum, errors.New("ipAddress is error1")
	}
	item2, err := strconv.Atoi(items[2])
	if err != nil {
		return ipNum, errors.New("ipAddress is error2")
	}
	item3, err := strconv.Atoi(items[3])
	if err != nil {
		return ipNum, errors.New("ipAddress is error3")
	}
	return uint32(item0<<24 | item1<<16 | item2<<8 | item3), nil
}

func IsTaobaoNick(taobaoNick string) bool {
	if len(taobaoNick) < 2 {
		return false
	}
	return regexp.MustCompile(`(^[\\u4e00-\\u9fa5\\w_—\\-，。…·〔〕（）！@￥%……&*？、；‘“]*$)`).MatchString(taobaoNick)
}

const zipOffset int = 19968

func ZipMd5(md5String string) (zipString string, err error) {
	md5Bytes := getHexBytes(md5String + "0")
	var data bytes.Buffer
	var total int = 0
	for index := 0; index < len(md5Bytes); index++ {
		switch index % 3 {
		case 0:
			intByte, err := strconv.Atoi(fmt.Sprintf("%d", md5Bytes[index]))
			if err != nil {
				return zipString, err
			}
			total = intByte
			break
		case 1:
			intByte, err := strconv.Atoi(fmt.Sprintf("%d", md5Bytes[index]))
			if err != nil {
				return zipString, err
			}
			total += intByte << 4
			break
		case 2:
			intByte, err := strconv.Atoi(fmt.Sprintf("%d", md5Bytes[index]))
			if err != nil {
				return zipString, err
			}
			total += (intByte << 8) + zipOffset
			uniString := fmt.Sprintf("\\u%x", total)
			chineseString, err := Unicode2String(uniString)
			if err != nil {
				return zipString, err
			}
			data.WriteString(chineseString)
			total = 0
			break
		}
	}
	zipString = data.String()
	return
}

func UnZipMd5(zipString string) (md5String string, err error) {
	var data bytes.Buffer
	unicodeString := String2Unicode(zipString)
	unicodeStrings := strings.Split(unicodeString, "\\u")
	for _, unicodeValue := range unicodeStrings {
		if strings.EqualFold(unicodeValue, "") {
			continue
		}
		dec, err := strconv.ParseInt(unicodeValue, 16, 32)
		if err != nil {
			return md5String, err
		}
		dec -= int64(zipOffset)
		data.WriteString(tenValue2Char(dec & 15))
		data.WriteString(tenValue2Char((dec >> 4) & 15))
		data.WriteString(tenValue2Char((dec >> 8) & 15))
	}
	md5String = SubString(data.String(), 0, 32)
	return
}

func Unicode2String(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}

func getHexBytes(str string) []byte {
	result := []byte(str)
	for index := 0; index < len(result); index++ {
		if result[index] < 58 {
			result[index] -= 48
		} else {
			result[index] -= 55
		}
	}
	return result
}

func tenValue2Char(ten int64) string {
	switch ten {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
		return strconv.FormatInt(ten, 10)
	case 10:
		return "A"
	case 11:
		return "B"
	case 12:
		return "C"
	case 13:
		return "D"
	case 14:
		return "E"
	case 15:
		return "F"
	}
	return ""
}
