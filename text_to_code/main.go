package main

import (
	"github.com/BaoBaoGitHub/chatgpt-baobao/text_to_code/code_generation"
	"github.com/BaoBaoGitHub/chatgpt-baobao/utils"
	"github.com/google/uuid"
	"path/filepath"
	"sync"
)

func main() {
	// 配置
	concurrentNum := 2        //并发量
	accessToken := []string{} // 访问Token
	baseURI := []string{}     // 代理URI

	datasetDir := "text_to_code/dataset/"
	fullPromptsDir := datasetDir + "full_prompts/"                        // 最好的prompts结果路径
	refDir := datasetDir + "ref/"                                         // 原始数据与标准答案路径
	concodePath := refDir + "test_shuffled_with_path_and_id_concode.json" // 数据源文件路径
	testConcodePath := refDir + "test_concode.json"                       //数据源测试文件
	testPath := refDir + "test.json"                                      //数据源文件路径
	testTestPath := refDir + "test_test.json"
	predictionPath := fullPromptsDir + "predictions.txt" // 预测代码部分最终存储路径
	answersPath := refDir + "answers.json"               // answers.json的路径

	// 测试标志
	if testFlag := true; testFlag {
		concodePath = testConcodePath
		testPath = testTestPath

	}

	// 0. 从concode中拿出code部分，从test.json中拿出nl部分，组成answers.json文件
	utils.GenerateAnswersFromJSONFile(concodePath, testPath, answersPath)

	// 1. 分割源文件
	splitConcodePath := utils.SplitFile(concodePath, concurrentNum)
	concurrentNum = len(splitConcodePath) //split文件时，若无法恰好分割，可能会多一个文件出来

	// 2. 必须要求accessToken与baseURI长度相等，且长度等于并发量（每个并发都需要有一个token）
	tokenLen := len(accessToken)
	for i := 0; i < concurrentNum-tokenLen; i++ {
		accessToken = append(accessToken, uuid.New().String())
		baseURI = append(baseURI, "https://freechat.lidong.xin")
	}

	// 3. 并发处理代码搜索工作
	var wg sync.WaitGroup
	wg.Add(concurrentNum)

	splitResponsePath := make([]string, concurrentNum)
	for i, srcPath := range splitConcodePath {
		go code_generation.CodeGenerationFromFile(srcPath, fullPromptsDir, accessToken[i%len(accessToken)], baseURI[i%len(baseURI)], wg.Done)
		splitResponsePath[i] = fullPromptsDir + utils.AddSuffix(filepath.Base(srcPath), "response")
	}

	// 4. 合并响应文件
	wg.Wait()
	responsePath := utils.MergeJSONFile(splitResponsePath)

	// 5. 删除中间文件
	defer utils.DeleteFiles(splitConcodePath)
	defer utils.DeleteFiles(splitResponsePath)

	// 6. 生成predictions文件
	utils.GetPredictionFromJSONFIle(responsePath, predictionPath)
}
