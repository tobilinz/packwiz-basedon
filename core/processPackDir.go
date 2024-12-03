package core

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

var ignoreDefaults = []string{
	// Defaults (can be overridden with a negating pattern preceded with !)

	// Exclude Git metadata
	".git/**",
	".gitattributes",
	".gitignore",

	// Exclude macOS metadata
	".DS_Store",

	// Exclude exported CurseForge zip files
	"/*.zip",

	// Exclude exported Modrinth packs
	"*.mrpack",

	// Exclude packwiz binaries, if the user puts them in their pack folder
	"packwiz.exe",
	"packwiz", // Note: also excludes packwiz/ as a directory - you can negate this pattern if you want a directory called packwiz
}

var ignoreFileName = ".packwizignore"

func readIgnoreFile(path string) (*gitignore.GitIgnore, error) {
	lines := ignoreDefaults

	data, err := os.ReadFile(path)
	if err == nil {
		//TODO: try to not split
		lines = append(lines, strings.Split(string(data), "\n")...)
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	return gitignore.CompileIgnoreLines(lines...), nil
}

func ProcessPackDir(pack *Pack, fn func(path string, info os.DirEntry, relPath string) error) error {
	ignore, err := readIgnoreFile(filepath.Join(pack.GetRootPath(), ignoreFileName))
	if err != nil {
		return err
	}

	return filepath.WalkDir(pack.GetRootPath(), func(path string, info os.DirEntry, err error) error {
		if err != nil {
			// TODO: Handle errors on individual files properly
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(pack.GetRootPath(), path)
		if err != nil {
			return err
		}

		if path == pack.SavePath || relPath == pack.Index.File || info.Name() == ignoreFileName {
			return nil
		}

		if ignore.MatchesPath(relPath) {
			if info.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		return fn(path, info, relPath)
	})
}
