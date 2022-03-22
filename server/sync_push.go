package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ecnepsnai/web/router"
)

func acceptWebhook(w http.ResponseWriter, r router.Request) {
	log.PDebug("Webhook recieved", map[string]interface{}{
		"method":    r.HTTP.Method,
		"path":      r.HTTP.URL.Path,
		"real_ip":   w.Header().Get("X-Real-IP"),
		"conn_addr": r.HTTP.RemoteAddr,
	})

	githubEvent := r.HTTP.Header.Get("X-GitHub-Event")
	if githubEvent != "release" {
		log.PWarn("Ignoring request", map[string]interface{}{
			"X-Github-Event": githubEvent,
		})
		w.WriteHeader(204)
		return
	}

	body, err := io.ReadAll(r.HTTP.Body)
	if err != nil {
		log.PError("Error reading body", map[string]interface{}{
			"error": err.Error(),
		})
		w.WriteHeader(400)
		return
	}

	hasher := hmac.New(sha256.New, []byte(Config.GithubWebhookSecret))
	if _, err := hasher.Write(body); err != nil {
		w.WriteHeader(500)
		return
	}
	actualHash := strings.ToLower(fmt.Sprintf("sha256=%x", hasher.Sum(nil)))
	providedHash := strings.ToLower(r.HTTP.Header.Get("X-Hub-Signature-256"))
	if actualHash != providedHash {
		log.PWarn("Invalid webhook signature", map[string]interface{}{
			"expected_hash": actualHash,
			"provided_hash": providedHash,
		})
		w.WriteHeader(401)
		return
	}

	event := GithubWebhookReleaseType{}
	if err := json.Unmarshal(body, &event); err != nil {
		log.PError("Error parsing JSON", map[string]interface{}{
			"error": err.Error(),
		})
		w.WriteHeader(400)
		return
	}

	log.PInfo("New release event", map[string]interface{}{
		"repo":    event.Repository.FullName,
		"release": event.Release.Name,
	})

	go downloadReleaseAssets(event.Repository, event.Release)
	w.WriteHeader(204)
}
