package utils

import (
	"bytes"
	"io"
	"os"
)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
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

func DeleteFiles(path []string) {
	for _, s := range path {
		if Exists(s) {
			os.Remove(s)
		}
	}
}
