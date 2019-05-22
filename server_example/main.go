package main

import (
	"encoding/json"
	"fmt"
	ChatBot "github.com/Rayer/chatbot"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func WebKeywordHandler(fullText string, keyword string, validKeyword bool) string {
	if validKeyword {
		return fmt.Sprintf("%s", keyword)
	}

	return fmt.Sprintf("<font color = 'red'>%s</font>", keyword)
}

type MessageDetail struct {
	Response string `json:"response"`
	ValidKeywordList []string `json:"validKeywords"`
	InvalidKeywordList []string `json:"invalidKeywords"`
}

func main() {

	conf := ChatBot.Configuration{ResetTimerSec:300, KeywordFormatter: WebKeywordHandler}
	ctm := ChatBot.NewContextManagerWithConfig(&conf)

	http.HandleFunc("/chatbot", func(writer http.ResponseWriter, request *http.Request) {

		if request.Method != http.MethodPost {
			writer.WriteHeader(404)
			writer.Write([]byte("Invalid method"))
			return
		}

		request.ParseForm()
		name := request.PostForm["name"][0]
		phrase := request.PostForm["phrase"][0]

		ctx := ctm.CreateUserContext(name, func() ChatBot.Scenario {
			return &RootScenario{}
		})

		dbg, _ := ctx.HandleMessage(phrase)

		origin := request.Header.Get("Origin")
		log.Info("Origin : %s", origin)

		writer.Header().Set("Access-Control-Allow-Origin", origin)
		writer.WriteHeader(200)

		output, validKeywords, invalidKeywords, _ := ctx.RenderMessageWithDetail()

		response := MessageDetail{
			output,
			validKeywords,
			invalidKeywords,
		}

		ret, err := json.Marshal(response)

		if err != nil {
			log.Fatal(err)
		}

		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(200)
		writer.Write([]byte(ret))
		log.Printf("Name : %s\nPhrase : %s\nRes : %s\nRet : %s", name, phrase, dbg, ret)

		//writer.Write([]byte(ret))
	})

	http.HandleFunc("/chatbot/detail", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		writer.Write([]byte("OK!"))
	})

	http.ListenAndServe(":12160", nil)
}