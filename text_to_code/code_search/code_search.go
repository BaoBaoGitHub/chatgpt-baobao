package code_search

import (
	"github.com/BaoBaoGitHub/chatgpt-baobao/chatGPT/chat"
	"github.com/BaoBaoGitHub/chatgpt-baobao/utils"
	"github.com/xyhelper/chatgpt-go"
	"log"
	"time"
)

// CodeSearchFromFile 是代码搜索函数，从fileAddr中读取json，获取query并将结果写入到fileAddr_response中。
// 如test_shuffled_with_path_and_id_concode.json的结果会写入到如test_shuffled_with_path_and_id_concode_response.json文件中
func CodeSearchFromFile(fileAddr string, accessToken string, baseURI string, done func()) string {
	defer done() //并发同步处理
	// chatgpt初始化
	//token := uuid.New().String()
	cli := chatgpt.NewClient(
		chatgpt.WithDebug(false),
		chatgpt.WithTimeout(120*time.Second),
		chatgpt.WithAccessToken(accessToken),
		chatgpt.WithBaseURI(baseURI),
	)
	conversationID := ""
	parentMessage := ""

	targetFileName := ""
	// 1 打开json，获取对象
	data := utils.ReadFromJsonFile(fileAddr)
	for _, content := range data {
		// 2 取出nl内容
		query := content["nl"].(string)
		// 3 输入到chatgpt
		//query = "java code for \"" + query + "\"" + " The java code should be in one code block."
		query = "java code for \"" + query + "\""
		log.Println(query)
		//text, err := cli.GetChatText(query, conversationID, parentMessage)
		// 封装了原来的GetChatText方法，保证可以访问
		text := chat.HandleChatRobustly(query, &conversationID, &parentMessage, accessToken, baseURI, cli)
		// 4 结果处理
		targetFileName = utils.AddSuffix(fileAddr, "response")
		utils.WriteToJSONFileFromString(targetFileName, text.Content, query)
		//设置连续对话
		conversationID = text.ConversationID
		parentMessage = text.MessageID
	}

	return targetFileName
}
