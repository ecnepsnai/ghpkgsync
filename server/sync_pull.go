package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func PullRepos() {
	for _, repo := range Config.GithubRepos {
		PullRepo(repo)
	}
}

func PullRepo(repo string) error {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s", repo)

	req, err := http.NewRequest("GET", apiURL+"/releases", nil)
	if err != nil {
		log.PError("", map[string]interface{}{
			"repo":  repo,
			"error": err.Error(),
		})
		return err
	}
	req.Header.Set("User-Agent", fmt.Sprint("ghrpmsync"))
	req.SetBasicAuth(Config.GithubUsername, Config.GithubAccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.PError("", map[string]interface{}{
			"repo":  repo,
			"error": err.Error(),
		})
		return err
	}
	if resp.StatusCode != 200 {
		log.PError("", map[string]interface{}{
			"repo":  repo,
			"error": fmt.Sprintf("http %d", resp.StatusCode),
		})
		return err
	}
	results := []GithubReleaseType{}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		log.PError("", map[string]interface{}{
			"repo":  repo,
			"error": err.Error(),
		})
		return err
	}

	for _, release := range results {
		downloadReleaseAssets(GithubRepositoryType{APIURL: apiURL}, release)
	}

	return nil
}
