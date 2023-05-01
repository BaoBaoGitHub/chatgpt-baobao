package main

import (
	"github.com/BaoBaoGitHub/chatgpt-baobao/chatGPT/chat"
	"github.com/BaoBaoGitHub/chatgpt-baobao/code_to_code/translation"
	"github.com/BaoBaoGitHub/chatgpt-baobao/utils"
	"github.com/google/uuid"
	"log"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	concurrentNum := 25 //并发量

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
		"c10fe326-e1ae-4573-90a5-33657047c25f", "d91cdae1-e2fd-4489-b47c-d7cca2fca705",
		"370020c6-188c-43e2-9c81-60206299166b", "e1a24418-d887-4b7a-ac3d-f6524eb74206",
		"2e58fc70-5695-4f1a-8556-5ae32d71bcc2", "553859b4-4491-4567-b981-ccae44f36c60",
		"6aa2c665-c45e-4591-9225-b09ec471d07a", "6133f778-2074-41f0-99db-27b0eabc3486",
		"4afa6f79-1dcf-40d5-97b1-0216ed125964", "b89f132f-7b42-4365-bb97-36c9b33edda2",
		"c1afbec3-e980-4934-9727-8be5f03874de", "992af00d-9b7c-40e6-8c8d-dafc8cc30d9a",
		"923f6ff6-9f2e-484c-a35a-b3b8421fb82d", "1c637b77-6591-4663-9d9a-5d7e73db8e96",
		"c0f4316e-3e50-4718-912b-f95152db660e", "f5339bbc-7f24-455a-8d88-769dd76c4bc7",
		"a5d5c29e-0eae-4902-8fb5-61aeeee9e0f4", "9196fbd7-d770-413d-9f40-727ac4b20e5a",
		"c3dc930b-1e86-4bec-8339-229653952bf7", "4bde5ec9-41af-49f3-9dcd-b4cb6010e463",
	}
	baseURI := []string{} // 代理URI
	for i := 0; i < len(accessToken); i++ {
		baseURI = append(baseURI, "https://p2.xyhelper.cn")
	}

	tokenInfo := chat.NewTokenInfo(accessToken, baseURI)

	//TODO 使用的是哪个prompt，一个是GuidedPromptsWithAPIAndException，一个是TaskPrompts
	promptsMode := chat.GuidedPromptsWithAPIAndException

	datasetDir := "code_to_code/dataset/"
	tgtDir := datasetDir + promptsMode + "/"     // 最好的prompts结果路径
	refDir := datasetDir + "ref/"                // 原始数据与标准答案路径
	csPath := refDir + "test.java-cs.txt.cs"     // cs源文件的path
	javaPath := refDir + "test.java-cs.txt.java" // java源文件的path
	testCSharpPath := refDir + "test.cs"         // test0.cs源文件的path
	testJavaPath := refDir + "test.java"
	predictionPath := tgtDir + "predictions.txt" //生成predictions.txt文件的path
	referencesPath := tgtDir + "references.txt"  // 根据javaPath生成的标准答案的path
	apiPath := refDir + "references_api.txt"
	testAPIPath := refDir + "test_references_api.txt"
	exceptionPath := refDir + "references_exception.txt"
	testExceptionPath := refDir + "test_references_exception.txt"

	//TODO 是否使用测试数据
	if testFlag := true; testFlag {
		csPath = testCSharpPath
		javaPath = testJavaPath
		apiPath = testAPIPath
		exceptionPath = testExceptionPath
	}

	filePathSuffix := "response"
	uri := "https://freechat.lidong.xin"

	// 0. 根据javaPath生成references.txt文件
	utils.GenerateReferencesFromPath(javaPath, referencesPath)

	// 1. 拆分数据集文件
	splitFilePaths := utils.SplitFile(csPath, concurrentNum)
	splitAPIPath := utils.SplitFile(apiPath, concurrentNum)
	splitExceptionPath := utils.SplitFile(exceptionPath, concurrentNum)

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
		//go translation.CodeTranslateFromFile(srcPath, tgtDir, accessToken[i], baseURI[i], filePathSuffix, wg.Done)
		go translation.CodeTranslateFromFileToekenInfoVersion(srcPath, tgtDir, promptsMode, tokenInfo, filePathSuffix, wg.Done, splitAPIPath[i], splitExceptionPath[i])
		// tgt文件路径
		targetFileName := tgtDir + utils.AddSuffix(filepath.Base(srcPath), "response")
		splitResponsePath[i] = strings.TrimSuffix(targetFileName, path.Ext(targetFileName)) + ".json"
	}
	wg.Wait()

	// 4. 删除中间文件
	defer utils.DeleteFiles(splitFilePaths)
	defer utils.DeleteFiles(splitResponsePath)
	defer utils.DeleteFiles(splitAPIPath)
	defer utils.DeleteFiles(splitExceptionPath)
	transitionJSONPath := utils.MergeJSONFile(splitResponsePath)

	// 5. 生成符合评估的predictions.txt文件
	utils.GetPredictionFromJSONFIle(transitionJSONPath, predictionPath)

	// 7. predictions中以类开头的百分比
	log.Println(utils.CalcClassNumFromPath(predictionPath))

	//log.Println(tokenInfo)
}
