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
	"net"
	"regexp"
	"strconv"
	"strings"
	"unicode"
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

//字符串逆序
func ReverseString(s string) string {
	str := []rune(s)
	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
	return string(str)
}

//获取文件扩展名
func GetFileSuffix(s string) string {
	re, _ := regexp.Compile(".(jpg|jpeg|png|gif|exe|doc|docx|ppt|pptx|xls|xlsx)")
	suffix := re.ReplaceAllString(s, "")
	return suffix
}

//生成指定范围内的int64数字
func RandInt64(min, max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < min {
		RandInt64(min, max)
	}
	return i.Int64()
}

//删除空格、换行、空格等字符
func Strim(str string) string {
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	return str
}

//字符串转成Unicode编码
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

//Unicode编码转成字符串
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

//html编码
func HTMLEncode(rs string) string {
	html := ""
	for _, r := range rs {
		html += "&#" + strconv.Itoa(int(r)) + ";"
	}
	return html
}

//获取一个Guid
func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(Base64(b))
}

//把IP地址转成数字
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

//判断是否是淘宝用户名
func IsTaobaoNick(taobaoNick string) bool {
	if len(taobaoNick) < 2 {
		return false
	}
	return regexp.MustCompile(`(^[\\u4e00-\\u9fa5\\w_—\\-，。…·〔〕（）！@￥%……&*？、；‘“]*$)`).MatchString(taobaoNick)
}

//判断是否是字母、数字
func IsLetterOrNumber(name string) bool {
	if len(name) < 1 {
		return false
	}
	return regexp.MustCompile(`(^[A-Za-z0-9]*$)`).MatchString(name)
}

/*
判断字符串是否全中文字符
*/
func IsAllChineseChar(str string) bool {
	for _, r := range str {
		if !unicode.Is(unicode.Scripts["Han"], r) {
			return false
		}
	}
	return true
}

const zipOffset int = 19968

//压缩md5或guid
func ZipMd5(md5String string) (zipString string, err error) {
	if len(md5String) != 16 && len(md5String) != 32 {
		return "", errors.New("源md5值长度不对")
	}
	var md5Bytes []byte
	switch len(md5String) {
	case 16:
		md5Bytes = getHexBytes(md5String + "00")
	case 32:
		md5Bytes = getHexBytes(md5String + "0")
	}
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

//解压缩汉字，结果为guid或md5值
func UnZipMd5(zipString string) (md5String string, err error) {
	if len(zipString) != 18 && len(zipString) != 33 {
		return "", errors.New("源zip值长度不对")
	}
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
	switch len(zipString) {
	case 18:
		md5String = SubString(data.String(), 0, 16)
	case 33:
		md5String = SubString(data.String(), 0, 32)
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

func GetCurrentIntranetIP() string {
	ipaddress := ""
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("获取当前内网IP出错：", err)
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if strings.HasPrefix(ipnet.IP.String(), "10.") || strings.HasPrefix(ipnet.IP.String(), "192.") {
					ipaddress = ipnet.IP.String()
					break
				}
			}
		}
	}
	if len(ipaddress) < 9 {
		fmt.Println("获取当前内网IP出错：没有找到IP")
		return ""
	}
	return ipaddress
}
