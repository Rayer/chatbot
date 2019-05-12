package ChatBot

type ScenarioState interface {
	InitScenarioState(scenario Scenario)
	RenderMessage() (string, error)
	HandleMessage(input string) (string, error)
	GetParentScenario() Scenario
	SetParentScenario(parent Scenario)
}

type KeywordAction func(keyword string, input string, scenario Scenario, state ScenarioState) (string, error)

type Keyword struct {
	Keyword string
	Action  KeywordAction
}

type KeywordHandler struct {
	keywordList []Keyword
	scenario    Scenario
	state       ScenarioState
}

type DefaultScenarioStateImpl struct {
	parent Scenario
}

func (dssi *DefaultScenarioStateImpl) GetParentScenario() Scenario {
	return dssi.parent
}

func (dssi *DefaultScenarioStateImpl) SetParentScenario(parent Scenario) {
	dssi.parent = parent
}
