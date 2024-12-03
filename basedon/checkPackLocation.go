package basedon

import (
	"fmt"
	"os"
	"path"
)

func CheckPackLocation(baseDir string) string {
	fmt.Printf("Checking, of a pack.toml file is present at %s...\n", packLocation)
	packPath := path.Join(baseDir, packLocation)

	_, err := os.Stat(packPath)
	switch {
	case err == nil:
		break
	case os.IsNotExist(err):
		fmt.Println("Pack.toml file does not exist at ", packPath)
		os.Exit(1)
	default:
		fmt.Println("An error occurred while checking if pack.toml file exists at ", packPath)
		os.Exit(1)
	}

	fmt.Println("Pack.toml file is present")
	return packPath
}
