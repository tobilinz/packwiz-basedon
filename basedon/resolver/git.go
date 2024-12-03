package resolver

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

func Git(uri string, destination string) {
	_, err := git.PlainClone(destination, false, &git.CloneOptions{
		URL: uri,
	})

	switch {
	case err == nil:
		break
	case err.Error() == "repository not found":
		fmt.Println("The specified URL does not lead to a git repository.")
		os.Exit(1)
	default:
		fmt.Println("Error cloning repository:", err)
		os.Exit(1)
	}
}

/*
func getGitPack(basedOnOptions BasedOn) {
	cacheDir := cmd.GetCacheDir()

	tmpDirPath := path.Join(cacheDir, "basedon-tmp/")
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Println("Error removing basedon-tmp dir:", err)
			os.Exit(1)
		}
	}(tmpDirPath)

	dirStat, err := os.Stat(cacheDir)
	if err != nil {
		fmt.Println("Error checking basedon-tmp dir:", err)
		os.Exit(1)
	}

	err = os.MkdirAll(tmpDirPath, dirStat.Mode())
	if err != nil {
		fmt.Println("Error creating basedon-tmp dir:", err)
		os.Exit(1)
	}

	fmt.Println("\nDownloading base git repository...")
	repository, err := git.PlainClone(tmpDirPath, false, &git.CloneOptions{
		URL:          basedOnOptions.URL,
		Progress:     os.Stdout,
		Tags:         git.AllTags,
		SingleBranch: true,
	})
	if err != nil {
		fmt.Printf("An error occured while trying to clone the git repository at %s. Error: %s\n", basedOnOptions.URL, err)
		os.Exit(1)
	}

	gitOptions := &basedOnOptions.Git

	if gitOptions.Tag != "" {
		fmt.Println("\nChecking out correct tag...")

		worktree, err := repository.Worktree()
		if err != nil {
			fmt.Println("Error getting worktree:", err)
			os.Exit(1)
		}

		tag, err := repository.Tag(gitOptions.Tag)
		if err != nil {
			fmt.Println("Error getting tag:", err)
			os.Exit(1)
		}

		err = worktree.Checkout(&git.CheckoutOptions{
			Hash: tag.Hash(),
		})
		if err != nil {
			fmt.Printf("Error checking out tag %s: %s\n", gitOptions.Tag, err)
			os.Exit(1)
		}
	}

	packPath := path.Join(cacheDir, "./pack/")

	srcPath := path.Join(tmpDirPath, gitOptions.Path)
	err = os.Rename(srcPath, packPath)
	if err != nil {
		fmt.Printf("Error renaming base directory %s: %s\n", srcPath, err)
		os.Exit(1)
	}
}
*/
