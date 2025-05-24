package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

var nikitosReplicas []string
var conversationHistory = make(map[int64][]map[string]string)

const maxHistory = 20

func getLocalOllamaResponse(chatID int64, prompt string) (string, error) {
	history := conversationHistory[chatID]
	history = append(history, map[string]string{"role": "user", "content": prompt})
	if len(history) > maxHistory {
		history = history[len(history)-maxHistory:]
	}
	conversationHistory[chatID] = history

	reqBody := map[string]interface{}{
		"model":    "nikitos",
		"messages": append([]map[string]string{{"role": "system", "content": "Ты — Никитос.  "}}, history...),
		"stream":   true,
	}

	data, _ := json.Marshal(reqBody)
	log.Println("➡️ Вызов Ollama через POST /api/chat")

	resp, err := http.Post("http://ollama:11434/api/chat", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var full strings.Builder
	decoder := json.NewDecoder(resp.Body)

	for decoder.More() {
		var chunk struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
			Done bool `json:"done"`
		}

		if err := decoder.Decode(&chunk); err != nil {
			log.Println("❌ Ошибка декодирования:", err)
			break
		}
		full.WriteString(chunk.Message.Content)
		if chunk.Done {
			break
		}
	}

	reply := full.String()
	conversationHistory[chatID] = append(conversationHistory[chatID], map[string]string{"role": "assistant", "content": reply})
	return reply, nil
}

func loadReplicas(filename string) []string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Println("не удалось загрузить реплики:", err)
		return []string{}
	}
	var list []string
	err = json.Unmarshal(data, &list)
	if err != nil {
		log.Println("ошибка парсинга JSON:", err)
		return []string{}
	}
	return list
}

func main() {
	log.Println("🔥 Nikitos бот стартует")
	rand.Seed(time.Now().UnixNano())
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не найден")
	}

	nikitosReplicas = loadReplicas("nikitos_replicas.json")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)

	var lastMessageTime = time.Now()

	for update := range updates {
		log.Printf("UPDATE: %+v\n", update)
		if update.Message == nil || update.Message.Text == "" {
			continue
		}

		text := update.Message.Text
		chat := update.Message.Chat
		user := update.Message.From

		shouldRespond := false
		if chat.IsPrivate() || strings.Contains(text, "@"+bot.Self.UserName) {
			shouldRespond = true
		} else if chat.IsGroup() || chat.IsSuperGroup() {
			if rand.Intn(10) == 0 {
				shouldRespond = true
			}
		}

		if !shouldRespond {
			if time.Since(lastMessageTime) > 15*time.Second && rand.Intn(5) == 0 {
				mention := fmt.Sprintf("[@%s](tg://user?id=%d)", user.FirstName, user.ID)
				randomComment := mention + ", ты это слышал?"
				msg := tgbotapi.NewMessage(chat.ID, randomComment)
				msg.ParseMode = "Markdown"
				bot.Send(msg)
				lastMessageTime = time.Now()
			}
			continue
		}

		userInput := strings.ReplaceAll(text, "@"+bot.Self.UserName, "")
		userInput = strings.TrimSpace(userInput)

		if len(nikitosReplicas) > 0 && rand.Intn(1000) < 10 {
			reply := nikitosReplicas[rand.Intn(len(nikitosReplicas))]
			log.Println("🟡 Редкий ответ из шаблона:", reply)
			msg := tgbotapi.NewMessage(chat.ID, reply)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			continue
		}

		reply, err := getLocalOllamaResponse(chat.ID, userInput)
		log.Printf("📥 Ответ от Ollama: %q", reply)
		if err != nil || strings.TrimSpace(reply) == "" {
			if len(nikitosReplicas) > 0 {
				reply = nikitosReplicas[rand.Intn(len(nikitosReplicas))]
			} else {
				reply = "Ну и зачем ты это сказал?"
			}
		}

		msg := tgbotapi.NewMessage(chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = "Markdown"
		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("Ошибка при отправке сообщения: %v", err)
		} else {
			log.Println("Сообщение отправлено успешно")
		}

		lastMessageTime = time.Now()
		time.Sleep(time.Duration(3+rand.Intn(5)) * time.Second)
	}
}
