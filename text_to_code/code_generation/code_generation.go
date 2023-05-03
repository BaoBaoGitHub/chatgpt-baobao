package code_generation

import (
	"bufio"
	"github.com/BaoBaoGitHub/chatgpt-baobao/chatGPT/chat"
	"github.com/BaoBaoGitHub/chatgpt-baobao/utils"
	"github.com/xyhelper/chatgpt-go"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// CodeGenerationFromFile 是代码搜索函数，从fileAddr中读取json，获取query并将结果写入到fileAddr_response中。
// 如test_shuffled_with_path_and_id_concode.json的结果会写入到如test_shuffled_with_path_and_id_concode_response.json文件中
func CodeGenerationFromFile(srcPath, tgtDir, promptMode string, accessToken, baseURI string, done func()) string {
	defer done() //并发同步处理
	// chatgpt初始化
	//token := uuid.New().String()
	// 从map里面拿一个空闲的token，然后把flag写为false

	cli := chat.NewDefaultClient(accessToken, baseURI)
	conversationID := new(string)
	parentMessage := new(string)

	// 获取srcPath文件名，再加上response后缀，再其前面拼接tgt
	targetFileName := tgtDir + utils.AddSuffix(filepath.Base(srcPath), "response")
	//logFileName := logDir + utils.AddSuffix(filepath.Base(srcPath), "log")
	//logFileName = strings.TrimSuffix(logFileName, path.Ext(logFileName)) + ".txt"
	//logger := log.New(utils.GetFileWriter(logFileName), "", log.LstdFlags)
	// 1 打开json，获取对象
	data := utils.ReadFromJsonFile(srcPath)
	for _, content := range data {
		query := chat.GenerateQueryBasedPromts(content, promptMode)
		//log.Println(query)
		//text, err := cli.GetChatText(query, conversationID, parentMessage)
		// 封装了原来的GetChatText方法，保证可以访问
		text := chat.HandleChatRobustly(query, conversationID, parentMessage, accessToken, baseURI, cli)
		//设置连续对话
		//*conversationID = text.ConversationID
		//*parentMessage = text.MessageID

		resContent := text.Content
		//if strings.Count(resContent, "```") == 1 {
		//	text = chat.HandleChatRobustly("continue", conversationID, parentMessage, accessToken, baseURI, cli)
		//	*conversationID = text.ConversationID
		//	*parentMessage = text.MessageID
		//}
		//resContent = resContent + text.Content
		// 4 结果处理
		utils.WriteToJSONFileFromString(targetFileName, resContent, query)

	}

	return targetFileName
}

// CodeGenerationFromFileTokeninfoVersion 是代码搜索函数，从fileAddr中读取json，获取query并将结果写入到fileAddr_response中。
// 如test_shuffled_with_path_and_id_concode.json的结果会写入到如test_shuffled_with_path_and_id_concode_response.json文件中
func CodeGenerationFromFileTokeninfoVersion(srcPath, tgtDir, promptMode string, tokenInfo *chat.TokenInfo, done func(), paths ...string) string {
	defer done() //并发同步处理
	// 从map里面拿一个空闲的token，然后把flag写为false
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

	// 获取srcPath文件名，再加上response后缀，再其前面拼接tgt
	targetFileName := tgtDir + utils.AddSuffix(filepath.Base(srcPath), "response")
	// 1 打开json，获取对象
	data := utils.ReadFromJsonFile(srcPath)
	for _, content := range data {
		apiScanner.Scan()
		exceptionScanner.Scan()
		query := chat.GenerateQueryBasedPromts(content, promptMode, apiScanner.Text(), exceptionScanner.Text())
		// 封装了原来的GetChatText方法，保证可以访问
		//text := chat.HandleChatRobustly(query, conversationID, parentMessage, accessToken, baseURI, cli)
		text, token, uri = chat.HandleChatRobustlyTokeninfoVersion(query, conversationID, parentMessage, token, uri, tokenInfo, cli)

		resContent := text.Content
		//if strings.Count(resContent, "```") == 1 {
		//	text = chat.HandleChatRobustly("continue", conversationID, parentMessage, accessToken, baseURI, cli)
		//	*conversationID = text.ConversationID
		//	*parentMessage = text.MessageID
		//}
		//resContent = resContent + text.Content
		// 4 结果处理
		utils.WriteToJSONFileFromString(targetFileName, resContent, query)
	}
	tokenInfo.ReleaseToken(token)
	return targetFileName
}
