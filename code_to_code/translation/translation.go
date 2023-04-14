package translation

import (
	"bufio"
	"github.com/BaoBaoGitHub/chatgpt-baobao/chatGPT/chat"
	"github.com/BaoBaoGitHub/chatgpt-baobao/utils"
	"github.com/xyhelper/chatgpt-go"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CodeTranslateFromFile(path, accessToken, baseURI, fileSuffix string, done func()) {
	defer done()
	// 0. chatGPT初始化
	cli := chatgpt.NewClient(
		chatgpt.WithDebug(false),
		chatgpt.WithTimeout(120*time.Second),
		chatgpt.WithAccessToken(accessToken),
		chatgpt.WithBaseURI(baseURI),
	)
	conversationID := ""
	parentMessage := ""

	// 1. 读取文件
	filePtr, err := os.Open(path)
	defer filePtr.Close()
	utils.FatalCheck(err)
	scanner := bufio.NewScanner(filePtr)
	for scanner.Scan() {
		line := scanner.Text()
		// 2.查询每一行代码
		query := "Translate following c# code surrounded ``` to java code.```" + line + "```"
		log.Println(query)
		chatText := chat.HandleChatRobustly(query, &conversationID, &parentMessage, accessToken, baseURI, cli)
		// 3. 获取响应文件名(json文件)
		respFilePath := utils.AddSuffix(path, fileSuffix)
		respFilePath = ModifyFileExtToJSON(respFilePath)
		// 4. 写入到响应文件中去
		utils.WriteToJSONFileFromString(respFilePath, chatText.Content, query)
	}
}

// ModifyFileExtToJSON 修改文件后缀名为JSON
func ModifyFileExtToJSON(path string) string {
	path = strings.TrimSuffix(path, filepath.Ext(path))
	path = path + ".json"
	return path
}
