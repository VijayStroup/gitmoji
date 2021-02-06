package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "fix",
	Short: "Prepend EMOJI to git commit message.",
	Run: func(cmd *cobra.Command, args []string) {
		commit := exec.Command("git", "status")
		out, err := commit.Output()
		if err != nil {
			// if the exitcode is 128, that is an indication of the current
			// directory not being a git repo, so let's tell the user
			if exitError, _ := err.(*exec.ExitError); exitError.ExitCode() == 128 {
				fmt.Println("gitm: 🚨 No git repo found in the current directory.")
			} else {
				// if another error occurs that is not checked for, alert the user
				// of the error
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println(string(out))
		}
		// if err := commit.Run(); err != nil {
		// 	fmt.Println("Cant run that command")
		// }
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
