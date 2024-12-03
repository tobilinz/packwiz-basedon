package basedon

import (
	"fmt"

	"github.com/packwiz/packwiz/cmd"
	"github.com/spf13/cobra"
)

var basedtest = &cobra.Command{
	Use:   "basedtest",
	Args:  cobra.NoArgs,
	Short: "Sets a base modpack for this modpack",
	Run: func(_cmd *cobra.Command, args []string) {
		fmt.Println("Note: This is a development command. Use this for testing purposes only.")
		//TODO:  This stuff needs to be recursive
		// resolver.Resolve(pack.BasedOn.Type, pack.BasedOn.Info)
		// merger.Merge()

		MergePacks()
	},
}

func init() {
	cmd.Add(basedtest)
}
