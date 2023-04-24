package translation

import (
	"bufio"
	"github.com/BaoBaoGitHub/chatgpt-baobao/chatGPT/chat"
	"github.com/BaoBaoGitHub/chatgpt-baobao/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func CodeTranslateFromFile(srcPath, tgtDir, accessToken, baseURI, fileSuffix string, done func()) {
	defer done()
	// 0. chatGPT初始化
	cli := chat.NewDefaultClient(accessToken, baseURI)
	conversationID := ""
	parentMessage := ""

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
		query := "Please help me translate the following C# code into Java code:\n" + line
		//log.Println(query)
		chatText := chat.HandleChatRobustly(query, &conversationID, &parentMessage, accessToken, baseURI, cli)
		//// 3. 获取响应文件名(json文件)
		//respFilePath := utils.AddSuffix(srcPath, fileSuffix)
		//respFilePath = ModifyFileExtToJSON(respFilePath)
		// 4. 写入到响应文件中去
		utils.WriteToJSONFileFromString(targetFileName, chatText.Content, query)
	}
}

// ModifyFileExtToJSON 修改文件后缀名为JSON
func ModifyFileExtToJSON(path string) string {
	path = strings.TrimSuffix(path, filepath.Ext(path))
	path = path + ".json"
	return path
}
