package main

import (
	"ChatBot"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)
import "net/http"

type IncomingMessage struct {
	User string `json:"user"`
	Text string `json:"text"`
}

func main() {
	ctm := ChatBot.NewContextManager()
	ctm.Configuration.ResetTimerSec = 300

	http.HandleFunc("/react", func(w http.ResponseWriter, r *http.Request) {
		incomingMessage := IncomingMessage{}
		err := r.ParseForm()
		log.Info(r.Form)

		if err != nil {
			log.Error("Error parse form : %v", err)
			return
		}

		for key := range r.Form {
			log.Println(key)
			//LOG: {"test": "that"}
			err := json.Unmarshal([]byte(key), &incomingMessage)
			if err != nil {
				log.Println(err.Error())
			}
		}



		ctx := ctm.GetUserContext(incomingMessage.User)

		if ctx == nil {
			ctx = ctm.CreateUserContext(incomingMessage.User, &WelcomeScenario{})
		}

		response, err := ctx.RenderMessage()
		if err != nil {
			log.Error("Error rendering message()")
		}




	})
}