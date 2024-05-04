package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v61/github"
)

const (
	TMP_REPO_PATH = "./tmp/fallout-saves"
)

type NewestFile struct {
	Name string `json:"name"`
}

func SyncNewestFileToDevice(config Config, client *github.Client) error {

	repo, _, err := client.Repositories.Get(context.Background(), config.CommiterName, config.RepoName)

	if err != nil {
		return err
	}

	contents, err := ReadWholeFile(TMP_REPO_PATH, "newest.json")

	if err != nil {
		return err
	}

	var newestJson NewestFile

	json.Unmarshal(contents, &newestJson)

	_, err = git.PlainClone(TMP_REPO_PATH, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: config.CommiterName,
			Password: config.GithubApiKey,
		},
		URL:      repo.GetCloneURL(),
		Progress: os.Stdout,
	})

	if err != nil {
		repo, err := git.PlainOpen(TMP_REPO_PATH)

		if err != nil {
			return err
		}

		workTree, err := repo.Worktree()
		if err != nil {
			return err
		}

		workTree.Pull(&git.PullOptions{
			Auth: &http.BasicAuth{
				Username: config.CommiterName,
				Password: config.GithubApiKey,
			},
			RemoteName: "origin",
		})

		fmt.Println("Temp files updated")
	}

	if err == nil {
		fmt.Println("Saves repo cloned to temp")
	}

	_, err = Copy(path.Join(TMP_REPO_PATH, newestJson.Name), path.Join(config.SaveLocation, newestJson.Name))

	return err
}

func CommitNewestFile(config Config, client *github.Client) error {
	newestFile, err := GetNewestFiles(config.SaveLocation)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	_, err = Copy(path.Join(config.SaveLocation, newestFile[0]), path.Join(TMP_REPO_PATH, newestFile[0]))

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	newest := NewestFile{
		Name: newestFile[0],
	}

	contents, err := json.Marshal(newest)

	if err != nil {
		return err
	}

	err = CreateNewFile(TMP_REPO_PATH, "newest.json", string(contents))

	if err != nil {
		return err
	}

	repo, err := git.PlainOpen(TMP_REPO_PATH)

	if err != nil {
		return err
	}

	err = repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: config.CommiterName,
			Password: config.GithubApiKey,
		},
		RemoteName: "origin",
	})

	return err
}
