package main

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

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/chat", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", gin.H{})
	})
	r.GET("/ws", handleWS)
	if err := r.Run(":80"); err != nil {
		log.Fatal(err)
	}
}
func handleWS(c *gin.Context) {
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
		log.Print(session)
		if err := json.Unmarshal(message, session); err != nil {
			log.Println(err)
			continue
		}
		if len(session.History) == 0 {
			log.Println("History is empty")
			continue
		}
		question := session.History[0]
		res, err := getCompletion(session, question)
		log.Print(res)
		if err != nil {
			message := err.Error()
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
	completion, err := client.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	answer := completion.Choices[0].Text
	session.History = append(session.History, question)
	session.History = append(session.History, answer)
	return answer, nil
}

type ChatMessage struct {
	User    string `json:"user"`
	Message string `json:"message"`
}
