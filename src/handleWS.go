package src

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	gogpt "github.com/sashabaranov/go-gpt3"
	"io/ioutil"
	"log"
	"net/http"
)

type ChatSession struct {
	History       []string
	HistoryAnswer string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWS(c *gin.Context) {
	//读取key.config中的密钥
	key, err := ioutil.ReadFile("key.config")
	if err != nil {
		log.Println(err)
		return
	}
	//设置密钥
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()
	for {
		_, message, err := ws.ReadMessage()
		log.Print(string(message))
		if err != nil {
			log.Println(err)
			break
		}
		session := &ChatSession{}
		if err := json.Unmarshal(message, session); err != nil {
			log.Println(err)
			continue
		}
		question := ""
		if len(session.History) == 1 {
			log.Println("History is empty")
			question = "Prompt(Ask): " + session.History[0] + "(Direct Answer,think step by step)"
		} else {
			question = session.History[len(session.History)-1]
			question = "Prompt(Ask): " + question
			session.History = session.History[:len(session.History)-1]
			//为前几个session的问题添加字符串：Previous Questions：
			for i := 0; i < len(session.History)-1; i++ {
				session.History[i] = "Prompt(Previous Questions): " + session.History[i]
			}
			PreviousAnswer := "Prompt(Your Previous Answer): " + session.HistoryAnswer
			session.History = append(session.History, PreviousAnswer)
			//session.History[len(session.History)-1] = "Previous Questions: " + session.History[len(session.History)-1]
		}
		res, err := getCompletion(session, question, string(key))
		if err != nil {
			message := err.Error()
			log.Println("Err:", message)
			_ = ws.WriteJSON(ChatMessage{
				User:    "ai",
				Message: message,
			})
			continue
		}
		//sessionData, _ := json.Marshal(session)
		//encodedSessionData := base64.StdEncoding.EncodeToString(sessionData)
		//cookie := http.Cookie{
		//	Name:  "chat_session",
		//	Value: encodedSessionData,
		//	Path:  "/",
		//}
		//http.SetCookie(c.Writer, &cookie)
		_ = ws.WriteJSON(ChatMessage{
			User:    "ai",
			Message: res,
		})
	}
}
func getCompletion(session *ChatSession, question string, key string) (string, error) {
	client := gogpt.NewClient("")
	ctx := context.Background()
	prompt := ""
	if len(session.History) != 1 {
		for _, q := range session.History {
			prompt += q + "\n"
		}
	}
	prompt += question + "\n"
	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 1024,
		Prompt:    prompt,
	}
	log.Print("Prompt:" + prompt)
	completion, err := client.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	answer := completion.Choices[0].Text
	session.History = append(session.History, question)
	session.History = append(session.History, answer)
	log.Print("Answer:", answer)
	//如果answer的第一个字符是“\n”，则去掉
	if answer[0] == '\n' {
		answer = answer[1:]
	}
	if answer[:7] == "Answer:" {
		answer = answer[7:]
	}
	return answer, nil
}

type ChatMessage struct {
	User    string `json:"user"`
	Message string `json:"message"`
}
