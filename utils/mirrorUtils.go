package utils

import (
	"os"
	"strings"

	"google.golang.org/api/sourcerepo/v1"
)

const (
	excludedFile string = "excluded.txt" // Replace with the text file of URLs for the repos that should be mirrored
)

// Returns a list of repos that aren't in the excluded file and
// are not in the google project
func GetNonMirroredRepos() ([]sourcerepo.Repo, error) {
	// Retrieve the repositories from the google project
	googleRepos, err := RequestGoogleData()
	if err != nil {
		return []sourcerepo.Repo{}, err
	}

	// Retrieve repositories from Github
	githubRepos := GHJSONData{}
	for i := 1; true; i++ {
		data, err := RequestGitHubData(i)
		if err != nil {
			return []sourcerepo.Repo{}, err
		}

		if len(data) == 0 || i > 5 {
			break
		}

		githubRepos = append(githubRepos, data...)
	}

	// Retrieve repo data from Bit Bucket
	bitbucketRepos := BBJSONData{}

	for i := 1; true; i++ {
		data, err := RequestBitBucketData(i)
		if err != nil {
			return []sourcerepo.Repo{}, err
		}

		if len(bitbucketRepos.Values) == 0 || i > 5 {
			break
		}

		bitbucketRepos.Values = append(bitbucketRepos.Values, data.Values...)
	}
	if err != nil {
		return []sourcerepo.Repo{}, err
	}

	// Retrieve a list of repos that must be excluded from mirroring
	body, err := os.ReadFile(excludedFile)
	if err != nil {
		return []sourcerepo.Repo{}, err
	}

	excludedRepos := strings.Split(string(body), "\n")

	var allRepos []sourcerepo.Repo

	// Iterate through the github and bitbucket repos, adding their repo item
	// to the allRepos slice
	for _, repo := range githubRepos {
		mirrorConfig := sourcerepo.MirrorConfig{}
		mirrorConfig.Url = repo.HTMLUrl
		mirrorConfig.DeployKeyId = "" // Code needed to get DeployKeyID
		mirrorConfig.WebhookId = ""   // Code needed to get WebhookID
		allRepos = append(allRepos, sourcerepo.Repo{Name: repo.Name, MirrorConfig: &mirrorConfig})
	}

	for _, repo := range bitbucketRepos.Values {
		// Create and set the values for google mirror config
		mirrorConfig := sourcerepo.MirrorConfig{}
		mirrorConfig.Url = repo.Links.Clone[0].HREF
		mirrorConfig.DeployKeyId = "" // Code needed to get DeployKeyID
		mirrorConfig.WebhookId = ""   // Code needed to get WebhookID
		allRepos = append(allRepos, sourcerepo.Repo{Name: repo.Name, MirrorConfig: &mirrorConfig})
	}

	var nonMirroredRepos []sourcerepo.Repo

	// Iterates through all repos for any that aren't excluded from the list
	// or already mirrored in Google and adds them to the list
	for _, repo := range allRepos {
		alreadyMirrored := false
		excluded := contains(excludedRepos, repo.MirrorConfig.Url)
		if !excluded {
			for _, googleRepo := range googleRepos {
				if googleRepo.MirrorConfig == nil {
					// not a mirror repo
					continue
				}
				if googleRepo.MirrorConfig.Url == repo.MirrorConfig.Url {
					alreadyMirrored = true
				}
			}
			if !alreadyMirrored {
				nonMirroredRepos = append(nonMirroredRepos, repo)
			}
		}

	}
	return nonMirroredRepos, nil
}

// Creates google mirrors for the repositories in the nonMirroredRepos array.
func MirrorRepos(nonMirroredRepos []sourcerepo.Repo) string {
	mirroredRepos := []sourcerepo.Repo{}
	for _, repo := range nonMirroredRepos {
		mirror, err := CreateGoogleMirror(repo)
		if err != nil {
			panic(err)
		}
		mirroredRepos = append(mirroredRepos, *mirror)
	}
	return mirrorData(mirroredRepos)
}

// Checks to see if string s is in string array a
func contains(a []string, s string) bool {
	for _, i := range a {
		if strings.TrimSpace(i) == strings.TrimSpace(s) {
			return true
		}
	}
	return false
}
