package main

import (
	"github.com/BaoBaoGitHub/chatgpt-baobao/text_to_code/code_generation"
	"github.com/BaoBaoGitHub/chatgpt-baobao/utils"
	"github.com/google/uuid"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	// 配置
	concurrentNum := 20 //并发量
	accessToken := []string{
		"b721d8c0-df4c-496a-a6d1-1fe46084d3c4", "3de1b933-fc23-40ba-a40a-ec753f33ded2",
		"f3325c34-cc73-433c-8eb2-3dc75c8b274a", "5a532386-b59c-49fc-80ba-51173ad36a55",
		"db47ecb2-b016-42dd-a38e-f384269f0dd1", "0144cdef-77da-4a8b-b919-c780708df555",
		"f5beb80d-d216-4986-bae8-07ee3d8cdbee", "7a6f7d81-0dac-42fb-8e78-9c9bc872a843",
		"e07c4fe5-ae84-4769-b15e-0e892075bb48", "f08dd134-d621-4591-8467-50c3c57b853c",
		"5a2cb1fd-7d63-4226-9db2-694f13414cca", "0bbed239-b52f-40fb-97fa-fa1580e35553",
		"87ffe270-4903-4b9f-a975-41223179673a", "3aee7976-7b54-44de-87de-927c0d483508",
		"c735df01-547e-4036-bb41-ecd25e29abcb", "853ec04b-08cd-4ff6-9bc7-b763bddf15f6",
		"39ec3546-463f-458d-ba86-23763a2c5f45", "a4420768-539f-43b0-b0ef-ca97ef168d70",
		"527660af-3fbb-4430-a3c9-7c116686c254", "70fbf744-46a5-4888-85c6-32515493a12a",
	} //token
	baseURI := []string{} // 代理URI
	for i := 0; i < len(accessToken); i++ {
		baseURI = append(baseURI, "https://personalchat.lidong.xin")
	}

	datasetDir := "text_to_code/dataset/"
	fullPromptsDir := datasetDir + "full_prompts/"                        // 最好的prompts结果路径
	refDir := datasetDir + "ref/"                                         // 原始数据与标准答案路径
	concodePath := refDir + "test_shuffled_with_path_and_id_concode.json" // 数据源文件路径
	testConcodePath := refDir + "test_concode.json"                       //数据源测试文件
	//testPath := refDir + "test.json"                                      //数据源文件路径
	//testTestPath := refDir + "test_test.json"
	predictionPath := fullPromptsDir + "predictions.txt" // 预测代码部分最终存储路径
	//answersPath := refDir + "answers.json"               // answers.json的路径
	refPath := refDir + "references.txt"
	//logDir := fullPromptsDir + "log/"

	// 测试标志
	if testFlag := false; testFlag {
		concodePath = testConcodePath
	}

	////删除logDir下所有日志文件
	//err := utils.DeleteAllFiles(logDir)
	//if err != nil {
	//	panic(err)
	//}

	// 0. 从concode中拿出code部分作为references.txt
	utils.GenRefFromConcode(concodePath, refPath)

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
	splitLogPath := make([]string, concurrentNum)
	for i, srcPath := range splitConcodePath {
		go code_generation.CodeGenerationFromFile(srcPath, fullPromptsDir, accessToken[i], baseURI[i], wg.Done)
		splitResponsePath[i] = fullPromptsDir + utils.AddSuffix(filepath.Base(srcPath), "response")
		logPath := fullPromptsDir + utils.AddSuffix(filepath.Base(srcPath), "log")
		logPath = strings.TrimSuffix(logPath, path.Ext(logPath)) + ".txt"
		splitLogPath = append(splitLogPath, logPath)
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
