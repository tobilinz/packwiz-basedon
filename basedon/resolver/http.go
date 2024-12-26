package resolver

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/packwiz/packwiz/core"
	"github.com/spf13/cobra"
)

var httpCommand = &cobra.Command{
	Use:   "http",
	Args:  cobra.NoArgs,
	Short: "Sets a base modpack for this modpack as http source",
}

var zipCommand = &cobra.Command{
	Use:   "zip [url]",
	Args:  cobra.ExactArgs(1),
	Short: "Sets a base modpack for this modpack as http zip source",
	Run: func(cmd *cobra.Command, args []string) {
		basePackURI := args[0]

		pack, err := core.LoadPack()
		if err != nil {
			fmt.Println("An error occurred while trying to load pack:", err.Error())
			os.Exit(1)
		}

		pack.BasedOn.Type = "http"
		pack.BasedOn.Info = map[string]string{
			"filetype": "zip",
			"url":      basePackURI,
		}

		err = pack.Write()
		if err != nil {
			fmt.Println("An error occurred while trying to write pack:", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	httpCommand.AddCommand(zipCommand)
}

func Http(info map[string]string) error {
	switch info["filetype"] {
	case "zip":
		return getZip(info["url"])
	case "directory":
		//TODO:
		return nil
	default:
		return errors.New("Filetype " + info["filetype"] + " is not supported by https.")
	}
}

func getZip(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	cachePath, err := core.GetProjectCache()
	if err != nil {
		return err
	}

	zipLocation := path.Join(cachePath, "base.zip")

	err = os.Remove(zipLocation)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	file, err := os.Create(zipLocation)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	reader, err := zip.NewReader(file, fileInfo.Size())
	if err != nil {
		return err
	}

	basePath := path.Join(cachePath, "base")
	err = os.RemoveAll(basePath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(basePath, 0744)
	if err != nil {
		return err
	}

	var baseName string
	for index, file := range reader.File {
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		if index == 0 {
			baseName = file.Name
			continue
		}

		rel, err := filepath.Rel(baseName, file.Name)
		if err != nil {
			return err
		}

		outputPath := filepath.Join(basePath, rel)

		if file.FileInfo().IsDir() {
			err = os.MkdirAll(outputPath, 0744)
			if err != nil {
				return err
			}
		} else {
			outputFile, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			defer outputFile.Close()

			_, err = io.Copy(outputFile, fileReader)
			if err != nil {
				return err
			}
		}
	}

	file.Close()
	os.Remove(file.Name())

	return nil
}
