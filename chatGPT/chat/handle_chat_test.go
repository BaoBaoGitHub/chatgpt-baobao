package chat

import (
	"github.com/google/uuid"
	"github.com/xyhelper/chatgpt-go"
	"log"
	"testing"
	"time"
)

var query = `java code for "Closes the current scans"`

// new chatgpt client
var token = uuid.New().String()

var cli = chatgpt.NewClient(
	chatgpt.WithDebug(false),
	chatgpt.WithTimeout(60*time.Second),
	chatgpt.WithAccessToken(token),
	chatgpt.WithBaseURI("https://freechat.lidong.xin"),
)

var conversationID = "test"
var parentMessage = "test"

var text *chatgpt.ChatText
var baseURI = "https://freechat.lidong.xin"

func TestHandleError(t *testing.T) {

	//HandleError(query, &conversationID, &parentMessage, token, baseURI, cli)
}

func TestHandlePanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			//text = HandleError(query, &conversationID, &parentMessage, token, baseURI, cli)
		}
	}()
	panic("test")
}
