# Github RPM Sync

This package provides a container image that will sync a RPM repo with various Github repositories. It is designed
to accept release events from Github repository webhooks where an rpm file may be included in the release assets.

## Container Usage

This container is designed to be run behind a reverse proxy, such as Nginx or Cloudflare.

### Variables

|Variable|Required|Description|
|-|-|-|
|`GITHUB_USERNAME`|Yes|Your Github username. Needed to query the API and download assets.|
|`GITHUB_ACCESS_TOKEN`|Yes|Your Github access token. Needed to query the API and download assets.|
|`GITHUB_WEBHOOK_SECRET`|Yes|The Webhook secrets for incoming events.|
|`GITHUB_REPOS`|Yes|Comma separated list of username/repo for repos to query on startup.|
|`YUM_REPO_ID`|No|A short, alphanumberic identifier for your YUM repo.|
|`YUM_REPO_DESCRIPTION`|No|The description of your YUM repo.|
|`YUM_REPO_BASEURL`|No|The baseurl for your YUM repo.|

If any of `YUM_REPO_ID`, `YUM_REPO_DESCRIPTION`, `YUM_REPO_BASEURL` are specified then all three are required, and a `.repo` file will be generated.

### Ports

- 80
- 443

The certificate presented by the TLS server on port 443 is a self-signed certificate generated when the application starts.

## Webhook Usage

The container provides a webhook receptor at `/gh_webhook` that can be used to have the container automatically update with new releases containg RPM assets.

Add a new Webhook in the Github repository with the following properties:
- **Payload URL:** `https://<container>/gh_webhook`
- **Content Type:** `application/json`
- **Secret:** Same value as provided to the `GITHUB_WEBHOOK_SECRET` variable in the container
- **Events:** Releases
