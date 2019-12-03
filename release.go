package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

func printReleases(ctx context.Context, headDate, tailDate time.Time) string {
	releases := fetchReleases(ctx, headDate)
	releases = filterReleases(releases, headDate, tailDate)
	return wrapReleasesReport(releases, headDate, tailDate, formatWords())
}

func fetchReleases(ctx context.Context, headDate time.Time) []*github.RepositoryRelease {
	// FIXME:暂时只取 100 个 Release（有可能30天的 Release 是超过 Release 的），后续再优化
	listOpts := &github.ListOptions{PerPage: 100}
	releases, _, err := client.Repositories.ListReleases(ctx, *owner, *repo, listOpts)
	if err != nil {
		panic(err)
	}
	return releases
}

func filterReleases(releases []*github.RepositoryRelease, headDate, tailDate time.Time) []*github.RepositoryRelease {
	intervalReleases := []*github.RepositoryRelease{}

	for _, r := range releases {
		if r.GetPublishedAt().After(headDate) && r.GetPublishedAt().Before(tailDate) {
			intervalReleases = append(intervalReleases, r)
		}
	}
	return intervalReleases
}

func wrapReleasesReport(releases []*github.RepositoryRelease, headDate, tailDate time.Time, lastStr string) string {
	releaseString := "# RELEASES\n"
	if len(releases) == 0 {
		releaseString += fmt.Sprintf("Last %s (%s~%s), no releases.\n",
			lastStr, headDate, tailDate)
	} else {
		releaseString += fmt.Sprintf("Last %s (%s~%s) there were/was %d releases.\n", lastStr, headDate, tailDate, len(releases))
		for _, item := range releases {
			releaseString += fmt.Sprintf(":rocket: [%s - (%s)](%s) by [%s](%s)\n",
				strings.Replace(item.GetTagName(), "\n", " ", -1),
				strings.Replace(item.GetName(), "\n", " ", -1),
				item.GetHTMLURL(),
				item.GetAuthor().GetLogin(),
				item.GetAuthor().GetHTMLURL())
			releaseString += fmt.Sprintf("%s\n", item.GetBody())
		}
	}
	return releaseString
}
