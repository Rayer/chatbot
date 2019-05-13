package ChatBot

type ScenarioState interface {
	InitScenarioState(scenario Scenario)
	RenderMessage() (string, error)
	HandleMessage(input string) (string, error)
	GetParentScenario() Scenario
}

type DefaultScenarioStateImpl struct {
	parent Scenario
	KeywordHandler
	Utilities
}

func (dssi *DefaultScenarioStateImpl) Init(scenario Scenario, state ScenarioState) {
	dssi.parent = scenario
	dssi.KeywordHandler.Init(scenario, state)
	dssi.Utilities.Init(scenario, state)
}

func (dssi *DefaultScenarioStateImpl) GetParentScenario() Scenario {
	return dssi.parent
}

