package ChatBot

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)


type KeywordhandlerTest struct {
	suite.Suite
	defaultKeywordHandler *KeywordHandler
	defaultAction         func(keyword string, input string, scenario Scenario, state ScenarioState) (s string, e error)
	defaultKeywords 	  []string
}

func TestKeywordHandler(t *testing.T) {
	suite.Run(t, new(KeywordhandlerTest))
}

func (k *KeywordhandlerTest)SetupTest() {
	k.defaultAction = func(keyword string, input string, scenario Scenario, state ScenarioState) (s string, e error) {
		return fmt.Sprintf("Get keyword %s", keyword), nil
	}

	k.defaultKeywords = []string{"KV1", "KV2", "KV3"}
	// TODO: We should add back mocked scenario and state
	k.defaultKeywordHandler = &KeywordHandler{
		keywordList:      make([]Keyword, 0, 0),
		scenario:         nil,
		state:            nil,
		initialized:      true,
		KeywordFormatter: DefaultKeywordFormatter,
	}
	k.defaultKeywordHandler.RegisterKeyword(&Keyword{
		Keyword: k.defaultKeywords[0],
		Action: k.defaultAction,
	})
	k.defaultKeywordHandler.RegisterKeyword(&Keyword{
		Keyword: k.defaultKeywords[1],
		Action: k.defaultAction,
	})
	k.defaultKeywordHandler.RegisterKeyword(&Keyword{
		Keyword: k.defaultKeywords[2],
		Action: k.defaultAction,
	})
}

func (k *KeywordhandlerTest)TearDownTest() {

}

