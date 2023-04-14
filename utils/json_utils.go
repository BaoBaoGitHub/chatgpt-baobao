package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Response 解析chatGPT响应为Response格式
type Response struct {
	Query   string `json:"query,omitempty"`
	Flag    bool   `json:"flag"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewSuccessfulResponse(query string, code string, message string) *Response {
	return &Response{Query: query, Flag: true, Code: code, Message: message}
}

func NewUnsuccessfulResponse(query string, message string) *Response {
	return &Response{Query: query, Flag: false, Code: "", Message: message}
}

func (r *Response) String() string {
	query := fmt.Sprintln(r.Query, "\n", strings.Repeat("*", 50))
	flag := fmt.Sprintln(r.Flag, "\n", strings.Repeat("*", 50))
	code := fmt.Sprintln(r.Code, "\n", strings.Repeat("*", 50))
	message := fmt.Sprintln(r.Message, "\n", strings.Repeat("*", 50))
	return query + flag + code + message
}

// ReadFromJsonFile 读取json文件并返回slice
// json文件中，每行一个json对象，对应一个map
func ReadFromJsonFile(fileName string) []map[string]any {
	f, err := os.Open(fileName)
	FatalCheck(err)
	defer f.Close()

	var data []map[string]any // 是一个silice，其中每一个元素均为map[string]any
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for i := 0; scanner.Scan(); i++ {
		lineData := make(map[string]any)
		json.Unmarshal([]byte(scanner.Text()), &lineData)
		data = append(data, lineData)
	}
	return data
}

// WriteToJSONFileFromString 从String解析得到struct，然后写到文件中去
func WriteToJSONFileFromString(fileName string, content string, query string) {
	// 检查文件是否存在
	if !Exists(fileName) {
		f, e := os.Create(fileName)
		f.Close()
		FatalCheck(e)
	}

	// 打开文件
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0666)
	FatalCheck(err)
	defer file.Close()

	// 转换string为json，如果string包括```，说明成功得到了code部分，此时应该返回
	resp := *ConvertStringToResponse(content, query)
	// 写入到文件中
	//log.Println(resp)
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.Encode(resp)
	file.Sync()
}

// ConvertStringToResponse 从chatGPT响应获取struct
func ConvertStringToResponse(content string, query string) *Response {
	var resp *Response
	if strings.Count(content, "```") > 1 { //可能有的响应里面只出现了一次```
		resp = NewSuccessfulResponse(query, GetCodeFromString(content), content)
	} else {
		resp = NewUnsuccessfulResponse(query, content)
	}
	return resp
}

// GetCodeFromString 从chatGPT响应解析出代码部分
func GetCodeFromString(content string) string {
	begin := strings.Index(content, "```")
	code := content[begin+3:]
	end := strings.Index(code, "```")
	code = code[:end]
	code = strings.TrimSpace(code)
	code = strings.TrimPrefix(code, "java\n") //有時候chatGPT返回結果有```java的样式（markdown的java代码段语法），需要将其去除
	return code
}

// MergeJSONFile 合并文件
func MergeJSONFile(path []string) {
	var allData []map[string]any

	for _, s := range path {
		data := ReadFromJsonFile(s)
		allData = append(allData, data...)
	}

	// 获取jsonfile文件名
	filePath := GetMergeFileName(path)
	// 将数据写入jsonfile
	WriteToJSONFileFromSlice(filePath, allData)
	log.Println("合并响应文件到：" + filePath)
}

func GetMergeFileName(path []string) string {
	// 获取后缀名
	s := path[0]
	ext := filepath.Ext(s)
	// 去除后缀名
	nameWithoutExt := strings.TrimSuffix(s, ext)
	// 获取最后两个下划线所在位置，删除其后元素
	twoIndex := LastTwoIndex(nameWithoutExt, "_")
	// 拼接最后一个下划线之前的+后缀
	res := nameWithoutExt[:twoIndex[0]] + nameWithoutExt[twoIndex[1]:] + ext

	return res
}

func LastTwoIndex(str string, subStr string) []int {

	var indices []int
	startIndex := len(str)

	for i := 0; i < 2; i++ {
		// 查找最后一个匹配子字符串的下标
		index := strings.LastIndex(str[:startIndex], subStr)

		// 如果没有找到匹配子字符串，则退出循环
		if index == -1 {
			break
		}

		// 记录匹配子字符串的下标
		indices = append(indices, index)

		// 更新下一次查找匹配子字符串的起始位置
		startIndex = index
	}

	// 反转下标顺序，使其按原始顺序排列
	indices[0], indices[1] = indices[1], indices[0]

	return indices

}

// WriteToJSONFileFromSlice
func WriteToJSONFileFromSlice(fileName string, data []map[string]any) {
	// 检查文件是否存在
	if !Exists(fileName) {
		f, e := os.Create(fileName)
		f.Close()
		FatalCheck(e)
	}

	// 打开文件
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0666)
	FatalCheck(err)
	defer file.Close()

	// 写入到文件中
	//fmt.Println(data)
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	for _, line := range data {
		encoder.Encode(line)
	}
	file.Sync()
}
