package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsNotExist(err) {
			//log.Println(path, "不存在！")
			return false
		} else {
			//log.Panic(err)
			return false
		}
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// lineCounter 根据文件中的\n计算文件行数，注意最后一行应该也有\n才行
func LineCounter(fileName string) (int, error) {
	r, err := os.Open(fileName)
	FatalCheck(err)
	defer r.Close()

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// SplitFile 将文件分割为指定的数量（如果不能整除，则分割为#{num}+1个文件），序号从0开始。
// 返回值为分割后文件的路径。
func SplitFile(fileName string, num int) []string {
	// 如果文件行数不能被恰好分割为行数相等的num个，则多一个文件
	lines, err := LineCounter(fileName)
	FatalCheck(err)
	linesInEveryFile := lines / num
	if lines%linesInEveryFile != 0 {
		num++
	}
	var res []string

	// 分割文件
	f, err := os.Open(fileName)
	FatalCheck(err)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	newFileName := fileName
	var file *os.File
	for i := 0; scanner.Scan(); i++ {
		if i%linesInEveryFile == 0 {
			// 添加文件后缀
			newFileName = AddSuffix(fileName, i/linesInEveryFile)
			// 关闭上一个文件
			if file != nil {
				file.Close()
			}
			// 创建一个文件
			tmpNewFile, err := os.Create(newFileName)
			tmpNewFile.Close()
			FatalCheck(err)
			res = append(res, newFileName)
			file, err = os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, 0666)
			FatalCheck(err)
		}
		//copy一行
		file.WriteString(scanner.Text() + "\n")
		file.Sync()
	}
	return res
}

// AddSuffix 为源文件添加后缀s。
// 如源文件名为test.txt,返回test_s.txt
func AddSuffix(fileName string, s any) string {
	ext := filepath.Ext(fileName)
	prefix := strings.TrimSuffix(fileName, ext)
	var res string

	switch s.(type) {
	case int:
		res = prefix + "_" + fmt.Sprintf("%d", s) + ext
	case string:
		res = prefix + "_" + fmt.Sprintf("%s", s) + ext
	default:
		res = ""
	}
	return res
}

func DeleteFiles(path []string) {
	for _, s := range path {
		if Exists(s) {
			//log.Println(s, "存在，调用删除函数")
			err := os.Remove(s)
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Println(s, "不存在，没有调用remove函数")
		}
	}
}

