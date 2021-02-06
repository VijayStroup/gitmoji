package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type command struct {
	name  string
	emoji string
}

var (
	commands map[string]string
)

func init() {
	commands := make(map[string]string)

	// initialize default commands
	commands["new"] = "✨"
	commands["fix"] = "🔧"
	commands["update"] = "☝️"

	// get commands from environment variables
	getEnvs(commands)

	// loop through commands and add them to rootCmd
	for k, v := range commands {
		rootCmd.AddCommand(makeCommand(k, v))
	}
}

// makeCommand generates a pointer to a cobra.Command to add to the rootCmd.
// This is done dynamically in the init() function.
func makeCommand(name, emoji string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Prepend %s to git commit message", emoji),
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
}

// getEnvs adds to a map of environment variables set by the user specific to
// Gitmoji. They are then processed to ensure they meet the correct format of
// "command:emoji". Example: "fix:🔧"
func getEnvs(c map[string]string) {
	for _, env := range os.Environ() {
		// check if env is not a Gitmoji env
		if !strings.HasPrefix(env, "GITM_") {
			continue
		}

		// get env name and value
		envSplit := strings.Split(env, "=")
		name := envSplit[0]
		value := strings.ReplaceAll(envSplit[1], " ", "")

		// check to see if value format is good
		v := strings.Split(value, ":")
		if len(v) != 2 {
			fmt.Printf("\033[31mEnvironment variable '%s' is is of wrong format: '%s'.\n", name, value)
			fmt.Printf("The correct format is 'command:emoji'.\nExample: 'fix:🔧'.\033[0m\n\n")
			os.Exit(1)
		}

		// check to see if name is lowercase
		if v[0] != strings.ToLower(v[0]) {
			fmt.Printf("\033[33mWarning: environment variable %s does not have key of type lowercase.\n", name)
			fmt.Printf("'%s' will be treated as '%s'. Change %s to supress this warning.\033[0m\n\n", v[0], strings.ToLower(v[0]), name)
		}
		v[0] = strings.ToLower(v[0])

		// all checks passed, add to commands map
		c[v[0]] = v[1]
	}
}
