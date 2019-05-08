package handler_test

import (
	"testing"

	"gotest.tools/assert"

	"github.com/doms/spongemock"
)

func TestSpongeMock(t *testing.T) {
	testCases := []struct{
		testStr string
		expectedStr string
	}{{
		testStr: "test",
		expectedStr: "TeSt",
	}, {
		testStr: "",
		expectedStr: "",
	}, {
		testStr: "こんにちは",
		expectedStr: "こんにちは",
	}, {
		testStr: "swaこんg",
		expectedStr: "SwAこんg",
	}}

	for _, test := range testCases {
		buf := handler.SpongeMock(test.testStr)
		assert.Equal(t, buf, test.expectedStr)
	}
}
