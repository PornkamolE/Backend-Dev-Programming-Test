package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

func main() {
	// Set Value
	const channelSecret = "Your_Secret_Channel"
	const channelToken = "Your_Access_Token"

	if channelSecret == "" || channelToken == "" {
		log.Fatalf("ERROR: Missing environment variables:\n LINE_CHANNEL_SECRET: %v\n LINE_CHANNEL_TOKEN: %v",
			channelSecret, channelToken)
	}

	var err error
	bot, err = linebot.New(channelSecret, channelToken)
	if err != nil {
		log.Fatalf("ERROR: Failed to initialize LINE Bot: %v", err)
	}

	r := gin.Default()

	//เพิ่ม route สำหรับ GET /
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "LINE Bot is running!"})
	})

	r.POST("/webhook", handleWebhook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Starting server on port:", port)
	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("ERROR: Failed to start server: %v", err)
	}
}

func handleWebhook(c *gin.Context) {
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			log.Println("ERROR: Invalid signature")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		} else {
			log.Println("ERROR: Failed to parse request:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				handleTextMessage(event, message)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func handleTextMessage(event *linebot.Event, message *linebot.TextMessage) {
	replyToken := event.ReplyToken
	var err error
	responseText := ""

	switch message.Text {
	case "text":
		responseText = "สวัสดี! นี่คือ ChatBot นี่คือตัวอย่างการตอบกลับ"

	case "button":
		button := createButton()
		_, err = bot.ReplyMessage(replyToken, linebot.NewTemplateMessage("ปุ่ม", button)).Do()

	case "quickreply":
		quickReply := createQuickReply()
		message := linebot.NewTextMessage("กรุณาเลือกด้วยค่ะ").WithQuickReplies(quickReply)
		_, err = bot.ReplyMessage(replyToken, message).Do()

	case "carousel":
		carousel := createCarousel()
		_, err = bot.ReplyMessage(replyToken, linebot.NewTemplateMessage("Carousel", carousel)).Do()

	default:
		responseText = fmt.Sprintf(" %s ", message.Text)
	}

	if responseText != "" {
		_, err = bot.ReplyMessage(replyToken, linebot.NewTextMessage(responseText)).Do()
	}

	if err != nil {
		log.Println("ERROR: Failed to send message:", err)
	}
}
func createButton() *linebot.ButtonsTemplate {
	return linebot.NewButtonsTemplate(
		"", "🔘 หัวข้อปุ่ม", "🔹 นี่คือตัวอย่างปุ่ม",
		linebot.NewMessageAction("✨ กดปุ่มนี้", "🎉 คุณกดปุ่ม!"),
	)
}
func createQuickReply() *linebot.QuickReplyItems {
	return linebot.NewQuickReplyItems(
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("🎯 เลือก A", "A")),
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("🔥 เลือก B", "B")),
	)
}
func createCarousel() *linebot.CarouselTemplate {
	return linebot.NewCarouselTemplate(
		linebot.NewCarouselColumn(
			"", "ดูดวงประจำวัน", "เลือกไพ่แห่งโชคชะตาเพื่อดูคำทำนายของคุณ",
			linebot.NewMessageAction("เลือก ดูดวงประจำวัน", "คุณเลือกดูดวงประจำวัน"),
		),
		linebot.NewCarouselColumn(
			"", "อ่านข่าวสาร IT", "ข่าวไอทีที่น่าสนใจ คลิกเลย!",
			linebot.NewMessageAction("เลือก อ่านข่าวสาร IT", "คุณเลือกอ่านข่าวสาร IT"),
		),
		linebot.NewCarouselColumn(
			"", "Gossip News!", "ข่าวซุบซิบ วงในดารา ดราม่าร้อนฉ่า",
			linebot.NewMessageAction("เลือก Gossip News!", "คุณเลือกGossip News!"),
		),
	)
}
