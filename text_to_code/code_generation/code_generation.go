package code_generation

import (
	"github.com/BaoBaoGitHub/chatgpt-baobao/chatGPT/chat"
	"github.com/BaoBaoGitHub/chatgpt-baobao/utils"
	"math/rand"
	"path/filepath"
	"time"
)

// CodeGenerationFromFile 是代码搜索函数，从fileAddr中读取json，获取query并将结果写入到fileAddr_response中。
// 如test_shuffled_with_path_and_id_concode.json的结果会写入到如test_shuffled_with_path_and_id_concode_response.json文件中
func CodeGenerationFromFile(srcPath, tgtDir string, accessToken, baseURI string, done func()) string {
	defer done() //并发同步处理
	// chatgpt初始化
	//token := uuid.New().String()
	cli := chat.NewDefaultClient(accessToken, baseURI)
	conversationID := ""
	parentMessage := ""

	// 获取srcPath文件名，再加上response后缀，再其前面拼接tgt
	targetFileName := tgtDir + utils.AddSuffix(filepath.Base(srcPath), "response")
	//logFileName := logDir + utils.AddSuffix(filepath.Base(srcPath), "log")
	//logFileName = strings.TrimSuffix(logFileName, path.Ext(logFileName)) + ".txt"
	//logger := log.New(utils.GetFileWriter(logFileName), "", log.LstdFlags)
	// 1 打开json，获取对象
	data := utils.ReadFromJsonFile(srcPath)
	for _, content := range data {
		query := chat.GenerateQueryBasedPromts(content)
		//log.Println(query)
		//text, err := cli.GetChatText(query, conversationID, parentMessage)
		// 封装了原来的GetChatText方法，保证可以访问
		text := chat.HandleChatRobustly(query, &conversationID, &parentMessage, accessToken, baseURI, cli)
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
		// 4 结果处理
		utils.WriteToJSONFileFromString(targetFileName, text.Content, query)
		//设置连续对话
		conversationID = text.ConversationID
		parentMessage = text.MessageID
	}

	return targetFileName
}
