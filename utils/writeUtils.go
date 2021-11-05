package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"google.golang.org/api/sourcerepo/v1"
)

// Writes all retrieved Git repos to a file
func WriteGitReposToFile(FileName string) {
	var repoData string

	repoData += "Repo name\t\tLast Updated\n------------------------------------------------"

	for i := 1; true; i++ {
		data, err := RequestGitHubData(i)
		if err != nil {
			log.Fatal(err)
		}

		for _, d := range data {
			repoData += fmt.Sprintf("%s\t\t%s\n", d.Name, d.UpdatedAt)
		}

		if len(data) == 0 || i > 5 {
			break
		}
	}

	err := writeToFile(fmt.Sprintf("%s.txt", FileName), repoData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Printed repositories to \"%s.txt\"\n", FileName)
}

func WriteMirrorsToFile(FileName string, mirrors []sourcerepo.Repo) {
	var data string

	for i, r := range mirrors {
		data += "Repo " + strconv.Itoa(i) + ":\n"
		data += "\tName: " + r.Name + "\n"
		data += "\tURL: " + r.Url + "\n"
		data += "\tMirrorConfig:\n"
		data += "\t\tMirror_URL: " + r.MirrorConfig.Url + "\n"
		data += "\t\tMirror_WebHookID: " + r.MirrorConfig.WebhookId + "\n"
		data += "\t\tMirror_DeployKeyID: " + r.MirrorConfig.DeployKeyId + "\n"
		data += "\n"
	}
	err := writeToFile(FileName, data)
	if err != nil {
		panic(err)
	}
}

func writeToFile(fileName string, data string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}

	return file.Sync()
}
