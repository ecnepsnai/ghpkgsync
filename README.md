# Github RPM Sync

This package provides a container image that will sync a RPM repo with various Github repositories. It is designed
to accept release events from Github repository webhooks where an rpm file may be included in the release assets.

## Usage

This container is designed to be run behind a reverse proxy, such as Nginx or Cloudflare.

### Variables

|Variable|Required|Description|
|-|-|-|
|GITHUB_USERNAME|Yes|Your Github username. Needed to query the API and download assets.|
|GITHUB_ACCESS_TOKEN|Yes|Your Github access token. Needed to query the API and download assets.|
|GITHUB_WEBHOOK_SECRET|Yes|The Webhook secrets for incoming events.|
|GITHUB_REPOS|Yes|Comma separated list of <username>/<repo> for repos to query on startup|
|YUM_REPO_ID|Yes|A short, alphanumberic identifier for your YUM repo|
|YUM_REPO_DESCRIPTION|Yes|The description of your YUM repo|
|YUM_REPO_BASEURL|Yes|The baseurl for your YUM repo|

### Ports

- 80
- 443
