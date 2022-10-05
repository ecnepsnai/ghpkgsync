package main

type GithubWebhookReleaseType struct {
	Action     string               `json:"action"`
	Release    GithubReleaseType    `json:"release"`
	Repository GithubRepositoryType `json:"repository"`
}

type GithubReleaseType struct {
	URL       string            `json:"url"`
	AssetsURL string            `json:"assets_url"`
	UploadURL string            `json:"upload_url"`
	HTMLURL   string            `json:"html_url"`
	ID        int               `json:"id"`
	TagName   string            `json:"tag_name"`
	Name      string            `json:"name"`
	Assets    []GithubAssetType `json:"assets"`
}

type GithubRepositoryType struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	APIURL   string `json:"url"`
}

type GithubAssetType struct {
	URL                string `json:"url"`
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	Size               uint64 `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}
