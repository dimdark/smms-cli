package smms_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/antchfx/htmlquery"
)

const (
	UploadUrl       = "https://sm.ms/api/upload"
	QueryHistoryUrl = "https://sm.ms/api/list"
	ClearHistoryUrl = "https://sm.ms/api/clear"
)

// 上传图片的结果
type UploadResult struct {
	// 上传图片的状态
	Code string   `json:"code"`
	Data DataInfo `json:"data,omitempty"`
	// 接收出错信息
	Msg string `json:"msg,omitempty"`
}

// 查询上传图片的历史结果
type UploadHistoryResult struct {
	Code string     `json:"code"`
	Data []DataInfo `json:"data"`
}

// 清除图片上传历史的结果
type ClearUploadHistoryResult struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

//
type DataInfo struct {
	FileName  string `json:"filename"`
	StoreName string `json:"storename"`
	Size      int    `json:"size"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Hash      string `json:"hash"`
	DeleteUrl string `json:"delete"`
	Url       string `json:"url"`
	Path      string `json:"path"`
	Ip        string `json:"ip"`
	TimeStamp int64  `json:"timestamp"`
}

func UploadPicture(filename string) (rst UploadResult, err error) {
	rst.Code = "Error"
	rst.Msg = "Internal function error!"

	if isPicture, suffix := CheckFileSuffix(filename); !isPicture {
		errMsg := fmt.Sprintf("File has an invalid extension %s, only supports file ext: jpeg, jpg, png, gif, bmp.\n", suffix)
		err := errors.New(errMsg)
		return rst, err
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	absFilename, _ := filepath.Abs(filename)
	formFile, err := writer.CreateFormFile("smfile", absFilename)
	if err != nil {
		return rst, err
	}
	srcFile, err := os.Open(absFilename)
	if err != nil {
		return rst, err
	}
	defer srcFile.Close()
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		return rst, err

	}

	_ = writer.WriteField("ssl", "0")
	// _ = writer.WriteField("format", "json")
	contentType := writer.FormDataContentType()
	// 发送之前必须调用Close()以写入结尾行
	_ = writer.Close()

	resp, err := http.Post(UploadUrl, contentType, buf)
	if err != nil {
		return rst, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &rst)
		if err != nil {
			return rst, err
		}
	}

	rst.Msg = ""
	return rst, nil
}

func DeleteUnloadedPicture(deletePictureUrl string) string {
	doc, _ := htmlquery.LoadURL(deletePictureUrl)
	rst := htmlquery.InnerText(htmlquery.FindOne(doc, "//div[@class='container']"))
	return rst
}

func ListUploadHistory() (rst UploadHistoryResult, err error) {
	rst.Code = "error"
	resp, err := http.Get(QueryHistoryUrl)
	if err != nil {
		return rst, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &rst); err == nil {
		return rst, nil
	} else {
		return rst, err
	}
}

func ClearUploadHistory() (rst ClearUploadHistoryResult, err error) {
	rst.Code = "error"
	resp, err := http.Get(ClearHistoryUrl)
	if err != nil {
		return rst, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &rst); err == nil {
		return rst, nil
	} else {
		return rst, err
	}
}

// check file suffix
func CheckFileSuffix(filename string) (bool, string) {
	suffix := path.Ext(filename)
	if (suffix == ".jpeg") || (suffix == ".jpg") ||
		(suffix == ".png") || (suffix == ".gif") ||
		(suffix == ".bmp") {
		return true, suffix
	}
	return false, ""
}
