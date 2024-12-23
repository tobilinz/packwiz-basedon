package basedon

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/packwiz/packwiz/basedon/merger"
	"github.com/packwiz/packwiz/basedon/resolver"
	"github.com/packwiz/packwiz/core"
)

var baseDisabledSuffix = ".base.disabled"

func MergePacks() core.Pack {
	cachePath, err := core.GetProjectCache()
	if err != nil {
		fmt.Println("An error occured while trying to get the cache directory")
		os.Exit(1)
	}

	thisPack, err := core.LoadPack()
	if err != nil {
		fmt.Println("Error loading pack:", err)
		os.Exit(1)
	}

	fmt.Println("Resolving base pack...")
	resolver.Resolve(thisPack.BasedOn.Type, thisPack.BasedOn.Info)

	basePath := filepath.Join(cachePath, "base", filepath.Dir(thisPack.BasedOn.PackLocation))
	basePack, err := core.LoadAnyPack(filepath.Join(basePath, filepath.Base(thisPack.BasedOn.PackLocation)), false)
	if err != nil {
		fmt.Println("Error loading pack:", err)
		os.Exit(1)
	}

	fmt.Println("Merging this and base...")
	mergedPath := filepath.Join(cachePath, "merged")

	err = os.RemoveAll(mergedPath)
	if err != nil {
		fmt.Println("Failed to unzip the base modpack ", err)
		os.Exit(1)
	}

	err = os.MkdirAll(mergedPath, 0744)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mergedPack := core.Pack{
		Name:        thisPack.Name,
		Author:      thisPack.Author,
		Version:     thisPack.Version,
		Description: thisPack.Description + "\nBased on " + basePack.Name,
		PackFormat:  thisPack.PackFormat,
		BasedOn:     core.BasedOn{},
		Index: struct {
			File       string `toml:"file"`
			HashFormat string `toml:"hash-format"`
			Hash       string `toml:"hash,omitempty"`
		}{
			File: thisPack.Index.File,
		},
		Versions: thisPack.Versions,
		Export:   thisPack.Export,
		Options:  thisPack.Options,
		SavePath: filepath.Join(mergedPath, "pack.toml"),
	}

	//TODO: Add .packwizignore contents
	//TODO: -> This can somehow be done with some gitignore api that is also used for index refresh
	merger.Merge(mergedPath, &basePack, &thisPack)

	err = core.ProcessPackDir(&mergedPack, func(path string, info os.DirEntry, relPath string) error {
		if !strings.HasSuffix(relPath, baseDisabledSuffix) {
			return nil
		}

		err := os.Remove(path)
		if err != nil {
			return err
		}

		err = os.Remove(path[:len(path)-len(baseDisabledSuffix)])
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Printf("Error removing disabled duplicate files: %s\n", err)
		os.Exit(1)
	}

	mergedIndexPath := filepath.Join(mergedPath, mergedPack.Index.File)
	_, err = os.Stat(mergedIndexPath)
	if os.IsNotExist(err) {
		// Create file
		err = os.WriteFile(mergedIndexPath, []byte{}, 0644)
		if err != nil {
			fmt.Printf("Error creating index file: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(mergedIndexPath + " created!")
	} else if err != nil {
		fmt.Printf("Error checking index file: %s\n", err)
		os.Exit(1)
	}

	mergedIndex, err := mergedPack.LoadIndex()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = mergedIndex.Refresh()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = mergedIndex.Write()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = mergedPack.UpdateIndexHash()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = mergedPack.Write()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return mergedPack
}
