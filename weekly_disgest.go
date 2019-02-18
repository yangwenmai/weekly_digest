package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	accessToken = flag.String("access_token", "", "Github API access token")
	owner       = flag.String("owner", "", "Github owner name")
	repo        = flag.String("repo", "", "Github owner's repo")
	endDate     = flag.String("end_date", "", "end date")
	interval    = flag.Int64("interval", 7, "past days base on end date")
)

const timeLayout = "2006-01-02 15:04:05"

func main() {
	flag.Parse()
	checkParams()

	tailDate, _ := time.Parse(timeLayout, *endDate)
	headDate := tailDate.Add(time.Duration(-*interval*24) * time.Hour)

	prs := genPRS()
	intervalPRs := filterPRS(prs, headDate, tailDate)
	pullRequestsString := printPRSReport(intervalPRs, headDate, tailDate, formatWords())
	
	// create a issue
	title := fmt.Sprintf("Weekly Digest (%d %s, %d - %d %s, %d)",
		headDate.Day(), headDate.Month(), headDate.Year(), tailDate.Day(), tailDate.Month(), tailDate.Year())

	body := ""
	bodyHead := fmt.Sprintf("Here's the **Weekly Digest** for [*%s/%s*](https://github.com/%s/%s):\n", *owner, *repo, *owner, *repo)

	body += bodyHead

	if len(pullRequestsString) > 0 {
		body += "\n --- \n"
		body += pullRequestsString
	}
	body += "\n --- \n"
	body += "\n"

	createIssue(title, body, "", []string{"weekly-digest"})
}

func checkParams() {
	if *accessToken == "" {
		panic("Github API access token can not be empty.")
	}
	if *owner == "" {
		panic("Github owner name can not be empty.")
	}
	if *repo == "" {
		panic("Github owner's repo can not be empty.")
	}
	if *endDate == "" {
		*endDate = time.Now().Format(timeLayout)
	}
	if *interval > 30 {
		panic("the max interval value is 30")
	}
}

func formatWords() string {
	lastStr := ""
	switch *interval {
	case 7:
		lastStr = "weeks"
		break
	case 14:
		lastStr = "2 weeks"
		break
	case 21:
		lastStr = "3 weeks"
		break
	case 30:
		lastStr = "months"
		break
	default:
		lastStr = fmt.Sprintf("%d days", *interval)
	}
	return lastStr
}
func genPRS() []*github.PullRequest {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// FIXME:暂时只取 100 个PR（有可能30天的PR是超过100的），后续再优化
	listOpts := github.ListOptions{PerPage: 100}
	opts := &github.PullRequestListOptions{State: "all", ListOptions: listOpts}

	prs, _, err := client.PullRequests.List(ctx, *owner, *repo, opts)
	if err != nil {
		panic(err)
	}

	return prs
}

func createIssue(title, body, assignee string, labels []string) *github.Issue {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	_ = client
	input := &github.IssueRequest{
		Title:    &title,
		Body:     &body,
		Assignee: &assignee,
		Labels:   &labels,
	}
	
	issue, _, err := client.Issues.Create(ctx, *owner, *repo, input)
	if err != nil {
		panic(err)
	}
	return issue
}

func filterPRS(prs []*github.PullRequest, headDate, tailDate time.Time) []*github.PullRequest {
	intervalPRs := []*github.PullRequest{}
	for _, pr := range prs {
		if pr.CreatedAt.After(headDate) && pr.CreatedAt.Before(tailDate) && pr.GetState() == "open" && pr.MergedAt == nil {
			intervalPRs = append(intervalPRs, pr)
			continue
		}
		if pr.UpdatedAt.After(headDate) && pr.UpdatedAt.Before(tailDate) && pr.GetState() == "open" && pr.MergedAt == nil {
			intervalPRs = append(intervalPRs, pr)
			continue
		}
		if pr.MergedAt != nil && pr.MergedAt.After(headDate) && pr.MergedAt.Before(tailDate) && pr.GetState() == "closed" {
			intervalPRs = append(intervalPRs, pr)
			continue
		}
	}
	return intervalPRs
}

