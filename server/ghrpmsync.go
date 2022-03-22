package main

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/web/router"
)

var log = logtic.Log.Connect("ghrpnsync")

var Config struct {
	GithubUsername      string
	GithubAccessToken   string
	GithubWebhookSecret string
	GithubRepos         []string
	YumRepoID           string
	YumRepoDescription  string
	YumRepoBaseurl      string
}

func main() {
	logtic.Log.Level = logtic.LevelDebug
	logtic.Log.Open()

	Config.GithubUsername = getEnv("GITHUB_USERNAME", true)
	Config.GithubAccessToken = getEnv("GITHUB_ACCESS_TOKEN", true)
	Config.GithubWebhookSecret = getEnv("GITHUB_WEBHOOK_SECRET", true)
	Config.GithubRepos = strings.Split(getEnv("GITHUB_REPOS", false), ",")
	Config.YumRepoID = getEnv("YUM_REPO_ID", true)
	Config.YumRepoDescription = getEnv("YUM_REPO_DESCRIPTION", true)
	Config.YumRepoBaseurl = getEnv("YUM_REPO_BASEURL", true)
	log.PInfo("Started ghrpmsync", map[string]interface{}{
		"GithubUsername":      Config.GithubUsername,
		"GithubAccessToken":   Config.GithubAccessToken,
		"GithubWebhookSecret": Config.GithubWebhookSecret,
		"GithubRepos":         Config.GithubRepos,
		"YumRepoID":           Config.YumRepoID,
		"YumRepoDescription":  Config.YumRepoDescription,
		"YumRepoBaseurl":      Config.YumRepoBaseurl,
	})

	info, err := os.Stat("repo")
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir("repo", os.ModePerm); err != nil {
				log.PPanic("error creating repo directory", map[string]interface{}{
					"error": err.Error(),
				})
			}
			log.Info("Making repo directory")
		} else {
			log.PPanic("repo directory not accessible", map[string]interface{}{
				"error": err.Error(),
			})
		}
	} else if !info.IsDir() {
		log.Panic("repo directory is not a directory")
	}

	makeRepoFile()
	pullRepos()

	runtime.GC()

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

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func getEnv(name string, required bool) string {
	value := os.Getenv(name)
	if value == "" && required {
		log.Fatal("Missing required environment variable: %s", name)
	}
	return value
}
