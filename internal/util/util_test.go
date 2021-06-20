package util_test

import (
	"testing"

	"github.com/automatico/jato/internal/util"
)

func TestUnderscorer(t *testing.T) {
	t.Parallel()
	type testCase struct {
		have string
		want string
	}
	testCases := []testCase{
		{have: "show version", want: "show_version"},
		{have: " show version ", want: "show_version"},
		{have: "Show run-config commands", want: "show_run_config_commands"},
	}

	for _, tc := range testCases {
		got := util.Underscorer(tc.have)
		if tc.want != got {
			t.Errorf("want %s, got %s", tc.want, got)
		}
	}

}

func TestTruncateOutput(t *testing.T) {
	t.Parallel()
	type testCase struct {
		have string
		want string
	}
	testCases := []testCase{
		{have: "a\r\nb\r\nc", want: "b"},
		{have: "", want: ""},
		{have: "a", want: "a"},
	}

	for _, tc := range testCases {
		got := util.TruncateOutput(tc.have)
		if tc.want != got {
			t.Errorf("want %s, got %s", tc.want, got)
		}
	}

}
