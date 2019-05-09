package handler_test

import (
	"strings"
	"testing"

	"gotest.tools/assert"

	handler "github.com/doms/spongemock"
)

func TestSpongeMock(t *testing.T) {
	testCases := []struct {
		testStr     string
		expectedStr string
	}{
		{
			testStr:     "test",
			expectedStr: "TeSt",
		}, {
			testStr:     "",
			expectedStr: "",
		}, {
			testStr:     "こんにちは",
			expectedStr: "こんにちは",
		}, {
			testStr:     "swaこんg",
			expectedStr: "SwAこんg",
		}, {
			testStr:     "hey @user, how's it going?",
			expectedStr: "HeY @user, HoW'S it GoinG?",
		}, {
			testStr:     "the party is happening in #party-room",
			expectedStr: "ThE PaRtY is HaPpEninG in #party-room",
		}, {
			testStr:     "<USERID|User>",
			expectedStr: "<USERID|User>",
		},
	}

	for _, test := range testCases {
		splitContent := strings.Split(test.testStr, " ")
		var res []string

		for _, word := range splitContent {
			res = append(res, handler.SpongeMock(word))
		}

		buf := strings.Join(res, " ")
		assert.Equal(t, buf, test.expectedStr)
	}
}
