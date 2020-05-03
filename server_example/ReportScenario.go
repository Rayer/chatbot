package main

import (
	ChatBot "github.com/rayer/chatbot"
)

type ReportScenario struct {
	ChatBot.DefaultScenarioImpl
	ThisWeekInDev []string
	ThisWeekDone  []string
}

func (rs *ReportScenario) InitScenario(uc *ChatBot.UserContext) error {
	rs.DefaultScenarioImpl.InitScenario(uc)
	rs.RegisterState("entry", &ReportEntryState{}, rs)
	rs.RegisterState("creating_done", &ReportCreatingDone{}, rs)
	rs.RegisterState("creating_indev", &ReportCreatingInDev{}, rs)
	rs.RegisterState("confirm", &ReportConfirm{}, rs)
	return nil
}

func (rs *ReportScenario) EnterScenario(source ChatBot.Scenario) error {
	return nil
}

func (rs *ReportScenario) ExitScenario(askFrom ChatBot.Scenario) error {
	return nil
}

func (rs *ReportScenario) DisposeScenario() error {
	return nil
}

func (rs *ReportScenario) Name() string {
	return "Weekly Report Scenario"
}

/*
States :
1. Entry - Greeting with current period report or re-create, if not, [Create Report]
2. CreatingDone
3. CreatingInDev
4. Review
*/

type ReportEntryState struct {
	ChatBot.DefaultScenarioStateImpl
}

func (res *ReportEntryState) InitScenarioState(scenario ChatBot.Scenario) {
	res.Init(scenario, res)
	res.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "create report",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			err := res.ChangeStateByName("creating_done")
			return "Ok let's creating a report", err
		},
	})
	res.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "view reports",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			return "Not really implemented in this prototype version... maybe later", nil
		},
	})
	res.KeywordHandler.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "exit",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			err := res.ReturnLastScenario()
			return "Let's back to previous session", err

		},
	})
}

func (res *ReportEntryState) RawMessage() (string, error) {
	/*
		Designed functionality :
		1. Let user view logs before (Not in this prototype)
		2. Show if log is submitted in this week. If so, show it and ask if it need to be recreate or exit
		3. If no report in this week, ask user to create one
	*/

	return "Hey, we don't see logs this week. Would you like to [create report]? or [view reports] in previous weeks? You also can [exit] if no longer need to operating with logs", nil
}

type ReportCreatingDone struct {
	ChatBot.DefaultScenarioStateImpl
}

func (rcd *ReportCreatingDone) InitScenarioState(scenario ChatBot.Scenario) {
	rcd.Init(scenario, rcd)
	rcd.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "good for now",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			err := rcd.ChangeStateByName("creating_indev")
			return "Done in done", err
		},
	})

	rcd.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			doneList := rcd.GetParentScenario().(*ReportScenario).ThisWeekDone
			rcd.GetParentScenario().(*ReportScenario).ThisWeekDone = append(doneList, input)
			return "Recorded (done) : " + input, nil
		},
	})
}

func (rcd *ReportCreatingDone) RawMessage() (string, error) {
	return "What task have been done in this week? or there is [good for now]?", nil
}

type ReportCreatingInDev struct {
	ChatBot.DefaultScenarioStateImpl
}

func (rcid *ReportCreatingInDev) InitScenarioState(scenario ChatBot.Scenario) {
	rcid.Init(scenario, rcid)
	rcid.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "good for now",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			rcid.ChangeStateByName("confirm")
			return "Done in dev", nil
		},
	})
	rcid.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			indevList := rcid.GetParentScenario().(*ReportScenario).ThisWeekInDev
			rcid.GetParentScenario().(*ReportScenario).ThisWeekInDev = append(indevList, input)

			return "Recorded (indev): " + input, nil
		},
	})
}

func (rcid *ReportCreatingInDev) RawMessage() (string, error) {
	return "What task is in dev this week? or it's [good for now]?", nil
}

type ReportConfirm struct {
	ChatBot.DefaultScenarioStateImpl
}

func (rc *ReportConfirm) InitScenarioState(scenario ChatBot.Scenario) {
	rc.Init(scenario, rc)
	rc.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "submit",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			err := rc.ReturnLastScenario()
			return "Submitted", err
		},
	})
	rc.RegisterKeyword(&ChatBot.Keyword{
		Keyword: "discard",
		Action: func(keyword string, input string, scenario ChatBot.Scenario, state ChatBot.ScenarioState) (s string, e error) {
			err := rc.ReturnLastScenario()
			return "Discarded", err
		},
	})
}

func (rc *ReportConfirm) RawMessage() (string, error) {
	doneList := rc.GetParentScenario().(*ReportScenario).ThisWeekDone
	indevList := rc.GetParentScenario().(*ReportScenario).ThisWeekInDev

	ret := "Will you [submit] or [discard] follow report entries : "
	ret += "Done : \n"
	for _, done := range doneList {
		ret += done + "\n"
	}

	ret += "In Dev : \n"
	for _, inDev := range indevList {
		ret += inDev + "\n"
	}

	return ret, nil

}
