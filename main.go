package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"os"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func main() {
	BOT_API := os.Getenv("BOT_API")
	
	bot, err := tgbotapi.NewBotAPI(BOT_API)
	if err != nil {
		log.Panic(err)
	}
	
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, MessageResuelt(update.Message.Text))
			msg.ReplyToMessageID = update.Message.MessageID 
			bot.Send(msg)
		}
	}
	http.ListenAndServe(":5000", nil)
}

func ConvertTimezoneName(city string)string{

	city_map := make(map[string]int)

	city_map["teh"] = 1
	city_map["ber"] = 2
	city_map["Ott"] = 3

	switch city_map[strings.ToLower(city)]{
	case 1:
		return "Tehran"
	case 2:
		return "Berlin"
	case 3:
		return "Ottawa"
	
	default:
		return "This bot dos not support you timezone"
	}

}

func ReqTimeZone(baselog string ,target string , time_ string ) string{
	TIME_API_KEY := os.Getenv("TIME_API_KEY")
	requestURL := "https://timezone.abstractapi.com/v1/convert_time/?api_key="+TIME_API_KEY+"&base_location="+baselog+"&base_datetime="+time.Now().Format("2006-01-02")+"T"+time_+"&target_location="+target
	resp, err := http.Get(requestURL)

	if err != nil {
        log.Fatalln(err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalln(err)
    }
    var data map[string]interface{}
    
    if err := json.Unmarshal([]byte(body), &data) ; err != nil {
        panic(err)
    }
	timese := data["target_location"].(map[string]any)
	timeser := timese["datetime"].(string)
    return timeser[11:16] + " " + target +"\n"
}

func ConvertTimeZone(baseloc string,time string) string {
	cityes := [3]string{"Tehran","Berlin","Ottawa"}
	var res string
	for _ , i := range cityes  {
		if i != baseloc {
			res = res + " " + ReqTimeZone(baseloc,i,time)
		}
	}
	return res
}

func MessageDecoder(msg string) (string,string){
	city := msg[:3]
	time := msg[4:]
	return city , time
}

func MessageResuelt(pm string)string{
	baseloc_city,baseloc_time := MessageDecoder(pm)
	msg := ConvertTimeZone(ConvertTimezoneName(baseloc_city),baseloc_time)
	return msg

}