func GenerateReferencesFromPath(sourcePath string, destPath string) {
	// 打开源文件
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	if Exists(destPath) {
		os.Remove(destPath)
	}

	// 创建目标文件
	dstFile, err := os.Create(destPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	// 拷贝文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}

}

// getFileWriter 返回一个文件写入器
func GetFileWriter(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	return f
}

func DeleteAllFiles(dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// CalcClassNumFromPath 计算文件中不符合要求的行数
func CalcClassNumFromPath(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return "100%"
	}
	defer f.Close()

	lineCnt := 0
	// 创建Scanner
	scanner := bufio.NewScanner(f)
	cnt := 0
	// 逐行读取文件
	for scanner.Scan() {
		line := scanner.Text()
		lineCnt++
		if strings.HasPrefix(line, "class") || strings.HasPrefix(line, "public class") ||
			strings.HasPrefix(line, "import") || strings.HasPrefix(line, "``` import") {
			cnt++
		}
	}
	return fmt.Sprintf("%.2f", 100.0*float64(cnt)/float64(lineCnt)) + "%"
}

// AddSpace 从srcPath中读取每行，并在适当位置添加空格
func AddSpace(srcPath, dstPath string) {
	srcFile, err := os.Open(srcPath)
	FatalCheck(err)
	defer srcFile.Close()
	scanner := bufio.NewScanner(srcFile)

	dstFile, err := os.Create(dstPath)
	FatalCheck(err)
	defer dstFile.Close()
	writer := bufio.NewWriter(dstFile)

	for scanner.Scan() {
		line := scanner.Text()
		res := parseSpace(line)
		writer.WriteString(res + "\n")
	}
	writer.Flush()
}

// parseSpace 为s中适当位置添加空格
func parseSpace(s string) string {
	s = strings.ReplaceAll(s, "(", " ( ")
	s = strings.ReplaceAll(s, ")", " ) ")
	s = strings.ReplaceAll(s, ".", " . ")
	s = strings.ReplaceAll(s, ",", " , ")
	s = strings.ReplaceAll(s, "<", " < ")
	s = strings.ReplaceAll(s, ">", " > ")
	s = strings.ReplaceAll(s, ";", " ; ")
	s = strings.ReplaceAll(s, "[", " [ ")
	s = strings.ReplaceAll(s, "]", " ] ")
	s = strings.ReplaceAll(s, "++", " ++ ")
	s = strings.ReplaceAll(s, "--", " -- ")
	s = strings.ReplaceAll(s, "@", " @ ")
	s = strings.ReplaceAll(s, "...", " ... ")
	return s
}

func deleteSpace(s string) string {
	s = strings.ReplaceAll(s, " ( ", "(")
	s = strings.ReplaceAll(s, " ) ", ")")
	s = strings.ReplaceAll(s, " . ", ".")
	s = strings.ReplaceAll(s, " , ", ",")
	s = strings.ReplaceAll(s, " < ", "<")
	s = strings.ReplaceAll(s, " > ", ">")
	s = strings.ReplaceAll(s, " ; ", ";")
	s = strings.ReplaceAll(s, " [ ", "[")
	s = strings.ReplaceAll(s, " ] ", "]")
	s = strings.ReplaceAll(s, " ++ ", "++")
	s = strings.ReplaceAll(s, " -- ", "--")
	s = strings.ReplaceAll(s, " @ ", "@")
	s = strings.ReplaceAll(s, " ... ", "...")
	return s
}

func deleteOverride(srcPath string, dstPath string) {
	srcFile, err := os.Open(srcPath)
	FatalCheck(err)
	defer srcFile.Close()
	scanner := bufio.NewScanner(srcFile)

	dstFile, err := os.Create(dstPath)
	FatalCheck(err)
	defer dstFile.Close()
	writer := bufio.NewWriter(dstFile)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "@Override ", "")
		writer.WriteString(line + "\n")
	}

	writer.Flush()
}

// GetPredictionFromJSONFIle 修改json文件的格式，只保留代码部分，且为txt格式。将\n改为空格
func GetPredictionWithoutCommentsFromJSONFIle(sourcePath string, tgtDir string) string {
	contentSlice := ReadFromJsonFile(sourcePath)

	tgtPath := filepath.Join(tgtDir, "predictions_without_comments.txt")
	// 打开文件
	file, err := os.Create(tgtPath)
	FatalCheck(err)
	writer := bufio.NewWriter(file)
	defer file.Close()

	var code string
	for _, content := range contentSlice {
		code = ""
		if content["flag"].(bool) == true {
			code = fmt.Sprintf("%s", content["code"])
		} else if content["flag"].(bool) == false {
			message := content["message"].(string)
			if strings.Contains(message, "\n\npublic") && strings.Contains(message, "\n}\n\n") {
				begin := strings.Index(message, "\n\npublic")
				end := strings.Index(message, "\n}\n\n") + 4
				code = message[begin:end]
			} else {
				code = message
			}
		} else {
			panic("flag不存在！")
		}
		// 删掉注释，修改代码格式
		code = removeComments(code)
		code = ModifyCodeFormat(code)
		file.WriteString(strings.TrimSpace(code) + "\n")
	}
	writer.Flush()
	return tgtPath
}

// GetPredictionWithoutCommentsWithSpaceFromJSONFile json中拿出代码，然后创建一个删除注释并且添加了空格的目标文件
func GetPredictionWithoutCommentsWithSpaceFromJSONFile(srcPath, tgtDir string) {
	predictionWithoutCommentsPath := GetPredictionWithoutCommentsFromJSONFIle(srcPath, tgtDir)
	AddSpace(predictionWithoutCommentsPath, AddSuffix(predictionWithoutCommentsPath, "with_space"))
}

// removeComments 从code里面remove掉单行注释
func removeComments(code string) string {
	var codeWithoutComments []string
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "//") {
			codeWithoutComments = append(codeWithoutComments, line)
		}
	}
	return strings.Join(codeWithoutComments, "\n")
}
