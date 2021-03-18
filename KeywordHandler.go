package ChatBot

import (
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type KeywordAction func(keyword string, input string, scenario Scenario, state ScenarioState) (string, error)
type KeywordFormatter func(fullMessage string, keyword string, isValidKeyword bool) string

type Keyword struct {
	Keyword string
	Action  KeywordAction
}

type KeywordHandler struct {
	keywordList      []Keyword
	scenario         Scenario
	state            ScenarioState
	initialized      bool
	KeywordFormatter KeywordFormatter
}

func NewKeywordHandler(scenario Scenario, state ScenarioState) *KeywordHandler {
	return &KeywordHandler{scenario: scenario, state: state, KeywordFormatter: GetConfiguration().KeywordFormatter, initialized:true}
}

func (kh *KeywordHandler) Init(scenario Scenario, state ScenarioState) {
	kh.initialized = true
	kh.scenario = scenario
	kh.state = state
	kh.KeywordFormatter = GetConfiguration().KeywordFormatter
}

func (kh *KeywordHandler) checkInitialized() {
	if !kh.initialized {
		panic("KeywordHandler is not initialized yet")
	}
}

func (kh *KeywordHandler) RegisterKeyword(keyword *Keyword) {
	kh.checkInitialized()
	if kh.keywordList == nil {
		kh.keywordList = []Keyword{}
	}
	kh.keywordList = append(kh.keywordList, *keyword)
}

func (kh *KeywordHandler) ParseAction(input string) (string, error) {
	kh.checkInitialized()
	for _, kw := range kh.keywordList {
		//bypass default one
		if kw.Keyword == "" {
			continue
		}
		if strings.Contains(strings.ToLower(input), strings.ToLower(kw.Keyword)) {
			ret, err := kw.Action(kw.Keyword, input, kh.scenario, kh.state)
			if err != nil {
				return "", errors.Wrap(err, "Error parsing action : "+kw.Keyword)
			}
			return ret, nil
		}
	}

	//if we have default keyword
	for _, kw := range kh.keywordList {
		if kw.Keyword == "" {
			ret, err := kw.Action(kw.Keyword, input, kh.scenario, kh.state)
			if err != nil {
				return "", errors.Wrap(err, "Error parsing action : "+kw.Keyword)
			}
			return ret, nil
		}
	}

	return "No match keyword", nil
}

func (kh *KeywordHandler) TransformRawMessage(rawMessage string) (transformedMessage string, validKeywords []string, invalidKeywords []string) {

	kh.checkInitialized()
	transformedMessage = rawMessage
	r, _ := regexp.Compile(`\[([A-Za-z 0-9_]*)]`)
	keywords := r.FindAllString(rawMessage, -1)

	for _, k := range keywords {
		k = strings.Replace(k, "[", "", -1)
		k = strings.Replace(k, "]", "", -1)
		invalidKeywords = append(invalidKeywords, k)
	}


	for _, keywordDefine := range kh.keywordList {
		//TODO: Maybe we should use map to avoid O(n^2)?
		for _, keyword := range keywords {
			//Skip default keyword
			if keyword == "" {
				continue
			}

			originalKeyword := keyword
			keyword = strings.Replace(keyword, "[", "", -1)
			keyword = strings.Replace(keyword, "]", "", -1)

			//TODO: Do we need case sensitive?
			if strings.ToLower(keywordDefine.Keyword) == strings.ToLower(keyword) {
				transformedKeyword := kh.KeywordFormatter(rawMessage, keyword, true)
				transformedMessage = strings.Replace(transformedMessage, originalKeyword, transformedKeyword, -1)
				validKeywords = append(validKeywords, keyword)

				//Remove keyword in invalidKeywords
				for i, v := range invalidKeywords {
					if v == keyword {
						invalidKeywords = append(invalidKeywords[:i], invalidKeywords[i+1:]...)
						break
					}
				}
				break
			}
		}
	}

	//Transform invalid keywords
	for _, ik := range invalidKeywords {
		transformedMessage = strings.Replace(transformedMessage, "[" + ik + "]", kh.KeywordFormatter(rawMessage, ik, false), -1)
	}

	return transformedMessage, validKeywords, invalidKeywords

}
