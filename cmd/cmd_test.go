package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

type TestCommander struct {
	command_output_file string
}

func (c TestCommander) Execute(command string, args ...string) ([]byte, error) {
	output, err := os.ReadFile(c.command_output_file)
	if err != nil {
		return nil, fmt.Errorf("problem with reading fixture file, err: %v", err)
	}
	return output, nil
}

func TestCmdExec(t *testing.T) {
	testcases := []struct {
		name                string
		args                []string
		command_output_file string
		diff_file           string
		assert              bool
	}{
		{
			name:                "two_pods",
			args:                []string{"compare", "pods", "kuttle", "kuttle2"},
			command_output_file: "../test/fixtures/two_pods_output.yaml",
			diff_file:           "../test/fixtures/two_pods_diff.yaml",
			assert:              true,
		},
	}

	for _, tc := range testcases {
		var kubectl TestCommander
		kubectl.command_output_file = tc.command_output_file
		rootCmd.SetArgs(tc.args)
		err := rootCmd.ParseFlags(tc.args)
		if err != nil {
			t.Errorf("cannot parse flags")
		}
		cmdout := new(bytes.Buffer)
		compareCmd := addCompareCmd(kubectl)
		compareCmd.SetOut(cmdout)
		err = compareCmd.RunE(compareCmd, tc.args)
		if err != nil {
			t.Errorf("compare error: %v", err)
		}
		resStdout, err := io.ReadAll(cmdout)
		if err != nil {
			t.Errorf("error reading command output: %v", err)
		}
		diff_output, err := os.ReadFile(tc.diff_file)
		if err != nil {
			t.Errorf("problem with reading fixture file, err: %v", err)
		}
		if !bytes.Equal(resStdout, diff_output) {
			t.Error("wrong diff")
		}
	}
}
