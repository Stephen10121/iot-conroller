package function

import "fmt"

func RunCommand(command string, params Command) {
	fmt.Println("Running command: " + command)
}

func CustomCommandChecker(str string, commands map[string]Command) bool {
	_, ok := commands[str]

	return ok
}
