package cmd

import (
	"fmt"
	"os"
	"os/exec"

	neat "github.com/itaysk/kubectl-neat/cmd"
	"github.com/spf13/cobra"
	yd "github.com/sters/yaml-diff/yamldiff"
	"gopkg.in/yaml.v2"
)

type YamlItems struct {
	Items []map[string]interface{} `json:"items"`
}

type Commander interface {
	Execute(string, ...string) ([]byte, error)
}

type RealCommander struct{}

func (c RealCommander) Execute(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	ret, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("execute returns error: %v", err)
	}
	return ret, nil
}

var (
	kubectl RealCommander

	rootCmd = &cobra.Command{
		Use:     "kubectl-compare",
		Short:   "kubectl-compare - a tool to compare two objects",
		Example: "kubectl-compare compare pods --namespace test pod1 pod2",
	}

	Version = "v0.0.0+unknown"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print kubectl-compare version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("kubectl-compare version: %s\n", Version)
		},
	}
)

func addCompareCmd(kubectl Commander) *cobra.Command {

	compareCmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare two objects",
		Args:  cobra.MatchAll(),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := append([]string{"get", "-o", "yaml"}, args...)
			output, err := kubectl.Execute("kubectl", cmdArgs...)
			if err != nil {
				return fmt.Errorf("kubectl returns error: %v", err)
			}
			yamlBytes, err := getYamlForCompare(output)
			if err != nil {
				return err
			}
			yamlNeatVersions := make([][]byte, 2)
			for i := 0; i < 2; i++ {
				yamlNeatVersions[i], err = neat.NeatYAMLOrJSON([]byte(yamlBytes[i]), "yaml")
				if err != nil {
					return fmt.Errorf("cannot neat yaml output for compare")
				}
			}
			yamlRawLists := make([]yd.RawYamlList, 2)
			for i := 0; i < 2; i++ {
				yamlRawLists[i], err = yd.Load(string(yamlNeatVersions[i]))
				if err != nil {
					return fmt.Errorf("cannot prepare yaml output for compare")
				}
			}
			compareResult := yd.Do(yamlRawLists[0], yamlRawLists[1])
			for _, c := range compareResult {
				cmd.Print(c.Dump())
			}

			return nil
		},
	}

	return compareCmd

}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(addCompareCmd(kubectl))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "there was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func getYamlForCompare(data []byte) ([][]byte, error) {
	var yamlItems YamlItems
	yamlItemsString := make([][]byte, 2)
	err := yaml.Unmarshal(data, &yamlItems)
	if err != nil {
		return nil, fmt.Errorf("yaml unmarshal error: %v", err)
	}
	if len(yamlItems.Items) != 2 {
		return nil, fmt.Errorf("two objects can be only compared")
	}
	for i := 0; i < 2; i++ {
		yamlItemsString[i], _ = yaml.Marshal(yamlItems.Items[i])
	}
	return yamlItemsString, nil
}
