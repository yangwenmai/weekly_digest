package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

func printContributors(ctx context.Context, headDate, tailDate time.Time, commits []*github.RepositoryCommit) string {
	contributors := filterContributors(commits)
	return wrapContributorsReport(contributors, headDate, tailDate, formatWords())
}

func filterContributors(commits []*github.RepositoryCommit) []*github.Contributor {
	contributors := []*github.Contributor{}
	for _, r := range commits {
		contributor := &github.Contributor{}
		contributor.Login = r.Author.Login
		contributor.HTMLURL = r.Author.HTMLURL
		contributors = append(contributors, contributor)
	}
	return contributors
}

func wrapContributorsReport(contributors []*github.Contributor, headDate, tailDate time.Time, lastStr string) string {
	contributorsString := "# CONTRIBUTORS\n"
	if len(contributors) == 0 {
		contributorsString += fmt.Sprintf("Last %s (%s~%s), no contributors.\n",
			lastStr, headDate, tailDate)
	} else {
		contributorsString += fmt.Sprintf("Last %s (%s~%s) there were/was %d contributors.\n", lastStr, headDate, tailDate, len(contributors))
		sts := map[string][]*github.Contributor{}
		for _, item := range contributors {
			sts[item.GetLogin()] = append(sts[item.GetLogin()], item)
		}
		for _, item := range sts {
			contributorsString += fmt.Sprintf(":bust_in_silhouette: [%s](%s) [%d commits]\n",
				item[0].GetLogin(), item[0].GetHTMLURL(), len(item))
		}
	}
	return contributorsString
}
