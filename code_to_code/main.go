package main

import (
	"github.com/BaoBaoGitHub/chatgpt-baobao/code_to_code/translation"
	"github.com/BaoBaoGitHub/chatgpt-baobao/utils"
	"github.com/google/uuid"
	"sync"
)

func main() {
	testFlag := false
	testPath := "code_to_code/dataset/test.cs"
	cs_path := "code_to_code/dataset/references.txt"
	predictionPath := "code_to_code/dataset/evaluator/predictions.txt"
	javaCodePath := "code_to_code/dataset/test.java-cs.txt.java"
	referencesPath := "code_to_code/dataset/evaluator/references.txt"
	if testFlag {
		cs_path = testPath
	}
	concurrentNum := 10
	var wg sync.WaitGroup
	var splitFilePaths []string
	var respFilePaths []string
	filePathSuffix := "response"
	uri := "https://freechat.lidong.xin"

	// 1. 拆分数据集文件
	splitFilePaths = utils.SplitFile(cs_path, concurrentNum)
	concurrentNum = len(splitFilePaths)

	// 2. 并发处理代码翻译
	for _, path := range splitFilePaths {
		wg.Add(1)
		go translation.CodeTranslateFromFile(path, uuid.New().String(), uri, filePathSuffix, wg.Done)
	}
	wg.Wait()

	// 3. 合并响应文件
	// 获取响应文件名（这里写的是真烂啊。。）
	for _, path := range splitFilePaths {
		respFilePath := utils.AddSuffix(path, filePathSuffix)
		respFilePath = translation.ModifyFileExtToJSON(respFilePath)
		respFilePaths = append(respFilePaths, respFilePath)
	}
	transitionJSONPath := utils.MergeJSONFile(respFilePaths)

	// 4. 删除中间文件
	defer utils.DeleteFiles(splitFilePaths)
	defer utils.DeleteFiles(respFilePaths)

	// 5. 生成符合评估的references.txt文件与predictions.txt文件
	utils.GenerateReferencesFromPath(javaCodePath, referencesPath)
	utils.GetPredictionFromJSONFIle(transitionJSONPath, predictionPath)
}
