package main

import (
	"github.com/BaoBaoGitHub/chatgpt-baobao/text_to_code/code_search"
	"github.com/BaoBaoGitHub/chatgpt-baobao/utils"
	"github.com/google/uuid"
	"sync"
)

func main() {
	// 配置
	testFlag := false
	concurrentNum := 20                                                        //并发量
	path := "text_to_code/dataset/test_shuffled_with_path_and_id_concode.json" //文件路径
	accessToken := []string{}                                                  // chatGPT token
	baseURI := []string{}                                                      // plus会员URI
	testPath := "text_to_code/dataset/test_file.json"

	// 测试标签
	if testFlag == true {
		path = testPath
	}

	// 1. 分割源文件
	splitFilePath := utils.SplitJsonFile(path, concurrentNum)
	concurrentNum = len(splitFilePath) //split文件时，若无法恰好分割，可能会多一个文件出来

	// 2. 必须要求accessToken与baseURI长度相等，且长度等于并发量（每个并发都需要有一个token）
	tokenLen := len(accessToken)
	for i := 0; i < concurrentNum-tokenLen; i++ {
		accessToken = append(accessToken, uuid.New().String())
		baseURI = append(baseURI, "https://freechat.lidong.xin")
	}

	// 3. 并发处理代码搜索工作
	var wg sync.WaitGroup
	wg.Add(concurrentNum)

	var finalRespFilePath []string
	for i, everyPath := range splitFilePath {
		go code_search.CodeSearchFromFile(everyPath, accessToken[i%len(accessToken)], baseURI[i%len(baseURI)], wg.Done)
		finalRespFilePath = append(finalRespFilePath, utils.AddSuffix(everyPath, "response"))
	}
	// 4. 合并响应文件
	wg.Wait()
	utils.MergeJSONFile(finalRespFilePath)

	// 5. 删除中间文件(还没实现)
	defer utils.DeleteFiles(splitFilePath)
	defer utils.DeleteFiles(finalRespFilePath)
}
