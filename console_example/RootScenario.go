package main

import (
	ChatBot "github.com/Rayer/chatbot"
	"github.com/fatih/color"
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

func keywordTransformer(fullMessage string, keyword string) string {
	yellow := color.New(color.BgRed).SprintFunc()
	return yellow(keyword)
}

//It's Scenario State
//The only state of the root scenario
type EntryState struct {
	ChatBot.DefaultScenarioStateImpl
	ChatBot.KeywordHandler
}

func (es *EntryState) InitScenarioState(scenario ChatBot.Scenario) {
	es.KeywordHandler.Init(scenario, es)
	es.RegisterKeyword(&ChatBot.Keyword{Keyword:"submit report", Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (string, error) {
		es.GetParentScenario().GetUserContext().InvokeNextScenario(&ReportScenario{}, ChatBot.Stack)
		return "Go to report scenario", nil
	}})

	es.RegisterKeyword(&ChatBot.Keyword{Keyword:"manage broadcasts", Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
		es.GetParentScenario().ChangeStateByName("second")
		return "Exit with 2", nil
	}})

	es.KeywordHandler.OnEachKeyword = keywordTransformer
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
	ChatBot.KeywordHandler
}

func (ss *SecondState) InitScenarioState(scenario ChatBot.Scenario) {
	ss.KeywordHandler.Init(scenario, ss)
}

func (ss *SecondState) RenderMessage() (string, error) {
	raw := "This is second message, you can only [exit] in order to get out of here"
	return ss.KeywordHandler.TransformRawMessage(raw)
}

func (ss *SecondState) HandleMessage(input string) (string, error) {
	if strings.Contains(input, "exit") {
		ss.GetParentScenario().ChangeStateByName("entry")
		return "Exiting...", nil
	}
	return "Not exit, stay here.", nil
}

func (rs *RootScenario) Name() string {
	return "RootScenario"
}
