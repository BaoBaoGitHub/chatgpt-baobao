package chat

import (
	"errors"
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
	//text, err := cli.GetChatText(query, *conversationIDPtr, *parentMessagePtr)
	text, err := cli.GetChatText(query)
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(15)+15) * time.Second)
	// 如果出错，就新建chatGPT对话
	if err != nil {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(30)+30) * time.Second)
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
			//cli = NewDefaultClient(uuid.New().String(), "https://freechat.xyhelper.cn")	//境外服务器
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
	//text, err := cli.GetChatText(query, *conversationIDPtr, *parentMessagePtr)
	text, err := cli.GetChatText(query)
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

// HandleChatRobustlyTokeninfoVersion 若chatgpt访问出错，则新建一个chatgpt连接，以保证chatgpt健壮访问
func HandleChatRobustlyTokeninfoVersion(
	query string,
	conversationIDPtr *string,
	parentMessagePtr *string,
	accessToken, baseURI string,
	tokenInfo *TokenInfo,
	cli *chatgpt.Client) (*chatgpt.ChatText, string, string) {

	var text *chatgpt.ChatText
	var (
		newToken string = accessToken
		newURI   string = baseURI
	)
	//如果panic，就新建chatGPT对话
	defer func() {
		if r := recover(); r != nil {
			rand.Seed(time.Now().UnixNano())
			time.Sleep(time.Duration(rand.Intn(60)+60) * time.Second)
			text, newToken, newURI = HandleErrorTokeninfoVersion(query, conversationIDPtr, parentMessagePtr, accessToken, baseURI, tokenInfo, cli, r)
		}
	}()
	// 正常访问chatgpt
	//text, err := cli.GetChatText(query, *conversationIDPtr, *parentMessagePtr)

	text, err := cli.GetChatText(query)
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(10)+10) * time.Second)
	// 如果出错，就新建chatGPT对话
	if err != nil {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(15)+15) * time.Second)
		text, newToken, newURI = HandleErrorTokeninfoVersion(query, conversationIDPtr, parentMessagePtr, accessToken, baseURI, tokenInfo, cli, err)
	}
	return text, newToken, newURI
}

// HandleChatRobustlyTokeninfoVersion 若chatgpt访问出错，则新建一个chatgpt连接，以保证chatgpt健壮访问
func HandleChatRobustlyTokeninfoTestVersion(
	query string,
	conversationIDPtr *string,
	parentMessagePtr *string,
	accessToken, baseURI string,
	tokenInfo *TokenInfo,
	cli *chatgpt.Client,
	i int,
) (*chatgpt.ChatText, string, string) {

	var text *chatgpt.ChatText
	var (
		newToken string = accessToken
		newURI   string = baseURI
	)
	//如果panic，就新建chatGPT对话
	defer func() {
		if r := recover(); r != nil {
			rand.Seed(time.Now().UnixNano())
			time.Sleep(time.Duration(rand.Intn(60)+60) * time.Second)
			text, newToken, newURI = HandleErrorTokeninfoVersion(query, conversationIDPtr, parentMessagePtr, accessToken, baseURI, tokenInfo, cli, r)
		}
	}()
	// 正常访问chatgpt
	//text, err := cli.GetChatText(query, *conversationIDPtr, *parentMessagePtr)

	text, err := cli.GetChatText(query)
	if i%3 == 0 {

	} else if i%3 == 1 {
		err = errors.New("429")
	} else if i%3 == 2 {
		err = errors.New("501")
	}
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(10)+10) * time.Second)
	// 如果出错，就新建chatGPT对话
	if err != nil {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(15)+15) * time.Second)
		text, newToken, newURI = HandleErrorTokeninfoVersion(query, conversationIDPtr, parentMessagePtr, accessToken, baseURI, tokenInfo, cli, err)
	}

	if i == 9 {
		tokenInfo.ReleaseToken(newToken)
	}
	return text, newToken, newURI
}

// HandleErrorTokeninfoVersion tokenInfo版本
func HandleErrorTokeninfoVersion(
	query string, conversationIDPtr *string, parentMessagePtr *string, accessToken, baseURI string, tokenInfo *TokenInfo, cli *chatgpt.Client, err any,
) (*chatgpt.ChatText, string, string) {
	log.Println(accessToken, err)
	var newToken = accessToken
	var newURI = baseURI
	if err2, ok := err.(error); ok {
		//log.Println(accessToken, err2)
		if strings.Contains(err2.Error(), "429") || strings.Contains(err2.Error(), "202") {
			//cli = NewDefaultClient(uuid.New().String(), "https://freechat.lidong.xin")
			//indexOfToken, ok := tokenInfo.getIndexOfToken(accessToken)
			//if !ok {
			//	log.Fatalln(accessToken, "不在tokeninfo中！")
			//}
			time.Sleep(time.Second * 10)
			newToken, newURI = tokenInfo.Handle429(accessToken)
			cli = NewDefaultClient(newToken, newURI)
			//cli = NewDefaultClient(uuid.New().String(), "https://freechat.xyhelper.cn")	//境外服务器
		} else {
			//log.Println(accessToken,baseURI)
			//cli = NewDefaultClient(accessToken, baseURI)
			newToken, newURI = tokenInfo.Handle429(accessToken)
			cli = NewDefaultClient(newToken, newURI)
		}
	} else {
		//log.Println(accessToken, err)
		// 新建cli
		newToken, newURI = tokenInfo.Handle429(accessToken)
		cli = NewDefaultClient(newToken, newURI)
	}
	// 修改conversationID和parentMessage
	*conversationIDPtr = ""
	*parentMessagePtr = ""
	// 再次访问并返回结果
	//text, err := cli.GetChatText(query, *conversationIDPtr, *parentMessagePtr)
	text, err := cli.GetChatText(query)
	// 如果仍然出错，递归解决错误
	if err != nil {
		return HandleErrorTokeninfoVersion(query, conversationIDPtr, parentMessagePtr, accessToken, baseURI, tokenInfo, cli, err)
	}
	return text, newToken, newURI
}

