package chat

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// token池，为了更好地处理429错误
type TokenInfo struct {
	sm         sync.Mutex
	token, uri []string
	flag       []bool // flag代表token是否空闲
	cntOf429   []int  // cnt代表429次数
}

func (r *TokenInfo) SetFlag(flag []bool) {
	r.flag = flag
}

func (r *TokenInfo) SetCntOf429(cntOf429 []int) {
	r.cntOf429 = cntOf429
}

func NewTokenInfo(token []string, uri []string) *TokenInfo {
	if len(token) != len(uri) {
		log.Fatalln("tokeninfo初始化失败因为token和uri长度不等")
	}
	flag := make([]bool, len(token))
	for i := range flag {
		flag[i] = true
	}
	cntOf429 := make([]int, len(token))
	return &TokenInfo{token: token, uri: uri, flag: flag, cntOf429: cntOf429}
}

func (r *TokenInfo) Len() int {
	if len(r.token) != len(r.uri) || len(r.token) != len(r.flag) || len(r.token) != len(r.cntOf429) {
		log.Fatalln(fmt.Sprintf("TokenInfo出错:%v", r))
	}
	return len(r.token)
}

// Use 使用一个可用的最小的cntOf429来访问
func (r *TokenInfo) Use() (string, string) {
	r.sm.Lock()
	defer r.sm.Unlock()
	index := r.getSpareIndex()
	r.flag[index] = false
	return r.token[index], r.uri[index]
}

func (r *TokenInfo) getIndexOfToken(token string) (int, bool) {
	var indexOfToken int
	for i, s := range r.token {
		if s == token {
			indexOfToken = i
			return indexOfToken, true
		}
	}
	return -1, false
}

func (r *TokenInfo) getCntOf429ForToken(token string) int {
	index, ok := r.getIndexOfToken(token)
	if !ok {
		log.Fatalln(token, "不在tokenInfo中")
		return -1
	}
	return r.cntOf429[index]
}

// Use 使用一个可用的最小的cntOf429来访问
func (r *TokenInfo) Handle429(tokenOf429 string) (string, string) {
	r.sm.Lock()
	defer r.sm.Unlock()

	var indexOfToken, ok = r.getIndexOfToken(tokenOf429)
	if !ok {
		log.Fatalln(tokenOf429, "未找到")
	}
	r.cntOf429[indexOfToken]++
	r.flag[indexOfToken] = true

	index := r.getSpareIndex()
	r.flag[index] = false
	log.Println(fmt.Sprintf("由于429或202错误accesstoken从%s切换到%s", tokenOf429, r.token[index]))
	return r.token[index], r.uri[index]
}

// getSpareIndex 从最小的可用的indexSlice中抽一个下标出来
func (r *TokenInfo) getSpareIndex() int {
	indexSlice := r.getMinValIndexSlice()
	rr := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := rr.Intn(len(indexSlice))
	return indexSlice[i]
}

// getMinValIndexSlice 获取最小的可用的indexSlice
func (r *TokenInfo) getMinValIndexSlice() []int {
	var minVal int
	var minIndexList []int
	for i, f := range r.flag {
		if f == true {
			minVal = r.cntOf429[i]
			minIndexList = []int{i}
			break
		}

	}

	// 遍历切片，查找最小值及其下标
	for i, val := range r.cntOf429 {
		if i <= minIndexList[0] {
			continue
		}
		if val < minVal && r.flag[i] == true {
			minVal = val
			minIndexList = []int{i}
		} else if val == minVal && r.flag[i] == true {
			minIndexList = append(minIndexList, i)
		}
	}

	return minIndexList
}

func (r *TokenInfo) ReleaseToken(token string) {
	r.sm.Lock()
	defer r.sm.Unlock()

	indexOfToken, ok := r.getIndexOfToken(token)
	if ok {
		r.flag[indexOfToken] = true
	} else {
		log.Fatalln(token, "没找到！")
	}
}
