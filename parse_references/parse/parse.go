package parse

import (
	"bufio"
	"encoding/json"
	"github.com/BaoBaoGitHub/chatgpt-baobao/chatGPT/chat"
	"github.com/BaoBaoGitHub/chatgpt-baobao/utils"
	"github.com/xyhelper/chatgpt-go"
	"log"
	"os"
	"strings"
	"unicode"
)

const (
	APIMode       = "api"
	ExceptionMode = "exception"
)

func ParseAPI(srcPath, tgtPath string, tokenInfo *chat.TokenInfo, done func()) {
	defer done()

	f, err := os.Open(srcPath)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	//tgtFile, err := os.Create(tgtPath)
	//if err != nil {
	//	log.Panic(err)
	//}
	//defer tgtFile.Close()

	token, uri := tokenInfo.Use()
	client := chat.NewDefaultClient(token, uri)
	var text *chatgpt.ChatText

	scanner := bufio.NewScanner(f)
	//writer := bufio.NewWriter(tgtFile)
	var jsonSlice []map[string]any
	var line map[string]any
	for scanner.Scan() {
		code := scanner.Text()
		query := "List used methods with name only in the following Java methods and do not explain:\n" + code
		text, token, uri = chat.HandleChatRobustlyTokeninfoVersion(query, new(string), new(string), token, uri, tokenInfo, client)
		content := text.Content
		// 如果content中的一行是以- 开头 或者是数字. 开头，那么就是我们想要的
		line = make(map[string]any)
		line["query"] = query
		line["message"] = content
		linesInContent := strings.Split(content, "\n")
		var flag = false
		var methods []string
		for _, v := range linesInContent {
			if strings.TrimSpace(v) == "" {
				continue
			}
			v = strings.TrimSpace(v)
			if strings.HasPrefix(v, "- ") || unicode.IsDigit([]rune(v)[0]) {
				flag = true
				begin := strings.Index(v, " ") + 1
				methods = append(methods, strings.Trim(strings.TrimSpace(v[begin:]), "`"))
			}
		}
		line["code"] = methods
		line["flag"] = flag
		jsonSlice = append(jsonSlice, line)
		//log.Println(content)
	}
	utils.WriteToJSONFileFromSlice(tgtPath, jsonSlice)
}

func ParseException(srcPath, tgtPath string, tokenInfo *chat.TokenInfo, done func()) {
	defer done()

	f, err := os.Open(srcPath)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	token, uri := tokenInfo.Use()
	client := chat.NewDefaultClient(token, uri)
	var text *chatgpt.ChatText

	scanner := bufio.NewScanner(f)
	var jsonSlice []map[string]any
	var line map[string]any
	for scanner.Scan() {
		code := scanner.Text()
		query := "Does the following java code have exception handling? Only answer yes ro no and do not explain.\n" + code
		text, token, uri = chat.HandleChatRobustlyTokeninfoVersion(query, new(string), new(string), token, uri, tokenInfo, client)
		content := text.Content
		line = make(map[string]any)
		line["query"] = query
		line["message"] = content

		flag := strings.HasPrefix(strings.ToLower(content), "yes")
		line["flag"] = flag
		line["code"] = flag
		jsonSlice = append(jsonSlice, line)
		//log.Println(content)
	}
	utils.WriteToJSONFileFromSlice(tgtPath, jsonSlice)
}

// GenResFromJSONFile 从源json文件的code列中读取内容写入到tgtPath中（txt文件）
func GenResFromJSONFile(srcPath string, tgtPath string, mode string) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		log.Panic(err)
	}
	defer srcFile.Close()

	tgtFile, err := os.Create(tgtPath)
	if err != nil {
		log.Panic(err)
	}
	defer tgtFile.Close()

	scanner := bufio.NewScanner(srcFile)
	writer := bufio.NewWriter(tgtFile)
	for scanner.Scan() {
		text := scanner.Text()
		var data map[string]any
		json.Unmarshal([]byte(text), &data)

		switch mode {
		case APIMode:
			{
				if data["flag"].(bool) {
					code, ok := data["code"].([]any)
					if ok {
						var codeStrSlice []string
						for _, v := range code {
							s := v.(string)
							codeStrSlice = append(codeStrSlice, s)
						}
						s := strings.Join(codeStrSlice, ",")
						writer.WriteString(s + "\n")
					} else {
						writer.WriteString("\n")
					}
				} else {
					writer.WriteString("\n")
				}
			}
		case ExceptionMode:
			{
				flag, ok := data["code"].(bool)
				if ok {
					if flag {
						writer.WriteString("true\n")
					} else {
						writer.WriteString("false\n")
					}
				} else {
					writer.WriteString("false\n")
				}
			}
		}
	}
	writer.Flush()
}
