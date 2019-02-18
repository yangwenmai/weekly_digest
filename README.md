# weekly_digest

Weekly digest on your GitHub repository 📆

Inspired by [probot/weekly-digest](https://github.com/probot/weekly-digest).

## Features

When you get the Weekly Digest, and installed it, you can get this information:

- Issues created in the last one week
    - Open Issues
    - Closed Issues
    - Noisy Issue
    - Liked Issue
- Pull requests opened, updated, or merged in the last pull request
    - Opened Pull Requests
    - Updated Pull Requests
    - Merged Pull Requests
- Commits made in the master branch, in the last week
- Contributors, adding contributions in the last week
- Stargazers, or the fans of your repositories, who really loved your repo
- Releases, of the project you are working on

## usage

```bash
$ git clone github.com/yangwenmai/weekly_digest
$ cd weekly_digest
$ go run github_weekly_disgest.go --access_token=<your personal github access_token> --owner=yangwenmai --repo=weekly_digest --end_date="2019-02-19 08:00:00" --interval=7
```

## Reference

1. [Github Pull Request API](https://developer.github.com/v3/pulls/)
2. [Create Your Personal Access Token](https://github.com/settings/tokens/new)

## Credits 

This project is developed and maintained by [yangwenmai](https://github.com/yangwenmai).

I would like to thanks [abhijeetps](https://github.com/abhijeetps) for this project. ❤️
