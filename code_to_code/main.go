package main

import (
	"github.com/BaoBaoGitHub/chatgpt-baobao/code_to_code/translation"
	"github.com/BaoBaoGitHub/chatgpt-baobao/utils"
	"github.com/google/uuid"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	concurrentNum := 5 //并发量
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

	datasetDir := "code_to_code/dataset/"
	fullPromptsDir := datasetDir + "full_prompts/" // 最好的prompts结果路径
	refDir := datasetDir + "ref/"                  // 原始数据与标准答案路径
	csPath := refDir + "test.java-cs.txt.cs"       // cs源文件的path
	javaPath := refDir + "test.java-cs.txt.java"   // java源文件的path
	testCSharpPath := refDir + "test.cs"           // test.cs源文件的path
	testJavaPath := refDir + "test.java"
	predictionPath := fullPromptsDir + "predictions.txt" //生成predictions.txt文件的path
	referencesPath := refDir + "references.txt"          // 根据javaPath生成的标准答案的path

	if testFlag := true; testFlag {
		csPath = testCSharpPath
		javaPath = testJavaPath
	}

	filePathSuffix := "response"
	uri := "https://freechat.lidong.xin"

	// 0. 根据javaPath生成references.txt文件
	utils.GenerateReferencesFromPath(javaPath, referencesPath)

	// 1. 拆分数据集文件
	splitFilePaths := utils.SplitFile(csPath, concurrentNum)
	concurrentNum = len(splitFilePaths)

	// 2. 必须要求accessToken与baseURI长度相等，且长度等于并发量（每个并发都需要有一个token）
	tokenLen := len(accessToken)
	for i := 0; i < concurrentNum-tokenLen; i++ {
		accessToken = append(accessToken, uuid.New().String())
		baseURI = append(baseURI, uri)
	}

	var wg sync.WaitGroup
	wg.Add(concurrentNum)

	splitResponsePath := make([]string, concurrentNum)
	// 2. 并发处理代码翻译
	for i, srcPath := range splitFilePaths {
		go translation.CodeTranslateFromFile(srcPath, fullPromptsDir, accessToken[i], baseURI[i], filePathSuffix, wg.Done)
		// tgt文件路径
		targetFileName := fullPromptsDir + utils.AddSuffix(filepath.Base(srcPath), "response")
		splitResponsePath[i] = strings.TrimSuffix(targetFileName, path.Ext(targetFileName)) + ".json"
	}
	wg.Wait()

	// 3. 合并响应文件
	// 获取响应文件名（这里写的是真烂啊。。）
	//for _, path := range splitFilePaths {
	//	respFilePath := utils.AddSuffix(path, filePathSuffix)
	//	respFilePath = translation.ModifyFileExtToJSON(respFilePath)
	//	respFilePaths = append(respFilePaths, respFilePath)
	//}

	transitionJSONPath := utils.MergeJSONFile(splitResponsePath)

	// 4. 删除中间文件
	defer utils.DeleteFiles(splitFilePaths)
	defer utils.DeleteFiles(splitResponsePath)

	// 5. 生成符合评估的predictions.txt文件
	utils.GetPredictionFromJSONFIle(transitionJSONPath, predictionPath)
}
