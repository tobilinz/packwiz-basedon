package resolver

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Resolve(sourceType string, source map[string]string) {
	switch sourceType {
	case "http":
		err := Http(source)
		if err != nil {
			fmt.Println("An error occured while trying to resolve the base pack.", err.Error())
		}
	default:
		fmt.Println("Cannot parse a base modpack with source of type ", sourceType)
		os.Exit(1)
		break
	}
}

func GetResolverCommands() []*cobra.Command {
	return []*cobra.Command{
		httpCommand,
	}
}
