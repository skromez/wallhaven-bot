package main

import (
	"encoding/json"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type Image struct {
	ID        string
	Url       string
	Path      string
	Favorites int
	Resolution string
}

func main() {
	err := godotenv.Load()
	bot, err := tg.NewBotAPI(os.Getenv("API_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("bot is working")
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		fmt.Printf("user's message: %v \n", update.Message.Text)

		request, err := http.NewRequest("GET", "http://localhost:3000/find-images?category=anime&q="+update.Message.Text, nil)
		if err != nil {
			fmt.Println("failed")
			log.Fatal(err)
		}

		resp, err := tg.HttpClient.Do(&http.Client{}, request)
		if err != nil {
			fmt.Println("failed")
			log.Fatal(err)
		}
		var images []Image

		if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
			fmt.Println(resp.Status)
			fmt.Println("failed 1")
		}
		var imagesTest []interface{}
		fmt.Println(images)
		for _, image := range images {
			fmt.Println(image.Path)
			fmt.Println(image.Resolution)
			newImage := tg.NewInputMediaPhoto(image.Path)
			imagesTest = append(imagesTest, newImage)
		}
		if len(imagesTest) == 0 {
			msg := tg.NewMessage(update.Message.Chat.ID, "nothing found")
			_, err = bot.Send(msg)
		}
		msg := tg.NewMediaGroup(update.Message.Chat.ID, imagesTest)

		_, err = bot.Send(msg)
		if err != nil {
			fmt.Println("failed")
			msg := tg.NewMessage(update.Message.Chat.ID, "nothing found")
			_, err = bot.Send(msg)
		}
	}
}
