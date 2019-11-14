package golibs

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	math_rand "math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var (
	GetFileSuffixRegex         = regexp.MustCompile(".(jpg|jpeg|png|gif|exe|doc|docx|ppt|pptx|xls|xlsx)")
	IsTaobaoNickRegex          = regexp.MustCompile(`(^[\\u4e00-\\u9fa5\\w_—\\-，。…·〔〕（）！@￥%……&*？、；‘“]*$)`)
	IsSubTaobaoNickRegex       = regexp.MustCompile(`(^[\\u4e00-\\u9fa5\\w_—\\-，。…·〔〕（）！@￥%……&*？、；‘“:]*$)`)
	IsVersionRegex             = regexp.MustCompile(`(^[0-9.]*$)`)
	IsUrlRegex                 = regexp.MustCompile(`(^[a-zA-z]+://[^\s]*$)`)
	IsNumberRegex              = regexp.MustCompile(`(^[0-9]*$)`)
	IsAsciiRegex               = regexp.MustCompile(`(^[\x00-\xff]*$)`)
	IsMultipNumberRegex        = regexp.MustCompile(`(^[0-9,]*$)`)
	IsLetterOrNumberRegex      = regexp.MustCompile(`(^[A-Za-z0-9_]*$)`)
	IsLetterOrNumber1Regex     = regexp.MustCompile(`(^[A-Za-z0-9_-]*$)`)
	IsHanOrLetterOrNumberRegex = regexp.MustCompile("^[A-Za-z0-9_\u4e00-\u9fa5-]*$")
	IsGeneralStringRegex       = regexp.MustCompile("^[A-Za-z0-9_\\-#+./:\u4e00-\u9fa5]*$")
	IsStandardTimeRegex        = regexp.MustCompile(`^[1-9]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])\s+(20|21|22|23|[0-1]\d):[0-5]\d:[0-5]\d$`)
	IsIPAddressRegecx          = regexp.MustCompile(`^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}$`)
	IsIntranetIPRegex          = regexp.MustCompile(`^(127\.0\.0\.1)|(10\.\d{1,3}\.\d{1,3}\.\d{1,3})|(172\.((1[6-9])|(2\d)|(3[01]))\.\d{1,3}\.\d{1,3})|(192\.168\.\d{1,3}\.\d{1,3})$`)
	IsEmailRegex               = regexp.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
	IsMobileRegex              = regexp.MustCompile("^[1](([3][0-9])|([4][5-9])|([5][0-3,5-9])|([6][2,5,6,7])|([7][0-8])|([8][0-9])|([9][1,3,5,8,9]))[0-9]{8}$")
)

//使用 utf8.RuneCountInString()统计字符串长度
func Length(str string) int {
	return utf8.RuneCountInString(str)
}

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
	suffix := GetFileSuffixRegex.ReplaceAllString(s, "")
	return suffix
}

//生成指定范围内的int64数字
func RandInt64(min, max int64) int64 {
	max_bigint := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, max_bigint)
	if i.Int64() < min {
		RandInt64(min, max)
	}
	return i.Int64()
}

//删除空格、换行、空格等字符
func Strim(s string) string {
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	return s
}

//字符串转成Unicode编码
func String2Unicode(s string) string {
	data := NewStringBuilder()
	for _, r := range s {
		rint := int64(r)
		if rint < 128 {
			data.Append("\\u00").Append(strconv.FormatInt(rint, 16))
		} else {
			data.Append("\\u").Append(strconv.FormatInt(rint, 16))
		}
	}
	return data.ToString()
}

