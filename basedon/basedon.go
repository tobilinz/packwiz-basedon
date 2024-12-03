package basedon

import (
	"fmt"
	"os"

	"github.com/packwiz/packwiz/basedon/resolver"
	"github.com/packwiz/packwiz/cmd"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/cobra"
)

var basedonCmd = &cobra.Command{
	Use:   "basedon",
	Args:  cobra.NoArgs,
	Short: "Sets a base modpack for this modpack",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		pack, err := core.LoadPack()
		if err != nil {
			fmt.Println("An error occurred while trying to load pack:", err)
			os.Exit(1)
		}

		pack.BasedOn.PackLocation = packLocation

		err = pack.Write()
		if err != nil {
			fmt.Println("An error occurred while trying to write pack:", err)
			os.Exit(1)
		}
	},
}

var packLocation string

func init() {
	cmd.Add(basedonCmd)
	basedonCmd.PersistentFlags().StringVarP(&packLocation, "pack-location", "p", cmd.GetRootFlag("pack-file").DefValue, "Sets the location of the pack.toml file relative to the pack root of the base modpack")
	basedonCmd.AddCommand(resolver.GetResolverCommands()...)
}

func AddBasedOn(newcommand *cobra.Command) {
	basedonCmd.AddCommand(newcommand)
}
