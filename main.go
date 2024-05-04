package main

import (
	"context"
	"fmt"
	"os"

	"github.com/JGugino/fallout-3-cloud-saves/cmd"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v61/github"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("You must enter a mode.")
		fmt.Println("Usage: fcs {mode} - Modes: (init, upload, sync)")
		os.Exit(0)
	}

	mode := os.Args[1]

	//load config
	config, err := cmd.LoadConfig()

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// create github API client
	client := github.NewClient(nil).WithAuthToken(config.GithubApiKey)

	// ### MODE - INIT ###
	if mode == "init" {
		savesRepo := &github.Repository{
			Name:    github.String(config.RepoName),
			Private: github.Bool(true),
		}

		repo, _, err := client.Repositories.Get(context.Background(), config.CommiterName, config.RepoName)

		if err != nil {
			repo, _, _ = client.Repositories.Create(context.Background(), "", savesRepo)
		}

		_, err = git.PlainClone(cmd.TMP_REPO_PATH, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: config.CommiterName,
				Password: config.GithubApiKey,
			},
			URL:      repo.GetCloneURL(),
			Progress: os.Stdout,
		})

		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		fmt.Println("Saves repo created")
		os.Exit(0)
	}

	// ### MODE - UPLOAD ###
	if mode == "upload" {
		err = cmd.CommitNewestFile(config, client)

		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		fmt.Println("New save file uploaded")
		os.Exit(0)
	}

	// ### MODE - SYNC ###
	if mode == "sync" {
		err = cmd.SyncNewestFileToDevice(config, client)

		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		fmt.Println("Newest save file has been synced")
	}
}
