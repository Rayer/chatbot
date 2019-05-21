package main

import (
	"bufio"
	"fmt"
	ChatBot "github.com/Rayer/chatbot"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)


func KeywordTransformer(fulltext string, keyword string) string {
	return fmt.Sprintf("<spin color=\"red\">%s</spin>", keyword)
}


func main() {
	logrus.SetLevel(logrus.WarnLevel)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter ID [BotSpec] : ")
	id, err := reader.ReadString('\n')
	//Trim \n
	id = strings.Replace(id, "\n", "", -1)

	if err != nil {
		panic(err.Error())
	}
	if id == "" {
		id = "BotSpec"
	}
	fmt.Println("Welcome " + id + ", start invoking session...")
	conf := ChatBot.Configuration{
		ResetTimerSec: 300,
		KeywordFormatter: KeywordTransformer,
	}
	ctx := ChatBot.NewContextManagerWithConfig(&conf)
	utx := ctx.CreateUserContext(id, func() ChatBot.Scenario {
		return &RootScenario{}
	})

	fmt.Println(utx.RenderMessage())
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if text == "exitloop" {
			break
		}
		fmt.Println(utx.HandleMessage(text))
		fmt.Println(utx.RenderMessage())
	}
}
