package handler_test

import (
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
			testStr:     "  test  ",
			expectedStr: "  TeSt  ",
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
		}, {
			testStr:     "20",
			expectedStr: "20",
		}, {
			testStr:     "20_20",
			expectedStr: "20_20",
		}, {
			testStr:     "2020",
			expectedStr: "2020",
		}, {
			testStr:     "2",
			expectedStr: "2",
		},
	}

	for _, test := range testCases {
		assert.Equal(t, handler.SpongeMock(test.testStr), test.expectedStr)
	}
}
