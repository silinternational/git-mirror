package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sourcerepo/v1"
)

var (
	gitHubUrl            string
	bitbucketUrl         string
	googleProjectString  string
	googleCredentialFile string
	googleAdminEmail     string
)

func init() {
	gitHubUrl = os.Getenv("GITHUB_URL")
	bitbucketUrl = os.Getenv("BITBUCKET_URL")
	googleProjectString = os.Getenv("GOOGLE_PROJECT_STRING")
	googleCredentialFile = os.Getenv("GOOGLE_CREDENTIAL_FILE")
	googleAdminEmail = os.Getenv("GOOGLE_ADMIN_EMAIL")
}

var (
	ErrorMissingGitHubURL    = errors.New("missing GITHUB_URL value")
	ErrorMissingBitbucketURL = errors.New("missing BITBUCKET_URL value")
)

// RequestGitHubData returns GH repo data that has items listed in GHItems struct
func RequestGitHubData(page int) (GHJSONData, error) {
	if gitHubUrl == "" {
		return nil, ErrorMissingGitHubURL
	}

	res, err := http.Get(fmt.Sprintf("%s?sort=full_name&per_page=100&page=%d", gitHubUrl, page))
	if err != nil {
		return GHJSONData{}, err
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return GHJSONData{}, err
	}

	if res.StatusCode != http.StatusOK {
		return GHJSONData{}, fmt.Errorf("error returned from GitHub: %s", body)
	}

	data := GHJSONData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return GHJSONData{}, err
	}

	return data, nil
}

// RequestBitBucketData returns BB repo data that has items listed in GHItems struct
func RequestBitBucketData(page int) (BBJSONData, error) {
	if bitbucketUrl == "" {
		return BBJSONData{}, ErrorMissingBitbucketURL
	}

	res, err := http.Get(fmt.Sprintf("%s?page=%d", bitbucketUrl, page))
	if err != nil {
		return BBJSONData{}, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return BBJSONData{}, err
	}

	if res.StatusCode != http.StatusOK {
		return BBJSONData{}, fmt.Errorf("error returned from Bitbucket: %s", body)
	}

	data := BBJSONData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return BBJSONData{}, err
	}

	return data, nil
}

// RequestGoogleData gets a slice of all available repositories on the Google project
func RequestGoogleData() ([]*sourcerepo.Repo, error) {
	if googleProjectString == "" || googleCredentialFile == "" || googleAdminEmail == "" {
		return nil, fmt.Errorf("Google configuration is incomplete")
	}

	ctx := context.Background()
	b, err := os.ReadFile(googleCredentialFile)
	if err != nil {
		return nil, err
	}
	config, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx)

	sourcerepoService, err := sourcerepo.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	resp, err := sourcerepoService.Projects.Repos.List(googleProjectString).Do()
	if err != nil {
		return nil, err
	}

	return resp.Repos, nil
}

func CreateGoogleMirror(repo sourcerepo.Repo) (*sourcerepo.Repo, error) {
	ctx := context.Background()
	b, err := os.ReadFile(googleCredentialFile)
	if err != nil {
		return &sourcerepo.Repo{}, err
	}
	config, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return &sourcerepo.Repo{}, err
	}

	client := config.Client(ctx)
	config.Subject = googleAdminEmail

	sourcerepoService, err := sourcerepo.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return &sourcerepo.Repo{}, err
	}

	mirror, err := sourcerepoService.Projects.Repos.Create(googleProjectString, &sourcerepo.Repo{
		Name:         repo.Name,
		MirrorConfig: repo.MirrorConfig,
	}).Do()
	if err != nil {
		return &sourcerepo.Repo{}, err
	}

	return mirror, nil
}

// Item is the single repository data structure
type GHItem struct {
	ID      int
	Name    string
	HTMLUrl string `json:"html_url"`
	// Owner           Owner
	// Description     string
	CreatedAt string `json:"created_at"`
	// StargazersCount int    `json:"stargazers_count"`
	UpdatedAt string `json:"updated_at"`
}

// GHJSONData contains the GitHub API response
type GHJSONData []GHItem

// CloneLink will have the HTTPS BitBucket clone link to the repo
type CloneLink struct {
	HREF string
}

// Links is the list of links for a BitBucket Repo
type Links struct {
	Clone []CloneLink
}

// BBItem is a single repository data structure for Bit Bucket.
type BBItem struct {
	Name      string
	Links     Links
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}

type BBResponse struct {
	Values []BBItem
}

// BBJSONData contains the BitBucket API response
type BBJSONData BBResponse