//Unicode编码转成字符串
func Unicode2String(s string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(s, `\u`, ``, -1))
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
func HTMLEncode(s string) string {
	html := ""
	for _, r := range s {
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
func GetIPNums(s string) (ipNum uint32, err error) {
	if strings.EqualFold(s, "") {
		return ipNum, errors.New("ipAddress is null")
	}
	items := strings.Split(s, ".")
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

//获取IP地址，不包含端口号
func GetIPAddressNotPort(ip_address string) string {
	if !strings.Contains(ip_address, ":") {
		return ip_address
	}
	start := strings.Index(ip_address, ":")
	if start <= 2 {
		return ip_address
	}
	return SubString(ip_address, 0, start)
}

//判断是否是淘宝用户名
func IsTaobaoNick(s string) bool {
	if len(s) < 2 {
		return false
	}
	return IsTaobaoNickRegex.MatchString(s)
}

//判断是否是淘宝用户名（子帐号）
func IsSubTaobaoNick(s string) bool {
	if len(s) < 2 {
		return false
	}
	return IsSubTaobaoNickRegex.MatchString(s)
}

//判断是否版本号
func IsVersion(s string) bool {
	if len(s) < 1 {
		return false
	}
	return IsVersionRegex.MatchString(s)
}

//判断是否网址
func IsUrl(s string) bool {
	if len(s) < 1 {
		return false
	}
	return IsUrlRegex.MatchString(s)
}

//是否数字
func IsNumber(s string) bool {
	if len(s) < 1 {
		return false
	}
	return IsNumberRegex.MatchString(s)
}

//是否Ascii字符
func IsAscii(s string) bool {
	if len(s) < 1 {
		return false
	}
	return IsAsciiRegex.MatchString(s)
}

//是否多数字(用逗号间隔)
func IsMultipNumber(s string) bool {
	if len(s) < 1 {
		return false
	}
	return IsMultipNumberRegex.MatchString(s)
}

//判断是否由字母、数字、下划线组成
func IsLetterOrNumber(s string) bool {
	if len(s) < 1 {
		return false
	}
	return IsLetterOrNumberRegex.MatchString(s)
}

//判断是否由字母、数字、下划线组成
func IsLetterOrNumber1(s string) bool {
	if len(s) < 1 {
		return false
	}
	return IsLetterOrNumber1Regex.MatchString(s)
}

//判断是否由汉字、字母、数字、下划线组成
func IsHanOrLetterOrNumber(s string) bool {
	if len(s) < 1 {
		return false
	}
	return IsHanOrLetterOrNumberRegex.MatchString(s)
}

//判断是否由汉字、字母、数字、下划线、中杠等组成
func IsGeneralString(s string) bool {
	if len(s) < 1 {
		return false
	}
	return IsGeneralStringRegex.MatchString(s)
}

//判断是否标准时间格式
func IsStandardTime(s string) bool {
	if len(s) != 19 {
		return false
	}
	return IsStandardTimeRegex.MatchString(s)
}

// 是否IPv4地址
func IsIPAddress(s string) bool {
	if len(s) < 7 {
		return false
	}
	return IsIPAddressRegecx.MatchString(s)
}

// 是否内网IP地址
func IsIntranetIP(s string) bool {
	if len(s) < 7 {
		return false
	}
	return IsIntranetIPRegex.MatchString(s)
}

//是否email
func IsEmail(s string) bool {
	if len(s) < 1 {
		return false
	}
	return IsEmailRegex.MatchString(s)
}

//是否手机号
func IsMobile(s string) bool {
	if len(s) < 1 {
		return false
	}
	return IsMobileRegex.MatchString(s)
}

/*
判断字符串是否全中文字符
*/
func IsAllChineseChar(s string) bool {
	for _, r := range s {
		if !unicode.Is(unicode.Scripts["Han"], r) {
			return false
		}
	}
	return true
}

//是否utf-8编码字符串
func IsUtf8(s string) bool {
	count := 0
	for _, v := range s {
		if int(v) > 65530 {
			count++
		}
	}
	return count == 0
}

const zip_offset int = 19968

//压缩md5或guid
func ZipMd5(md5String string) (zip_string string, err error) {
	if len(md5String) != 16 && len(md5String) != 32 {
		return "", errors.New("源md5值长度不对")
	}
	var md5_bytes []byte
	switch len(md5String) {
	case 16:
		md5_bytes = getHexBytes(md5String + "00")
	case 32:
		md5_bytes = getHexBytes(md5String + "0")
	}
	var data bytes.Buffer
	var total int = 0
	for index := 0; index < len(md5_bytes); index++ {
		switch index % 3 {
		case 0:
			intByte, err := strconv.Atoi(fmt.Sprintf("%d", md5_bytes[index]))
			if err != nil {
				return zip_string, err
			}
			total = intByte
			break
		case 1:
			intByte, err := strconv.Atoi(fmt.Sprintf("%d", md5_bytes[index]))
			if err != nil {
				return zip_string, err
			}
			total += intByte << 4
			break
		case 2:
			intByte, err := strconv.Atoi(fmt.Sprintf("%d", md5_bytes[index]))
			if err != nil {
				return zip_string, err
			}
			total += (intByte << 8) + zip_offset
			uni_string := fmt.Sprintf("\\u%x", total)
			chinese_string, err := Unicode2String(uni_string)
			if err != nil {
				return zip_string, err
			}
			data.WriteString(chinese_string)
			total = 0
			break
		}
	}
	zip_string = data.String()
	return
}

//解压缩汉字，结果为guid或md5值
func UnZipMd5(zip_string string) (md5_string string, err error) {
	if len(zip_string) != 18 && len(zip_string) != 33 {
		return "", errors.New("源zip值长度不对")
	}
	var data bytes.Buffer
	unicode_string := String2Unicode(zip_string)
	unicode_strings := strings.Split(unicode_string, "\\u")
	for _, unicode_alue := range unicode_strings {
		if strings.EqualFold(unicode_alue, "") {
			continue
		}
		dec, err := strconv.ParseInt(unicode_alue, 16, 32)
		if err != nil {
			return md5_string, err
		}
		dec -= int64(zip_offset)
		data.WriteString(ten_value_to_char(dec & 15))
		data.WriteString(ten_value_to_char((dec >> 4) & 15))
		data.WriteString(ten_value_to_char((dec >> 8) & 15))
	}
	switch len(zip_string) {
	case 18:
		md5_string = SubString(data.String(), 0, 16)
	case 33:
		md5_string = SubString(data.String(), 0, 32)
	}
	return
}

func getHexBytes(str string) []byte {
	result := StringToSliceByte(str)
	for index := 0; index < len(result); index++ {
		if result[index] < 58 {
			result[index] -= 48
		} else {
			result[index] -= 55
		}
	}
	return result
}

func ten_value_to_char(ten int64) string {
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

//获取当前内网IP
func GetCurrentIntranetIP() string {
	ip_address := ""
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("获取当前内网IP出错：", err)
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if strings.HasPrefix(ipnet.IP.String(), "10.") || strings.HasPrefix(ipnet.IP.String(), "192.") || strings.HasPrefix(ipnet.IP.String(), "172.") {
					ip_address = ipnet.IP.String()
					break
				}
			}
		}
	}
	if len(ip_address) < 9 {
		fmt.Println("获取当前内网IP出错：没有找到IP")
		return ""
	}
	return ip_address
}

//序列化为json
func ToJson(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}

//格式化为Json字符串
func FormatJson(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return err.Error()
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "\t")
	if err != nil {
		return err.Error()
	}
	return out.String()
}

