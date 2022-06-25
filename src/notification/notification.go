package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
)

type MessageObj struct {
	UserID  string `json:"userId"`
	Message string `json:"message"`
}

/* メッセージ送信 */
func sendMessage(bot *linebot.Client, p *events.SQSMessage) error {
	// 取得したMessageをデコードする。。
	fmt.Println("*** message decode")

	var message MessageObj

	err := json.Unmarshal([]byte(p.Body), &message)
	if err != nil {
		return fmt.Errorf("*** message decode error.[%s]", err)
	}

	fmt.Println("*** push")
	if _, err := bot.PushMessage(message.UserID, linebot.NewTextMessage(message.Message)).Do(); err != nil {
		log.Fatal(err)
		return fmt.Errorf("massages push error.[%s]", message.UserID)
	}

	return nil
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {

	fmt.Println("*** linebot new")
	// LINEのBotの設定
	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		return fmt.Errorf("linebot new error.[%s]", err)
	}

	// メッセージをSQSから取得
	for _, message := range sqsEvent.Records {
		// メッセージ送信
		err := sendMessage(bot, &message)
		if err != nil {
			return fmt.Errorf("sendMessage error.[%s]", err)
		}
	}

	// 終了
	fmt.Println("*** end")
	return nil
}

func main() {
	lambda.Start(handler)
}
