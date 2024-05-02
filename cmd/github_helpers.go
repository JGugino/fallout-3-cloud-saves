package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v61/github"
)

func SyncNewestFileToDevice(config Config, client *github.Client) error {

	commits, _, err := client.Repositories.ListCommits(context.Background(), config.CommiterName, config.RepoName, nil)

	if err != nil {
		return err
	}

	var modTime time.Time
	var shas []string

	for i := 0; i < len(commits); i++ {
		commit := commits[i]

		if !commit.Commit.Committer.Date.Before(modTime) {
			if commit.Commit.Committer.Date.Time.After(modTime) {
				modTime = commit.Commit.Committer.GetDate().Time
				shas = shas[:0]
			}
			shas = append(shas, *commit.SHA)
		}
	}

	commit, _, err := client.Repositories.GetCommit(context.Background(), config.CommiterName, config.RepoName, shas[0], nil)

	if err != nil {
		return err
	}

	resp, err := http.Get(commit.Files[0].GetRawURL())

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	fmt.Println(resp.Body)

	// cmd := exec.Command("curl", "-O", "--header", "Authorization:", fmt.Sprintf("token %s", config.GithubApiKey))
	// cmd.Run()

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
