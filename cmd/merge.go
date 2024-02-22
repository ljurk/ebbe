package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ljurk/ebbe/helper"
	"github.com/spf13/cobra"
)

var (
	mergeCmd = &cobra.Command{
		Use:   "merge",
		Short: "Merge multiple command together, later arguments will be subtracted from the previous ones",
		Run: func(cmd *cobra.Command, args []string) {
			var messages []string
			var tmp []string
			var err error
			input, _ := cmd.Flags().GetStringSlice("input")
			commands, _ := cmd.Flags().GetBool("commands")

			for _, i := range input {
				if commands {
					tmp, err = runCommand(i)
					if err != nil {
						fmt.Println("Error running ", err)
						return
					}
				} else {
					tmp, err = helper.ReadFromFile(i)
					if err != nil {
						fmt.Printf("%v\n", err)
						return
					}
				}
				messages = helper.MergeLists(tmp, messages)
			}
			for _, i := range messages {
				fmt.Print(i)
			}
		},
	}
)

func runCommand(commandString string) ([]string, error) {
	fields := strings.Fields(commandString)
	if len(fields) == 0 {
		return nil, errors.New("empty command string")
	}

	// Extract the command and its arguments
	command := fields[0]
	args := fields[1:]

	// Run the command
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Split the output into lines
	lines := strings.Split(string(output), "\n")
	var result []string
	for _, line := range lines {
		if line != "" {
			result = append(result, line+"\n")
		}
	}
	return result, nil
}

func init() {
	mergeCmd.Flags().StringSliceP("input", "i", []string{}, "Specify multiple files")
	mergeCmd.Flags().BoolP("commands", "c", false, "Pass commands to input instead of files")
	err := mergeCmd.MarkFlagRequired("input")
	if err != nil {
		fmt.Println(err)
	}
	rootCmd.AddCommand(mergeCmd)
}
