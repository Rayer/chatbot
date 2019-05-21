package main

import (
	"fmt"
	ChatBot "github.com/Rayer/chatbot"
	"log"
	"net/http"
)

func WebKeywordHandler(fullText string, keyword string) string {
	return fmt.Sprintf("<font color=\"red\">%s</font>", keyword)
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


		writer.WriteHeader(200)

		ret, _ := ctx.RenderMessage()

		writer.Write([]byte(ret))
		log.Printf("Name : %s\nPhrase : %s\nRes : %s\nRet : %s", name, phrase, dbg, ret)

		//writer.Write([]byte(ret))
	})

	http.HandleFunc("/chatbot/lastresponse", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		writer.Write([]byte("OK!"))
	})

	http.ListenAndServe(":12160", nil)
}