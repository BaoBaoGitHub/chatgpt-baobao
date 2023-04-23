package chat

import (
	"github.com/google/uuid"
	"github.com/xyhelper/chatgpt-go"
	"log"
	"math/rand"
	"strings"
	"time"
)

// HandleChatRobustly 若chatgpt访问出错，则新建一个chatgpt连接，以保证chatgpt健壮访问
func HandleChatRobustly(
	query string,
	conversationIDPtr *string,
	parentMessagePtr *string,
	accessToken string,
	baseURI string,
	cli *chatgpt.Client) *chatgpt.ChatText {

	var text *chatgpt.ChatText
	//如果panic，就新建chatGPT对话
	defer func() {
		if r := recover(); r != nil {
			rand.Seed(time.Now().UnixNano())
			time.Sleep(time.Duration(rand.Intn(60)+60) * time.Second)
			text = HandleError(query, conversationIDPtr, parentMessagePtr, accessToken, baseURI, cli, r)
		}
	}()
	// 正常访问chatgpt
	text, err := cli.GetChatText(query, *conversationIDPtr, *parentMessagePtr)
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(5)+5) * time.Second)
	// 如果出错，就新建chatGPT对话
	if err != nil {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(30)+30) * time.Second)
		// 从免费池新建一个
		text = HandleError(query, conversationIDPtr, parentMessagePtr, accessToken, baseURI, cli, err)
	}
	return text
}

func HandleError(
	query string,
	conversationIDPtr *string,
	parentMessagePtr *string,
	accessToken string,
	baseURI string,
	cli *chatgpt.Client,
	err any) *chatgpt.ChatText {
	log.Println(accessToken, err)
	if err2, ok := err.(error); ok {
		//log.Println(accessToken, err2)
		if strings.Contains(err2.Error(), "429") {
			cli = NewDefaultClient(uuid.New().String(), "https://freechat.lidong.xin")
		} else {
			//log.Println(accessToken,baseURI)
			cli = NewDefaultClient(accessToken, baseURI)
		}
	} else {
		//log.Println(accessToken, err)
		// 新建cli
		cli = NewDefaultClient(accessToken, baseURI)
	}
	// 修改conversationID和parentMessage
	*conversationIDPtr = ""
	*parentMessagePtr = ""
	// 再次访问并返回结果
	text, err := cli.GetChatText(query, *conversationIDPtr, *parentMessagePtr)
	// 如果仍然出错，递归解决错误
	if err != nil {
		return HandleError(query, conversationIDPtr, parentMessagePtr, accessToken, baseURI, cli, err)
	}
	return text
}

func NewDefaultClient(accessToken, baseURI string) *chatgpt.Client {
	return chatgpt.NewClient(
		chatgpt.WithDebug(false),
		chatgpt.WithTimeout(180*time.Second),
		chatgpt.WithAccessToken(accessToken),
		chatgpt.WithBaseURI(baseURI),
		chatgpt.WithModel("text-davinci-002-render-sha"),
	)
}
