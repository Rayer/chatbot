package ChatBot

type ScenarioState interface {
	InitScenarioState(scenario Scenario)
	RenderMessage() (string, error)
	RawMessage() (string, error)
	HandleMessage(input string) (string, error)
	GetParentScenario() Scenario
}

type DefaultScenarioStateImpl struct {
	parent Scenario
	KeywordHandler
	Utilities
	host ScenarioState
}

func (dssi *DefaultScenarioStateImpl) Init(scenario Scenario, state ScenarioState) {
	dssi.parent = scenario
	dssi.KeywordHandler.Init(scenario, state)
	dssi.Utilities.Init(scenario, state)
	dssi.host = state //workaround.....
}

func (dssi *DefaultScenarioStateImpl) GetParentScenario() Scenario {
	return dssi.parent
}

func (dssi *DefaultScenarioStateImpl) RenderMessage() (string, error) {
	message, err := dssi.host.RawMessage()
	if err != nil {
		return message, err
	}
	return dssi.TransformRawMessage(message)
}


