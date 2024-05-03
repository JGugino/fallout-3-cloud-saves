package cmd

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v61/github"
)

func SyncNewestFileToDevice(config Config, client *github.Client) error {

	repo, _, err := client.Repositories.Get(context.Background(), config.CommiterName, config.RepoName)

	if err != nil {
		return err
	}

	_, err = git.PlainClone("./tmp/fallout-saves", false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: config.CommiterName,
			Password: config.GithubApiKey,
		},
		URL:      repo.GetCloneURL(),
		Progress: os.Stdout,
	})

	if err != nil {
		repo, err := git.PlainOpen("./tmp/fallout-saves")

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

	newestFile, err := GetNewestFiles("./tmp/fallout-saves")

	if err != nil {
		return err
	}

	_, err = Copy(path.Join("./tmp/fallout-saves", newestFile[0]), path.Join(config.SaveLocation, newestFile[0]))

	return err
}

func CommitNewestFile(config Config, client *github.Client) error {
	newestFile, err := GetNewestFiles(config.SaveLocation)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	saveContents, err := ReadWholeFile(config.SaveLocation, "/"+newestFile[0])

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	saveFile := &github.RepositoryContentFileOptions{
		Message: github.String("Cloud Save Sync - Create"),
		Committer: &github.CommitAuthor{
			Name:  github.String(config.CommiterName),
			Email: github.String(config.CommiterEmail),
		},
		Content: saveContents,
	}

	_, _, err = client.Repositories.CreateFile(context.Background(), config.CommiterName, config.RepoName, newestFile[0], saveFile)

	return err
}

func CommitUpdatedFile(config Config, client *github.Client) error {
	newestFile, err := GetNewestFiles(config.SaveLocation)

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	saveContents, err := ReadWholeFile(config.SaveLocation, "/"+newestFile[0])

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	//TODO: Finish updating

	commitInfo, _, _ := client.Repositories.ListCommits(context.Background(), config.CommiterName, config.RepoName, nil)

	fmt.Println(commitInfo[0])

	saveFile := &github.RepositoryContentFileOptions{
		Message: github.String("Cloud Save Sync - Update"),
		Committer: &github.CommitAuthor{
			Name:  github.String(config.CommiterName),
			Email: github.String(config.CommiterEmail),
		},
		Content: saveContents,
	}

	_, _, err = client.Repositories.UpdateFile(context.Background(), config.CommiterName, config.RepoName, newestFile[0], saveFile)

	return err
}
