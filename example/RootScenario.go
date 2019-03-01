package example

import (
	"ChatBot"
	"github.com/pkg/errors"
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

func (es *EntryState) RenderMessage() (string, error) {
	return "Hey it's BossBot! Are you going to [submit report], [manage broadcasts] or [check]?", nil
}

func (es *EntryState) HandleMessage(input string) (string, error) {
	if strings.Contains(input, "submit report") {
		err := es.GetParentScenario().GetUserContext().InvokeNextScenario(&ReportScenario{}, ChatBot.Stack)
		if err != nil {
			return "Error invoking report.... call Rayer to take a look?", errors.Wrap(err, "Error while invoking Report Scenario")
		}
		return "Ok let's make report then", nil
	} else if strings.Contains(input, "manage broadcast") {
		es.GetParentScenario().ChangeStateByName("second")
		return "Oops", nil
	} else if strings.Contains(input, "check") {
		es.GetParentScenario().ChangeStateByName("second")
		return "Oops", nil
	}

	return "I am not really understand....", nil
}

type SecondState struct {
	ChatBot.DefaultScenarioStateImpl
}

func (ss *SecondState) RenderMessage() (string, error) {
	return "Here is under construction, you can only [exit] in order to get out of here", nil
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
