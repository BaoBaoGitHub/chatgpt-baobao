package translation

import (
	"bufio"
	"fmt"
	"github.com/BaoBaoGitHub/chatgpt-baobao/chatGPT/chat"
	"github.com/BaoBaoGitHub/chatgpt-baobao/utils"
	"github.com/xyhelper/chatgpt-go"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func CodeTranslateFromFile(srcPath, tgtDir, accessToken, baseURI, fileSuffix string, done func()) {
	defer done()
	// 0. chatGPT初始化
	cli := chat.NewDefaultClient(accessToken, baseURI)
	conversationID := new(string)
	parentMessage := new(string)

	// tgt文件路径
	targetFileName := tgtDir + utils.AddSuffix(filepath.Base(srcPath), "response")
	targetFileName = strings.TrimSuffix(targetFileName, path.Ext(targetFileName)) + ".json"

	// 1. 读取文件
	filePtr, err := os.Open(srcPath)
	defer filePtr.Close()
	utils.FatalCheck(err)
	scanner := bufio.NewScanner(filePtr)
	for scanner.Scan() {
		line := scanner.Text()
		// 2.查询每一行代码
		query := "Translate C# code into Java code:\n" + line
		//log.Println(query)
		chatText := chat.HandleChatRobustly(query, conversationID, parentMessage, accessToken, baseURI, cli)
		//// 3. 获取响应文件名(json文件)
		//respFilePath := utils.AddSuffix(srcPath, fileSuffix)
		//respFilePath = ModifyFileExtToJSON(respFilePath)
		// 4. 写入到响应文件中去
		utils.WriteToJSONFileFromString(targetFileName, chatText.Content, query)
	}
}

// CodeTranslateFromFileToekenInfoVersion 代码翻译的池化版本
func CodeTranslateFromFileToekenInfoVersion(srcPath, tgtDir, promptsMode string, tokenInfo *chat.TokenInfo, fileSuffix string, done func(), paths ...string) {
	defer done()
	// 0. chatGPT初始化
	token, uri := tokenInfo.Use()
	cli := chat.NewDefaultClient(token, uri)
	conversationID := new(string)
	parentMessage := new(string)
	var text *chatgpt.ChatText

	var apiScanner *bufio.Scanner
	var exceptionScanner *bufio.Scanner
	for _, path := range paths {
		if strings.Contains(path, "api") {
			f, err := os.Open(path)
			if err != nil {
				log.Panic(err)
			}
			defer f.Close()
			apiScanner = bufio.NewScanner(f)
		} else if strings.Contains(path, "exception") {
			f, err := os.Open(path)
			if err != nil {
				log.Panic(err)
			}
			defer f.Close()
			exceptionScanner = bufio.NewScanner(f)
		}
	}

	// tgt文件路径
	targetFileName := tgtDir + utils.AddSuffix(filepath.Base(srcPath), "response")
	targetFileName = strings.TrimSuffix(targetFileName, path.Ext(targetFileName)) + ".json"

	// 1. 读取文件
	filePtr, err := os.Open(srcPath)
	defer filePtr.Close()
	utils.FatalCheck(err)
	scanner := bufio.NewScanner(filePtr)
	for scanner.Scan() {
		apiScanner.Scan()
		exceptionScanner.Scan()
		line := scanner.Text()
		// 2.查询每一行代码
		query := GenQueryBasedPrompts(line, apiScanner.Text(), exceptionScanner.Text(), promptsMode)
		//log.Println(query)
		text, token, uri = chat.HandleChatRobustlyTokeninfoVersion(query, conversationID, parentMessage, token, uri, tokenInfo, cli)

		//// 3. 获取响应文件名(json文件)
		//respFilePath := utils.AddSuffix(srcPath, fileSuffix)
		//respFilePath = ModifyFileExtToJSON(respFilePath)
		// 4. 写入到响应文件中去
		utils.WriteToJSONFileFromString(targetFileName, text.Content, query)
	}

	tokenInfo.ReleaseToken(token)
}

// CodeTranslateFromFileToekenInfoVersion 代码翻译的池化版本
func CodeTranslateFromFileToekenInfoVersionWithSession(srcPath, tgtDir, promptsMode string, tokenInfo *chat.TokenInfo, fileSuffix string, done func(), paths ...string) {
	defer done()
	// 0. chatGPT初始化
	token, uri := tokenInfo.Use()
	cli := chat.NewDefaultClient(token, uri)
	conversationID := new(string)
	parentMessage := new(string)
	var text *chatgpt.ChatText

	var apiScanner *bufio.Scanner
	var exceptionScanner *bufio.Scanner
	for _, path := range paths {
		if strings.Contains(path, "api") {
			f, err := os.Open(path)
			if err != nil {
				log.Panic(err)
			}
			defer f.Close()
			apiScanner = bufio.NewScanner(f)
		} else if strings.Contains(path, "exception") {
			f, err := os.Open(path)
			if err != nil {
				log.Panic(err)
			}
			defer f.Close()
			exceptionScanner = bufio.NewScanner(f)
		}
	}

	// tgt文件路径
	targetFileName := tgtDir + utils.AddSuffix(filepath.Base(srcPath), "response")
	targetFileName = strings.TrimSuffix(targetFileName, path.Ext(targetFileName)) + ".json"

	// 1. 读取文件
	filePtr, err := os.Open(srcPath)
	defer filePtr.Close()
	utils.FatalCheck(err)
	scanner := bufio.NewScanner(filePtr)
	for scanner.Scan() {
		apiScanner.Scan()
		exceptionScanner.Scan()
		line := scanner.Text()
		// 2.查询每一行代码
		query := GenQueryBasedPrompts(line, apiScanner.Text(), exceptionScanner.Text(), promptsMode)
		//log.Println(query)
		text, token, uri = chat.HandleChatRobustlyTokeninfoVersionWithSession(query, conversationID, parentMessage, token, uri, tokenInfo, cli)
		if text != nil {
			*conversationID = text.ConversationID
			*parentMessage = text.MessageID
		} else {
			*conversationID = ""
			*parentMessage = ""
		}
		//*conversationID = text.ConversationID
		//*parentMessage = text.MessageID
		//// 3. 获取响应文件名(json文件)
		//respFilePath := utils.AddSuffix(srcPath, fileSuffix)
		//respFilePath = ModifyFileExtToJSON(respFilePath)
		// 4. 写入到响应文件中去
		utils.WriteToJSONFileFromString(targetFileName, text.Content, query)
	}

	tokenInfo.ReleaseToken(token)
}

// ModifyFileExtToJSON 修改文件后缀名为JSON
func ModifyFileExtToJSON(path string) string {
	path = strings.TrimSuffix(path, filepath.Ext(path))
	path = path + ".json"
	return path
}

// GenQueryBasedPrompts 生成query
func GenQueryBasedPrompts(code, api, exception string, promptsMode string) string {
	var res string
	switch promptsMode {
	case chat.TaskPrompts:
		{
			res = "Translate C# code into Java code:\n" + code
		}
	case chat.GuidedPromptsWithAPIAndException:
		{
			if strings.TrimSpace(api) == "" {
				api = ""
			} else {
				api = fmt.Sprintf("that calls %s", api)
			}
			if strings.TrimSpace(exception) == "true" {
				exception = ""
			} else {
				exception = "out"
			}
			res = fmt.Sprintf("Translate C# code into Java code %s with%s exception handling:\n%s", api, exception, code)
		}
	case chat.TaskPromptsWithBackticks:
		{
			res = "Translate C# code delimited by triple backticks into Java code.\nDo not provide annotation.\n" + fmt.Sprintf("```%s```", code)
		}
	case chat.TaskPromptsWithBackticksAndConciseness:
		{
			res = "Translate C# code delimited by triple backticks into concise Java code.\nDo not provide annotation.\n" + fmt.Sprintf("```%s```", code)
		}
		//TaskPromptsWithBackticksAndAnnotationAndAPI             = "task_prompts_backtick_annotation_api"
		//TaskPromptsWithBackticksAndAnnotationAndException       = "task_prompts_backtick_annotation_exception"
		//TaskPromptsWithBackticksAndAnnotationAndAPIAndException = "task_prompts_backtick_annotation_api_exception"
	case chat.TaskPromptsWithAnnotation:
		{
			res = "Translate C# code into Java code:\nDo not provide annotation.\n" + code
		}
	case chat.TaskPromptsWithBackticksAndAnnotationAndAPI:
		{
			if strings.TrimSpace(api) == "" {
				api = ""
			} else {
				api = fmt.Sprintf("that calls %s", api)
			}

			res = fmt.Sprintf("Translate C# code delimited by triple backticks into Java code %s.\nDo not provide annotation.\n", api) + fmt.Sprintf("```%s```", code)
		}
	case chat.TaskPromptsWithBackticksAndAnnotationAndException:
		{
			if strings.TrimSpace(exception) == "true" {
				exception = ""
			} else {
				exception = "out"
			}
			res = fmt.Sprintf("Translate C# code delimited by triple backticks into Java code with%s exception handling.\nD not provide annotation.\n", exception) + fmt.Sprintf("```%s```", code)
		}
	case chat.TaskPromptsWithBackticksAndAnnotationAndAPIAndException:
		{
			if strings.TrimSpace(api) == "" {
				api = ""
			} else {
				api = fmt.Sprintf("that calls %s", api)
			}
			if strings.TrimSpace(exception) == "true" {
				exception = ""
			} else {
				exception = "out"
			}
			res = fmt.Sprintf("Translate C# code delimited by triple backticks into Java code %s with%s exception handling.\nDo not provide annotation.\n", api, exception) + fmt.Sprintf("```%s```", code)
		}
	default:
		{
			res = ""
			log.Panic(promptsMode, "生成的query为空")
		}
	}
	return res
}
