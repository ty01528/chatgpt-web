package src

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	gogpt "github.com/sashabaranov/go-gpt3"
	"log"
	"net/http"
)

type ChatSession struct {
	History []string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWS(c *gin.Context) {
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
		log.Print("session:", session)
		if err := json.Unmarshal(message, session); err != nil {
			log.Println(err)
			continue
		}
		question := ""
		if len(session.History) == 0 {
			log.Println("History is empty")
			question = session.History[0]
			continue
		}
		log.Print("session.History:", session.History)
		question = session.History[len(session.History)-1]
		res, err := getCompletion(session, question)
		if err != nil {
			message := err.Error()
			log.Println("Err:", message)
			_ = ws.WriteJSON(ChatMessage{
				User:    "ai",
				Message: message,
			})
			continue
		}
		sessionData, _ := json.Marshal(session)
		encodedSessionData := base64.StdEncoding.EncodeToString(sessionData)
		cookie := http.Cookie{
			Name:  "chat_session",
			Value: encodedSessionData,
			Path:  "/",
		}
		http.SetCookie(c.Writer, &cookie)
		_ = ws.WriteJSON(ChatMessage{
			User:    "ai",
			Message: res,
		})
	}
}
func getCompletion(session *ChatSession, question string) (string, error) {
	client := gogpt.NewClient("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")
	ctx := context.Background()
	prompt := ""
	for _, q := range session.History {
		prompt += q + "\n"
	}
	prompt += question + "\n"
	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 1024,
		Prompt:    prompt,
	}
	log.Print("Print:", prompt)
	completion, err := client.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	answer := completion.Choices[0].Text
	session.History = append(session.History, question)
	session.History = append(session.History, answer)
	log.Print("Answer:", answer)
	return answer, nil
}

type ChatMessage struct {
	User    string `json:"user"`
	Message string `json:"message"`
}
