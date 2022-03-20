package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/web/router"
)

var log = logtic.Log.Connect("ghrpnsync")

var Config struct {
	GithubUsername    string
	GithubAccessToken string
	WebhookSecret     string
	GithubRepos       []string
}

func main() {
	logtic.Log.Level = logtic.LevelDebug
	logtic.Log.Open()

	githubUsername := os.Getenv("GITHUB_USERNAME")
	githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	githubRepos := strings.Split(os.Getenv("GITHUB_REPOS"), ",")

	webhookSecret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if webhookSecret == "" {
		fmt.Fprintf(os.Stderr, "Environment variable GITHUB_WEBHOOK_SECRET required\n")
		os.Exit(1)
	}

	Config.GithubUsername = githubUsername
	Config.GithubAccessToken = githubAccessToken
	Config.WebhookSecret = webhookSecret
	Config.GithubRepos = githubRepos

	log.PInfo("Started ghrpmsync", map[string]interface{}{
		"GithubUsername":    Config.GithubUsername,
		"GithubAccessToken": Config.GithubAccessToken,
		"WebhookSecret":     Config.WebhookSecret,
		"GithubRepos":       Config.GithubRepos,
	})

	PullRepos()

	server := router.New()
	server.ServeFiles("repo", "/rpm")
	server.Handle("POST", "/gh_webhook", acceptWebhook)

	go func(s *router.Server) {
		if err := startHTTPS(s); err != nil {
			log.Panic("error starting https server: %s", err.Error())
		}
	}(server)
	go func(s *router.Server) {
		if err := startHTTP(s); err != nil {
			log.Panic("error starting http server: %s", err.Error())
		}
	}(server)
	for true {
		time.Sleep(1 * time.Minute)
	}
}

func downloadAsset(repo GithubRepositoryType, asset GithubAssetType) error {
	traceStart := time.Now()

	if !strings.HasSuffix(asset.Name, ".rpm") {
		return nil
	}

	assetPath := path.Join("repo", asset.Name)
	if fileExists(assetPath) {
		log.PDebug("Asset already downloaded", map[string]interface{}{
			"path": assetPath,
		})
		return nil
	}

	url := fmt.Sprintf("%s/releases/assets/%d", repo.APIURL, asset.ID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/octet-stream")
	req.SetBasicAuth(Config.GithubUsername, Config.GithubAccessToken)

	log.PDebug("Downloading asset", map[string]interface{}{
		"url":  url,
		"path": assetPath,
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("http %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(assetPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	len, err := io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	if uint64(len) != asset.Size {
		return fmt.Errorf("bad size")
	}

	log.PInfo("Downloaded asset", map[string]interface{}{
		"url":      url,
		"path":     assetPath,
		"duration": time.Since(traceStart).String(),
		"size":     logtic.FormatBytesB(uint64(len)),
		"size_b":   len,
	})

	return nil
}

func downloadReleaseAssets(repo GithubRepositoryType, release GithubReleaseType) {
	wg := sync.WaitGroup{}
	wg.Add(len(release.Assets))
	for _, asset := range release.Assets {
		go func(r GithubRepositoryType, a GithubAssetType) {
			if err := downloadAsset(r, a); err != nil {
				log.PError("Error downloading asset", map[string]interface{}{
					"asset": a.ID,
					"error": err.Error(),
				})
			}
			wg.Done()
		}(repo, asset)
	}
	wg.Wait()
	syncRepo()
}

func syncRepo() {
	cmd := exec.Command("/usr/bin/createrepo_c", ".")
	cmd.Dir = "repo"
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.PError("Error syncing repository", map[string]interface{}{
			"error":  err.Error(),
			"output": string(output),
		})
		return
	}
	log.Info("Synced repository")
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