func printPRSReport(intervalPRs []*github.PullRequest, headDate, tailDate time.Time, lastStr string) string {
	pullRequestsString := "# PULL REQUESTS\n"
	if len(intervalPRs) == 0 {
		pullRequestsString += fmt.Sprintf("Last %s (%s~%s), no pull requests were created, updated or merged.\n",
			lastStr, headDate, tailDate)
	} else {
		pullRequestsString += fmt.Sprintf("Last %s (%s~%s), %d pull request's were created, updated or merged.\n",
			lastStr, headDate, tailDate, len(intervalPRs))
		openPullRequest := []*github.PullRequest{}
		for _, item := range intervalPRs {
			if item.CreatedAt.After(headDate) && item.CreatedAt.Before(tailDate) && item.CreatedAt == item.UpdatedAt && item.MergedAt == nil {
				openPullRequest = append(openPullRequest, item)
			}
		}
		updatedPullRequest := []*github.PullRequest{}
		for _, item := range intervalPRs {
			if item.UpdatedAt.After(headDate) && item.CreatedAt.Before(tailDate) && item.CreatedAt != item.UpdatedAt && item.MergedAt == nil {
				updatedPullRequest = append(updatedPullRequest, item)
			}
		}
		mergedPullRequest := []*github.PullRequest{}
		for _, item := range intervalPRs {
			if item.MergedAt != nil && item.MergedAt.After(headDate) && item.MergedAt.Before(tailDate) {
				mergedPullRequest = append(mergedPullRequest, item)
			}
		}

		mergedPullRequestString := ""
		openPullRequestString := ""
		updatedPullRequestString := ""
		if len(mergedPullRequest) > 0 {
			mergedPullRequestString = "## MERGED PULL REQUEST\n"
			mergedPullRequestString += fmt.Sprintf("Last %s, %d pull request were/was merged.\n",
				lastStr, len(mergedPullRequest))
			for _, item := range mergedPullRequest {
				mergedPullRequestString += fmt.Sprintf(":purple_heart: #%d [%s](%s) merged at %s, by [%s](%s)\n",
					item.GetNumber(), strings.Replace(item.GetTitle(), "\n", " ", -1), item.GetHTMLURL(),
					item.GetMergedAt(), item.GetUser().GetLogin(), item.GetUser().GetHTMLURL())
			}
		}
		if len(openPullRequest) > 0 {
			openPullRequestString = "## OPEN PULL REQUEST\n"
			openPullRequestString += fmt.Sprintf("Last %s, %d pull request were/was opened.\n", lastStr, len(openPullRequest))
			for _, item := range openPullRequest {
				openPullRequestString += fmt.Sprintf(":green_heart: #%d [%s](%s) opened at %s, by [%s](%s)\n",
					item.GetNumber(), strings.Replace(item.GetTitle(), "\n", " ", -1), item.GetHTMLURL(),
					item.GetCreatedAt(), item.GetUser().GetLogin(), item.GetUser().GetHTMLURL())
			}
		}
		if len(updatedPullRequest) > 0 {
			updatedPullRequestString = "## UPDATED PULL REQUEST\n"
			updatedPullRequestString += fmt.Sprintf("Last %s, %d pull request were/was updated.\n", lastStr, len(updatedPullRequest))
			for _, item := range updatedPullRequest {
				updatedPullRequestString += fmt.Sprintf(":yellow_heart: #%d [%s](%s) updated at %s, by [%s](%s)\n",
					item.GetNumber(), strings.Replace(item.GetTitle(), "\n", " ", -1), item.GetHTMLURL(),
					item.GetUpdatedAt(), item.GetUser().GetLogin(), item.GetUser().GetHTMLURL())
			}
		}
		if len(openPullRequestString) > 0 {
			pullRequestsString += openPullRequestString
		}
		if len(updatedPullRequestString) > 0 {
			pullRequestsString += updatedPullRequestString
		}
		if len(mergedPullRequestString) > 0 {
			pullRequestsString += mergedPullRequestString
		}
	}
	return pullRequestsString
}
