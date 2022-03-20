package main

type GithubWebhookReleaseType struct {
	Action     string               `json:"action"`
	Release    GithubReleaseType    `json:"release"`
	Repository GithubRepositoryType `json:"repository"`
}

type GithubReleaseType struct {
	URL             string            `json:"url"`
	AssetsURL       string            `json:"assets_url"`
	UploadURL       string            `json:"upload_url"`
	HTMLURL         string            `json:"html_url"`
	ID              int               `json:"id"`
	NodeID          string            `json:"node_id"`
	TagName         string            `json:"tag_name"`
	TargetCommitish string            `json:"target_commitish"`
	Name            string            `json:"name"`
	Draft           bool              `json:"draft"`
	Prerelease      bool              `json:"prerelease"`
	CreatedAt       string            `json:"created_at"`
	PublishedAt     string            `json:"published_at"`
	Assets          []GithubAssetType `json:"assets"`
	TarballURL      string            `json:"tarbURL"`
	ZipballURL      string            `json:"zipball_url"`
	Body            string            `json:"body"`
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
	NodeID             string `json:"node_id"`
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	State              string `json:"state"`
	Size               uint64 `json:"size"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	BrowserDownloadURL string `json:"browser_download_url"`
}
