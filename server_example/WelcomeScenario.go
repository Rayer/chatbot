package main

import (
	"fmt"
	ChatBot "github.com/Rayer/chatbot"
)

type WelcomeScenario struct{
	ChatBot.DefaultScenarioImpl
}

func (ws *WelcomeScenario) InitScenario(uc *ChatBot.UserContext) error {
	ws.DefaultScenarioImpl.InitScenario(uc)
	ws.RegisterState("entry", &EntryState{}, ws)
	return nil
}

func (ws *WelcomeScenario) EnterScenario(source ChatBot.Scenario) error {
	return nil
}

func (ws *WelcomeScenario) ExitScenario(askFrom ChatBot.Scenario) error {
	return nil
}

func (ws *WelcomeScenario) DisposeScenario() error {
	return nil
}

func (ws *WelcomeScenario) Name() string {
	return "WelcomeScenario"
}

type EntryState struct {
	ChatBot.DefaultScenarioStateImpl
}

func (es *EntryState) InitScenarioState(scenario ChatBot.Scenario) {

}

func (es *EntryState) RenderMessage() (string, error) {
	name := es.GetParentScenario().GetUserContext().User
	return fmt.Sprintf("Hello %s!", name), nil

}

func (es *EntryState) HandleMessage(input string) (string, error) {
	return "You can review [system] statistics, [application] statistics or view [link lists].", nil
}
