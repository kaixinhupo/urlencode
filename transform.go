package urlencode

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"sort"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	utf8s   = "utf-8,utf8,UTF8,UTF-8"
	chinese = "gb2312,Gb2312,GB2312,gbk,GBK,Gbk,gb18030,GB18030,Gb18030"
)

var (
	rune2byte = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
)

// 把Http请求参数按照指定编码进行UrlEncode编码
//
//@Param encoding 支持uft8 gbk gb2312 gb18030编码方式
func UrlEncode(params map[string]string, encoding string) (string, error) {
	if params == nil {
		return "", nil
	}
	isChinese := strings.Contains(chinese, encoding)
	isUtf8 := strings.Contains(utf8s, encoding)

	if !(isChinese || isUtf8) {
		return "", errors.New("unrecognized encoding")
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	builder := strings.Builder{}
	for _, k := range keys {
		builder.WriteString(encode(k, isChinese))
		builder.WriteByte('=')
		builder.WriteString(encode(params[k], isChinese))
		builder.WriteByte('&')
	}
	str := builder.String()
	str = strings.TrimSuffix(str, "&")

	return str, nil
}

// 解码UrlEncode编码的字符串
//
// @Param encoding 支持uft8 gbk gb2312 gb18030编码方式
func UrlDecode(str string, encoding string) (map[string]string, error) {
	if str == "" {
		return nil, nil
	}
	isChinese := strings.Contains(chinese, encoding)
	isUtf8 := strings.Contains(utf8s, encoding)

	if !(isChinese || isUtf8) {
		return nil, errors.New("unrecognized encoding")
	}

	parts := strings.Split(str, "&")
	rst := make(map[string]string, len(parts))
	for _, p := range parts {
		kv := strings.Split(p, "=")
		key, _ := decode(kv[0], isChinese)
		if len(kv) > 1 {
			val, _ := decode(kv[1], isChinese)
			rst[key] = val
		} else {
			rst[key] = ""
		}
	}

	return rst, nil
}

// 对输入字符出按照指定编码进行UrlEncode处理
//
// @Param encoding 支持uft8 gbk gb2312 gb18030编码方式
func Encode(str string, encoding string) (string, error) {
	if str == "" {
		return "", nil
	}
	isChinese := strings.Contains(chinese, encoding)
	isUtf8 := strings.Contains(utf8s, encoding)

	if !(isChinese || isUtf8) {
		return "", errors.New("unrecognized encoding")
	}
	return encode(str, isChinese), nil
}

// 解码UrlEncode编码的字符串
//
// @Param encoding 支持uft8 gbk gb2312 gb18030编码方式
func Decode(str string, encoding string) (string, error) {
	if str == "" {
		return "", nil
	}
	isChinese := strings.Contains(chinese, encoding)
	isUtf8 := strings.Contains(utf8s, encoding)

	if !(isChinese || isUtf8) {
		return "", errors.New("unrecognized encoding")
	}
	return decode(str, isChinese)
}

func decode(str string, isChinese bool) (string, error) {
	runes := []rune(str)
	length := len(runes)
	buf := new(bytes.Buffer)
	i := 0
	for i < length {
		r := runes[i]
		if r != '%' {
			buf.WriteRune(r)
			i++
		} else {
			b, err := parseEncodeByte(runes[i+1], runes[i+2])
			if err != nil {
				return "", err
			}
			buf.WriteByte(b)
			i += 3
		}
	}
	if isChinese {
		data, _ := ioutil.ReadAll(transform.NewReader(buf, simplifiedchinese.GB18030.NewDecoder()))
		return string(data), nil
	} else {
		return buf.String(), nil
	}
}

func encode(str string, isChinese bool) string {
	builder := strings.Builder{}
	for _, r := range str {

		if !shouldEscape(r) {
			builder.WriteRune(r)
			continue
		}

		p := make([]byte, utf8.UTFMax)
		c := utf8.EncodeRune(p, r)
		var data []byte
		if isChinese {
			data, _ = ioutil.ReadAll(transform.NewReader(bytes.NewReader(p[0:c]), simplifiedchinese.GB18030.NewEncoder()))
		} else {
			data = p[0:c]
		}

		target := make([]byte, len(data)*2)
		hex.Encode(target, data)

		i := 0
		for ; i < len(target); i += 2 {
			builder.WriteByte('%')
			builder.Write([]byte{target[i], target[i+1]})
		}
	}
	return builder.String()
}

func shouldEscape(r rune) bool {
	return !((r >= '0' && r <= '9') ||
		(r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		r == '-' || r == '_' || r == '.')
}

func parseEncodeByte(high rune, low rune) (byte, error) {
	if high >= 65 && high <= 90 {
		high = high + 32
	}
	if low >= 65 && low <= 90 {
		low = low + 32
	}
	i := byteOfRune(high)
	if i == -1 {
		return 0, errors.New("invalid arguments")
	}
	j := byteOfRune(low)

	if j == -1 {
		return 0, errors.New("invalid arguments")
	}
	return byte((0x0f&i)<<4 + (0x0f & j)), nil
}

func byteOfRune(r rune) int {
	for i, b := range rune2byte {
		if r == rune(b) {
			return i
		}
	}
	return -1
}
