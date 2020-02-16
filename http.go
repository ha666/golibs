package golibs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout:     3 * time.Minute,
			TLSHandshakeTimeout: 10 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 10 * time.Minute,
				DualStack: true,
			}).DialContext,
		},
	}
}

//获取url对应的内容，返回信息：StatusCode，body，err
func Get(requestUrl string) (int, string, error) {
	reqest, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return 0, "", err
	}
	response, err := client.Do(reqest)
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()
	if err != nil {
		return 0, "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}
	return response.StatusCode, SliceByteToString(body), nil
}

// 带上Bearer Token，发起一个get请求
func GetByToken(requestUrl, token string) (int, string, error) {
	reqest, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return 0, "", err
	}
	reqest.Header.Set("Authorization", "Bearer "+token)
	response, err := client.Do(reqest)
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()
	if err != nil {
		return 0, "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}
	return response.StatusCode, SliceByteToString(body), nil
}

// 带上AuthToken，发起一个get请求
func GetByAuthToken(requestUrl, authToken string) (int, string, error) {
	reqest, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return 0, "", err
	}
	reqest.Header.Set("auth-token", authToken)
	response, err := client.Do(reqest)
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()
	if err != nil {
		return 0, "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}
	return response.StatusCode, SliceByteToString(body), nil
}

//获取url和参数列表对应的完整请求url
func BuildRequestUrl(requestUrl string, params url.Values) string {
	if len(params) <= 0 {
		return requestUrl
	}
	data := NewStringBuilder()
	data.Append(requestUrl)
	if !strings.Contains(requestUrl, "?") {
		data.Append("?")
	}
	has_param := false
	for k, v := range params {
		if len(k) > 0 && len(v[0]) > 0 {
			if has_param {
				data.Append("&")
			}
			data.Append(k)
			data.Append("=")
			data.Append(url.QueryEscape(v[0]))
			has_param = true
		}
	}
	return data.ToString()
}

//获取url对应的内容，返回信息：StatusCode，body，err
func Post(requestUrl string, params url.Values) (int, string, error) {
	reqest, err := http.NewRequest("POST", requestUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return 0, "", err
	}
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	reqest.Header.Set("User-Agent", "golang")
	response, err := client.Do(reqest)
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}
	return response.StatusCode, SliceByteToString(body), nil
}

//获取url对应的内容，返回信息：StatusCode，body，err
func PostByAuthToken(requestUrl string, authToken string, params url.Values) (int, string, error) {
	reqest, err := http.NewRequest("POST", requestUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return 0, "", err
	}
	reqest.Header.Set("auth-token", authToken)
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	reqest.Header.Set("User-Agent", "golang")
	response, err := client.Do(reqest)
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}
	return response.StatusCode, SliceByteToString(body), nil
}

//用Post方法获取url对应的内容，提交json，返回信息：StatusCode，body，err
func PostJson(requestUrl string, params map[string]string) (int, string, error) {
	req := bytes.NewBuffer(StringToSliceByte(ToJson(params)))
	reqest, err := http.NewRequest("POST", requestUrl, req)
	if err != nil {
		return 0, "", err
	}
	reqest.Header.Set("Content-Type", "application/json")
	response, err := client.Do(reqest)
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}
	return response.StatusCode, SliceByteToString(body), nil
}

// 带上Bearer Token，发起一个post请求，无内容
func PostByToken(requestUrl, token string) (int, string, error) {
	reqest, err := http.NewRequest("POST", requestUrl, nil)
	if err != nil {
		return 0, "", err
	}
	reqest.Header.Set("Authorization", "Bearer "+token)
	reqest.Header.Set("Content-Type", "application/json")
	response, err := client.Do(reqest)
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}
	return response.StatusCode, SliceByteToString(body), nil
}

// 带上Bearer Token，发起一个post请求，内容是json
func PostJsonByToken(requestUrl, token, jsonString string) (int, string, error) {
	req := bytes.NewBuffer(StringToSliceByte(jsonString))
	reqest, err := http.NewRequest("POST", requestUrl, req)
	if err != nil {
		return 0, "", err
	}
	reqest.Header.Set("Authorization", "Bearer "+token)
	reqest.Header.Set("Content-Type", "application/json")
	response, err := client.Do(reqest)
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}
	return response.StatusCode, SliceByteToString(body), nil
}

//用Post方法获取url对应的内容，提交body，返回信息：StatusCode，body，err
func PostBody(requestUrl string, reqBody string) (int, string, error) {
	req := bytes.NewBuffer(StringToSliceByte(reqBody))
	reqest, err := http.NewRequest("POST", requestUrl, req)
	if err != nil {
		return 0, "", err
	}
	reqest.Header.Set("Content-Type", "application/json")
	response, err := client.Do(reqest)
	if err != nil {
		return 0, "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}
	return response.StatusCode, SliceByteToString(body), nil
}

//获取url对应的内容,同时上传文件，返回信息：StatusCode，body，err
func PostFile(requestUrl string, params url.Values, field_name, path string) (int, string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile(field_name, path)
	if err != nil {
		return 0, "", err
	}

	//打开文件句柄操作
	fh, err := os.Open(path)
	if err != nil {
		return 0, "", err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return 0, "", err
	}

	//写入表单数据
	if len(params) > 0 {
		for k, v := range params {
			bodyWriter.WriteField(k, v[0])
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(requestUrl, contentType, bodyBuf)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}
	return resp.StatusCode, SliceByteToString(resp_body), nil
}

//获取当前连接的Http方法
func Method(r *http.Request) string {
	return r.Method
}

// 获取程序运行路径
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//解析请求参数到一个url.Values对象
func ParseRequestQuery(req string) (params url.Values, err error) {
	if Length(req) <= 0 {
		err = errors.New("没有找到参数")
		return
	}
	tmp := strings.Split(req, "&")
	if len(tmp) <= 0 {
		err = errors.New("没有找到参数")
		return
	}
	params = url.Values{}
	for i, v := range tmp {
		if v == "" {
			continue
		}
		start := strings.Index(v, "=")
		if start < 0 {
			err = errors.New(fmt.Sprintf("没有找到参数名,第%d个", i))
			return
		}
		key := SubString(v, 0, start)
		val := SubString(v, start+1, len(v)-start-1)
		if strings.Contains(val, "%") {
			val, err = url.QueryUnescape(val)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("解析参数%s的值出错:%s", key, err.Error()))
			}
		}
		params.Add(key, val)
	}
	return
}
