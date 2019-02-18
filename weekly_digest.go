package main

import (
	"context"
	"flag"
	"fmt"
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

var (
	client *github.Client
)

const timeLayout = "2006-01-02 15:04:05"

func main() {
	flag.Parse()
	checkParams()

	tailDate, _ := time.Parse(timeLayout, *endDate)
	headDate := tailDate.Add(time.Duration(-*interval*24) * time.Hour)

	title := fmt.Sprintf("Weekly Digest (%d %s, %d - %d %s, %d)",
		headDate.Day(), headDate.Month(), headDate.Year(), tailDate.Day(), tailDate.Month(), tailDate.Year())

	ctx := context.Background()
	client = NewClient(ctx)
	pullRequestsString := printPullRequests(ctx, headDate, tailDate)
	body := ""
	bodyHead := fmt.Sprintf("Here's the **Weekly Digest** for [*%s/%s*](https://github.com/%s/%s):\n", *owner, *repo, *owner, *repo)
	body += bodyHead
	if len(pullRequestsString) > 0 {
		body += "\n --- \n"
		body += pullRequestsString
	}
	commitsString := printCommits(ctx, headDate, tailDate)
	if len(commitsString) > 0 {
		body += "\n --- \n"
		body += commitsString
	}
	body += "\n --- \n"
	body += "\n"

	// fmt.Println(title, body)

	// create a issue
	createIssue(title, body, "", []string{"weekly-digest"})
}

func NewClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
	return client
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