//序列化-->zlib压缩-->混淆-->base64
func ToJsonZipConfusedBase64(obj interface{}) string {
	b, _ := json.Marshal(obj)
	b, _ = ZlibZipBytes(b)
	b = ConfusedTwo(b)
	return Base64(b)
}

//序列化-->混淆-->base64
func ToJsonConfusedBase64(obj interface{}) string {
	b, _ := json.Marshal(obj)
	b = ConfusedTwo(b)
	return Base64(b)
}

//混淆-->zlib压缩-->base64
func ToConfusedZipBase64(str string) string {
	b := ConfusedTwo(StringToSliceByte(str))
	b, _ = ZlibZipBytes(b)
	return Base64(b)
}

//混淆-->base64
func ToConfusedBase64(str string) string {
	b := ConfusedTwo(StringToSliceByte(str))
	return Base64(b)
}

//生成简单随机密码，短时间内会重复
func GetSimplePwd(lenth int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	runes := []rune(chars)
	len1 := len(runes)
	time.Sleep(time.Millisecond)
	math_rand.Seed(time.Now().UnixNano())
	pwd := NewStringBuilder()
	for index := 0; index < lenth; index++ {
		pwd.Append(string(runes[math_rand.Intn(len1)]))
	}
	return pwd.ToString()
}

//生成随机密码，短时间内会重复
func GetPwd(lenth int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%&*+-=?@^_~'"
	runes := []rune(chars)
	len1 := len(runes)
	time.Sleep(time.Millisecond)
	math_rand.Seed(time.Now().UnixNano())
	pwd := NewStringBuilder()
	for index := 0; index < lenth; index++ {
		pwd.Append(string(runes[math_rand.Intn(len1)]))
	}
	return pwd.ToString()
}
