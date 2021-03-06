package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

func printCommits(ctx context.Context, headDate, tailDate time.Time) (string, []*github.RepositoryCommit) {
	commits := fetchCommits(ctx, headDate)
	commits = filterCommits(commits, headDate, tailDate)
	return wrapCommitsReport(commits, headDate, tailDate, formatWords()), commits
}

func fetchCommits(ctx context.Context, headDate time.Time) []*github.RepositoryCommit {
	// FIXME:暂时只取 200 个 commits（有可能条件范围内的 commits 是超过 200 的），后续再优化
	listOpts := github.ListOptions{PerPage: 200}
	opts := &github.CommitsListOptions{SHA: "master", Since: headDate, ListOptions: listOpts}

	commits, _, err := client.Repositories.ListCommits(ctx, *owner, *repo, opts)
	if err != nil {
		panic(err)
	}
	return commits
}

func filterCommits(commits []*github.RepositoryCommit, headDate, tailDate time.Time) []*github.RepositoryCommit {
	intervalCommits := []*github.RepositoryCommit{}

	for _, commit := range commits {
		if commit.Author != nil && commit.Commit.GetCommitter().Date.After(headDate) && commit.Commit.GetCommitter().Date.Before(tailDate) {
			intervalCommits = append(intervalCommits, commit)
		}
	}
	return intervalCommits
}

func wrapCommitsReport(commits []*github.RepositoryCommit, headDate, tailDate time.Time, lastStr string) string {
	commitsString := "# COMMITS\n"
	if len(commits) == 0 {
		commitsString += fmt.Sprintf("Last %s (%s~%s), no commits.\n",
			lastStr, headDate, tailDate)
	} else {
		commitsString += fmt.Sprintf("Last %s (%s~%s) there were/was %d commits.\n", lastStr, headDate, tailDate, len(commits))
		for _, item := range commits {
			commitsString += fmt.Sprintf(":hammer_and_wrench: [%s](%s) by [%s](%s)\n",
				strings.Replace(item.Commit.GetMessage(), "\n", " ", -1), item.GetHTMLURL(), item.GetAuthor().GetLogin(), item.GetAuthor().GetHTMLURL())
		}
	}
	return commitsString
}