// HandleChatRobustlyTokeninfoVersion 若chatgpt访问出错，则新建一个chatgpt连接，以保证chatgpt健壮访问
func HandleChatRobustlyTokeninfoVersionWithSession(
	query string,
	conversationIDPtr *string,
	parentMessagePtr *string,
	accessToken, baseURI string,
	tokenInfo *TokenInfo,
	cli *chatgpt.Client) (*chatgpt.ChatText, string, string) {

	var text *chatgpt.ChatText
	var (
		newToken string = accessToken
		newURI   string = baseURI
	)
	//如果panic，就新建chatGPT对话
	defer func() {
		if r := recover(); r != nil {
			rand.Seed(time.Now().UnixNano())
			time.Sleep(time.Duration(rand.Intn(60)+60) * time.Second)
			text, newToken, newURI = HandleErrorTokeninfoVersionWithSession(query, conversationIDPtr, parentMessagePtr, accessToken, baseURI, tokenInfo, cli, r)
		}
	}()
	// 正常访问chatgpt
	text, err := cli.GetChatText(query, *conversationIDPtr, *parentMessagePtr)
	//if text != nil {
	//
	//	*conversationIDPtr = text.ConversationID
	//	*parentMessagePtr = text.MessageID
	//}

	//text, err := cli.GetChatText(query)
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(10)+10) * time.Second)
	// 如果出错，就新建chatGPT对话
	if err != nil {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Duration(rand.Intn(15)+15) * time.Second)
		text, newToken, newURI = HandleErrorTokeninfoVersionWithSession(query, conversationIDPtr, parentMessagePtr, accessToken, baseURI, tokenInfo, cli, err)
	}
	return text, newToken, newURI
}

// HandleErrorTokeninfoVersion tokenInfo版本
func HandleErrorTokeninfoVersionWithSession(
	query string, conversationIDPtr *string, parentMessagePtr *string, accessToken, baseURI string, tokenInfo *TokenInfo, cli *chatgpt.Client, err any,
) (*chatgpt.ChatText, string, string) {
	log.Println(accessToken, err)
	var newToken = accessToken
	var newURI = baseURI
	if err2, ok := err.(error); ok {
		//log.Println(accessToken, err2)
		if strings.Contains(err2.Error(), "429") || strings.Contains(err2.Error(), "202") || strings.Contains(err2.Error(), "404") {
			//cli = NewDefaultClient(uuid.New().String(), "https://freechat.lidong.xin")
			//indexOfToken, ok := tokenInfo.getIndexOfToken(accessToken)
			//if !ok {
			//	log.Fatalln(accessToken, "不在tokeninfo中！")
			//}
			time.Sleep(time.Second * 10)
			newToken, newURI = tokenInfo.Handle429(accessToken)
			cli = NewDefaultClient(newToken, newURI)
			//cli = NewDefaultClient(uuid.New().String(), "https://freechat.xyhelper.cn")	//境外服务器
		} else {
			//log.Println(accessToken,baseURI)
			//cli = NewDefaultClient(accessToken, baseURI)
			newToken, newURI = tokenInfo.Handle429(accessToken)
			cli = NewDefaultClient(newToken, newURI)
		}
	} else {
		//log.Println(accessToken, err)
		// 新建cli
		//cli = NewDefaultClient(accessToken, baseURI)
		newToken, newURI = tokenInfo.Handle429(accessToken)
		cli = NewDefaultClient(newToken, newURI)
	}
	// 修改conversationID和parentMessage
	*conversationIDPtr = ""
	*parentMessagePtr = ""
	// 再次访问并返回结果
	//text, err := cli.GetChatText(query, *conversationIDPtr, *parentMessagePtr)
	text, err := cli.GetChatText(query)
	//if text != nil {
	//	*conversationIDPtr = text.ConversationID
	//	*parentMessagePtr = text.MessageID
	//}
	// 如果仍然出错，递归解决错误
	if err != nil {
		return HandleErrorTokeninfoVersionWithSession(query, conversationIDPtr, parentMessagePtr, accessToken, baseURI, tokenInfo, cli, err)
	}
	return text, newToken, newURI
}
