package ChatBot

type Utilities struct {
	scenario Scenario
	state ScenarioState
	initialized bool
}

func (u *Utilities) Init(scenario Scenario, state ScenarioState) {
	u.scenario = scenario
	u.state = state
	u.initialized = true
}

func (u* Utilities) checkInitialized() {
	if !u.initialized {
		panic("Utilities is used but not yet initialized!")
	}
}

func (u *Utilities) InvokeNextScenario(scenario Scenario, strategy InvokeStrategy) error {
	u.checkInitialized()
	return u.state.GetParentScenario().GetUserContext().InvokeNextScenario(scenario, strategy)
}

func (u *Utilities) ChangeStateByName(stateName string) error {
	u.checkInitialized()
	return u.state.GetParentScenario().ChangeStateByName(stateName)
}

func (u *Utilities) ReturnLastScenario() error {
	u.checkInitialized()
	return u.state.GetParentScenario().GetUserContext().ReturnLastScenario()
}