func (k *KeywordhandlerTest)TestParseActionWithoutDefault() {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Only contains keyword",
			args:    args{
				input: k.defaultKeywords[0],
			},
			want:    fmt.Sprintf("Get keyword %s", k.defaultKeywords[0]),
			wantErr: false,
		},
		{
			name:    "Sentence contains one keyword",
			args:    args{
				input: "This is a phase that contains " + k.defaultKeywords[1] + " inside!",
			},
			want:    fmt.Sprintf("Get keyword %s", k.defaultKeywords[1]),
			wantErr: false,
		},
		{
			name:    "Sentence contains no keyword",
			args:    args{
				input: "ABCDEFG, There is no keyword!",
			},
			want:    "No match keyword",
			wantErr: false,
		},
		{
			name:    "Sentence contains multiple keywords",
			args:    args{
				input: fmt.Sprintf("Combined with multiple keyword, %s and %s", k.defaultKeywords[1], k.defaultKeywords[2]),
			},
			want:    fmt.Sprintf("Get keyword %s", k.defaultKeywords[1]),
			wantErr: false,
		},
	}
	kh := k.defaultKeywordHandler
	for _, tt := range tests {
		k.Run(tt.name, func() {
			got, err := kh.ParseAction(tt.args.input)
			if (err != nil) != tt.wantErr {
				k.T().Errorf("ParseAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				k.T().Errorf("ParseAction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (k *KeywordhandlerTest)TestParseActionWithDefaultKeyword() {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Only contains keyword",
			args:    args{
				input: k.defaultKeywords[0],
			},
			want:    fmt.Sprintf("Get keyword %s", k.defaultKeywords[0]),
			wantErr: false,
		},
		{
			name:    "Sentence contains one keyword",
			args:    args{
				input: "This is a phase that contains " + k.defaultKeywords[1] + " inside!",
			},
			want:    fmt.Sprintf("Get keyword %s", k.defaultKeywords[1]),
			wantErr: false,
		},
		{
			//TODO: It might be not correct behavior, need change and review
			name:    "Sentence contains no keyword",
			args:    args{
				input: "ABCDEFG, There is no keyword!",
			},
			want:    "Get keyword ",
			wantErr: false,
		},
		{
			name:    "Sentence contains multiple keywords",
			args:    args{
				input: fmt.Sprintf("Combined with multiple keyword, %s and %s", k.defaultKeywords[1], k.defaultKeywords[2]),
			},
			want:    fmt.Sprintf("Get keyword %s", k.defaultKeywords[1]),
			wantErr: false,
		},
	}
	kh := k.defaultKeywordHandler
	kh.RegisterKeyword(&Keyword{
		Keyword: "",
		Action:  k.defaultAction,
	})
	for _, tt := range tests {
		k.T().Run(tt.name, func(t *testing.T) {

			got, err := kh.ParseAction(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseAction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func (k *KeywordhandlerTest)TestTransformRawMessage() {

	type args struct {
		rawMessage string
	}
	tests := []struct {
		name                   string
		args                   args
		formatter			   KeywordFormatter
		wantTransformedMessage string
		wantValidKeywords      []string
		wantInvalidKeywords    []string
	}{
		{
			name:                   "DefaultSetup",
			args:                   args{
				rawMessage: fmt.Sprintf("Here is a message including [%s], [%s] and [%s]", k.defaultKeywords[0], k.defaultKeywords[1], k.defaultKeywords[2]),
			},
			formatter: DefaultKeywordFormatter,
			wantTransformedMessage: fmt.Sprintf("Here is a message including [%s], [%s] and [%s]", k.defaultKeywords[0], k.defaultKeywords[1], k.defaultKeywords[2]),
			wantValidKeywords:      k.defaultKeywords[0:3],
			wantInvalidKeywords:    []string{},
		},
		{
			name: "CustomizedFormatter",
			args: args{
				rawMessage: fmt.Sprintf("Here is a message including [%s], [%s] and [%s], [%s] is invalid.", k.defaultKeywords[0], k.defaultKeywords[1], k.defaultKeywords[2], "invalid one"),
			},
			formatter: func(fullMessage string, keyword string, isValidKeyword bool) string {
				if isValidKeyword {
					return fmt.Sprintf("{ %s }", keyword)
				}
				return fmt.Sprintf("<<< %s >>>", keyword)
			},
			wantTransformedMessage: fmt.Sprintf("Here is a message including { %s }, { %s } and { %s }, <<< %s >>> is invalid.", k.defaultKeywords[0], k.defaultKeywords[1], k.defaultKeywords[2], "invalid one"),
			wantValidKeywords:      k.defaultKeywords[0:3],
			wantInvalidKeywords:    []string{"invalid one"},
		},
		{
			name: "InvalidKeywords",
			args: args{
				rawMessage: fmt.Sprintf("Here is a message including [%s], [%s] and [%s] and [%s] should be invalid.", k.defaultKeywords[0], k.defaultKeywords[1], k.defaultKeywords[2], "invalid one"),
			},
			formatter: DefaultKeywordFormatter,
			wantTransformedMessage: fmt.Sprintf("Here is a message including [%s], [%s] and [%s] and <<%s>> should be invalid.", k.defaultKeywords[0], k.defaultKeywords[1], k.defaultKeywords[2], "invalid one"),
			wantValidKeywords:      k.defaultKeywords[0:3],
			wantInvalidKeywords:    []string{"invalid one"},
		},
	}
	for _, tt := range tests {
		k.T().Run(tt.name, func(t *testing.T) {
			kh := k.defaultKeywordHandler
			k.defaultKeywordHandler.KeywordFormatter = tt.formatter
			gotTransformedMessage, gotValidKeywords, gotInvalidKeywords := kh.TransformRawMessage(tt.args.rawMessage)
			if gotTransformedMessage != tt.wantTransformedMessage {
				t.Errorf("TransformRawMessage() gotTransformedMessage = %v, want %v", gotTransformedMessage, tt.wantTransformedMessage)
			}
			if !reflect.DeepEqual(gotValidKeywords, tt.wantValidKeywords) {
				t.Errorf("TransformRawMessage() gotValidKeywords = %v, want %v", gotValidKeywords, tt.wantValidKeywords)
			}
			if !reflect.DeepEqual(gotInvalidKeywords, tt.wantInvalidKeywords) {
				t.Errorf("TransformRawMessage() gotInvalidKeywords = %v, want %v", gotInvalidKeywords, tt.wantInvalidKeywords)
			}
		})
	}
}