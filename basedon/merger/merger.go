package merger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	tml "github.com/BurntSushi/toml"

	"github.com/packwiz/packwiz/core"
)

type Merger struct {
	Regex    *regexp.Regexp
	Encoding string
	Function func(base, this []byte, relPath string) ([]byte, error)
}

var baseMergers = []*Merger{json, options, properties, toml}
var mergers = append([]*Merger{fabricLoaderDependencies}, baseMergers...)

func Merge(mergedProjectPath string, basePack *core.Pack, thisPack *core.Pack) {
	attributions := make(map[string]string)

	err := copyDirectories(mergedProjectPath, basePack, func(path string, info os.DirEntry, relPath string) error {
		base, err := os.Open(path)
		if err != nil {
			return err
		}
		defer base.Close()

		mergedPath := filepath.Join(mergedProjectPath, relPath)
		merged, err := os.Create(mergedPath)
		if err != nil {
			return err
		}
		defer merged.Close()

		_, err = io.Copy(merged, base)
		if err != nil {
			return err
		}

		source, err := getSource(path, basePack.Name)
		if err != nil {
			return err
		}
		attributions[relPath] = source

		return nil
	})
	if err != nil {
		println("An error occured while attemting to copy the base project files into the merge directory.", err.Error())
		os.Exit(1)
	}

	err = copyDirectories(mergedProjectPath, thisPack, func(path string, info os.DirEntry, relPath string) error {
		thisContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		mergedPath := filepath.Join(mergedProjectPath, relPath)
		mergedContent, err := os.ReadFile(mergedPath)
		if os.IsNotExist(err) {
			err = os.WriteFile(mergedPath, thisContent, 0644)
			if err != nil {
				return err
			}
			source, err := getSource(path, thisPack.Name)
			if err != nil {
				return err
			}
			attributions[relPath] = source

			return nil
		}
		if err != nil {
			return err
		}

		var merger func(base, this []byte, relPath string) ([]byte, error) = nil
		for _, v := range mergers {
			if v.Regex.MatchString(relPath) {
				merger = v.Function
				break
			}
		}

		if merger == nil {
			fmt.Println("Failed to merge ", filepath.Base(path), " and ", filepath.Base(mergedPath), " because no merger for those file types could be found. Using the file declared in this modpack as fallback.")
			return nil
		}

		//TODO: rework error system. Add warnings too, because mergers often just want to warn and fall back to something
		mergedContent, err = merger(mergedContent, thisContent, relPath)
		if err != nil {
			return err
		}

		err = os.WriteFile(mergedPath, mergedContent, 0644)
		if err != nil {
			return err
		}

		source, err := getSource(path, thisPack.Name+" and "+basePack.Name)
		if err != nil {
			return err
		}
		attributions[relPath] = source

		return nil
	})
	if err != nil {
		println("An error occured while attemting to merge this projects files into the merge directory.", err.Error())
		os.Exit(1)
	}

	//TODO: Base modpack (or this modpack) might have changed source manually in source.md. Respect that
	attributionsList := "| Filepath Relative to This Directory | Source URL or Author Name |\n|--------|--------|\n"

	for path, source := range attributions {
		attributionsList += fmt.Sprintf("| %s | %s |\n", path, source)
	}

	attributionsList += fmt.Sprintf("\n*%s by %s is based on %s by %s.*", thisPack.Name, thisPack.Author, basePack.Name, basePack.Author)

	file, err := os.Create(filepath.Join(mergedProjectPath, "ATTRIBUTIONS.md"))
	if err != nil {
		println("failed to create file: %w", err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = file.WriteString(attributionsList)
	if err != nil {
		println("failed to write to file: %w", err)
		os.Exit(1)
	}
}

func getSource(path string, fallbackSource string) (string, error) {
	source := fallbackSource
	if strings.HasSuffix(path, ".pw.toml") {
		var data map[string]interface{}
		if _, err := tml.DecodeFile(path, &data); err != nil {
			return "", err
		}

		downloadSection, ok := data["download"].(map[string]interface{})
		if !ok {
			return "", errors.New("'download' section not found or invalid")
		}

		url, ok := downloadSection["url"].(string)
		if !ok {
			return "", errors.New("'url' not found or is not a string")
		}

		source = url
	}

	return source, nil
}

func copyDirectories(mergedProjectPath string, pack *core.Pack, fn func(path string, info os.DirEntry, relPath string) error) error {
	return core.ProcessPackDir(pack, func(path string, info os.DirEntry, relPath string) error {
		err := os.MkdirAll(filepath.Join(mergedProjectPath, filepath.Dir(relPath)), 0744)
		if err != nil {
			return err
		}

		return fn(path, info, relPath)
	})
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
