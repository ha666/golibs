package golibs

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//获取url对应的内容，返回信息：StatusCode，body，err
func Get(requestUrl string) (int, string, error) {
	response, err := http.Get(requestUrl)
	defer response.Body.Close()
	if err != nil {
		return 0, "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}
	return response.StatusCode, string(body), nil
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
	client := &http.Client{}
	reqest, err := http.NewRequest("POST", requestUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return 0, "", err
	}
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	reqest.Header.Set("User-Agent", "ha666")
	response, err := client.Do(reqest)
	if err != nil {
		return response.StatusCode, "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}
	return response.StatusCode, string(body), nil
}

//用Post方法获取url对应的内容，提交json，返回信息：StatusCode，body，err
func PostJson(requestUrl string, params map[string]string) (int, string, error) {
	client := &http.Client{}
	req := bytes.NewBuffer([]byte(ToJson(params)))
	reqest, err := http.NewRequest("POST", requestUrl, req)
	if err != nil {
		return 0, "", err
	}
	reqest.Header.Set("Content-Type", "application/json")
	response, err := client.Do(reqest)
	if err != nil {
		return response.StatusCode, "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, "", err
	}
	return response.StatusCode, string(body), nil
}

//用Post方法获取url对应的内容，提交body，返回信息：StatusCode，body，err
func PostBody(requestUrl string, reqBody string) (int, string, error) {
	client := &http.Client{}
	req := bytes.NewBuffer([]byte(reqBody))
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
	return response.StatusCode, string(body), nil
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
		return resp.StatusCode, "", err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}
	return resp.StatusCode, string(resp_body), nil
}

//获取当前连接的Http方法
func Method(r *http.Request) string {
	return r.Method
}
