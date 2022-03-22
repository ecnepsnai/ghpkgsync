package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/ecnepsnai/logtic"
)

func downloadAsset(repo GithubRepositoryType, asset GithubAssetType) (bool, error) {
	traceStart := time.Now()

	if !strings.HasSuffix(asset.Name, ".rpm") {
		return false, nil
	}

	assetPath := path.Join("repo", asset.Name)
	if fileExists(assetPath) {
		log.PDebug("Asset already downloaded", map[string]interface{}{
			"path": assetPath,
		})
		return false, nil
	}

	url := fmt.Sprintf("%s/releases/assets/%d", repo.APIURL, asset.ID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Accept", "application/octet-stream")
	req.SetBasicAuth(Config.GithubUsername, Config.GithubAccessToken)

	log.PDebug("Downloading asset", map[string]interface{}{
		"url":  url,
		"path": assetPath,
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("http %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(assetPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return false, err
	}
	defer f.Close()

	len, err := io.Copy(f, resp.Body)
	if err != nil {
		return false, err
	}

	if uint64(len) != asset.Size {
		return false, fmt.Errorf("bad size")
	}

	log.PInfo("Downloaded asset", map[string]interface{}{
		"url":      url,
		"path":     assetPath,
		"duration": time.Since(traceStart).String(),
		"size":     logtic.FormatBytesB(uint64(len)),
		"size_b":   len,
	})

	return true, nil
}

func downloadReleaseAssets(repo GithubRepositoryType, release GithubReleaseType) {
	wg := sync.WaitGroup{}
	needSync := false
	wg.Add(len(release.Assets))
	for _, asset := range release.Assets {
		go func(r GithubRepositoryType, a GithubAssetType) {
			didDownload, err := downloadAsset(r, a)
			if err != nil {
				log.PError("Error downloading asset", map[string]interface{}{
					"asset": a.ID,
					"error": err.Error(),
				})
			}
			if didDownload {
				needSync = true
			}
			wg.Done()
		}(repo, asset)
	}
	wg.Wait()
	if needSync {
		syncRepo()
	}
}
