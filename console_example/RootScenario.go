package main

import (
	ChatBot "github.com/Rayer/chatbot"
	log "github.com/sirupsen/logrus"
	"strings"
)

type RootScenario struct {
	ChatBot.DefaultScenarioImpl
}

func (rs *RootScenario) InitScenario(uc *ChatBot.UserContext) error {
	rs.DefaultScenarioImpl.InitScenario(uc)
	rs.RegisterState("entry", &EntryState{}, rs)
	rs.RegisterState("second", &SecondState{}, rs)
	return nil
}

func (rs *RootScenario) EnterScenario(source ChatBot.Scenario) error {
	log.Debugln("Entering root scenario")
	return nil
}

func (rs *RootScenario) ExitScenario(askFrom ChatBot.Scenario) error {
	log.Debugln("Exiting root scenario")
	return nil
}

func (rs *RootScenario) DisposeScenario() error {
	log.Debugln("Disposing root scenario")
	return nil
}

//It's Scenario State
//The only state of the root scenario
type EntryState struct {
	ChatBot.DefaultScenarioStateImpl
}

func (es *EntryState) InitScenarioState(scenario ChatBot.Scenario) {
	es.Init(scenario, es)
	es.RegisterKeyword(&ChatBot.Keyword{Keyword:"submit report", Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (string, error) {
		err := es.InvokeNextScenario(&ReportScenario{}, ChatBot.Stack)
		return "Go to report scenario", err
	}})

	es.RegisterKeyword(&ChatBot.Keyword{Keyword:"manage broadcasts", Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
		err := es.ChangeStateByName("second")
		return "Exit with 2", err
	}})
}

func (es *EntryState) RenderMessage() (string, error) {
	rawMsg := "Hey it's BossBot! Are you going to [submit report], [manage broadcasts] or [check]?"
	return es.TransformRawMessage(rawMsg)
}

func (es *EntryState) HandleMessage(input string) (string, error) {
	return es.KeywordHandler.ParseAction(input)
}

type SecondState struct {
	ChatBot.DefaultScenarioStateImpl
}

func (ss *SecondState) InitScenarioState(scenario ChatBot.Scenario) {
	ss.Init(scenario, ss)
}

func (ss *SecondState) RenderMessage() (string, error) {
	raw := "This is second message, you can only [exit] in order to get out of here"
	return ss.KeywordHandler.TransformRawMessage(raw)
}

func (ss *SecondState) HandleMessage(input string) (string, error) {
	if strings.Contains(input, "exit") {
		err := ss.ChangeStateByName("entry")
		return "Exiting...", err
	}
	return "Not exit, stay here.", nil
}

func (rs *RootScenario) Name() string {
	return "RootScenario"
}
