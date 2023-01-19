package main

import (
	"fmt"
	"strings"

	"github.com/silinternational/git-mirror/utils"
)

const (
	// filename to save the mirrored repos to
	fileName string = "Mirrored_Repos"
)

func main() {
	// Get a slice of non mirrored repos
	nonMirroredRepos, err := utils.GetNonMirroredRepos()
	if err != nil {
		panic(err)
	}

	if len(nonMirroredRepos) == 0 {
		fmt.Println("found no repos to mirror")
		return
	}

	// Output the names of all the non mirrored repos
	fmt.Print("The following repositories are not mirrored: \n\n")
	for i, repo := range nonMirroredRepos {
		fmt.Printf("%d: %s\n", i, repo.MirrorConfig.Url)
	}

	// Request confirmation to confirm
	var answer string
	fmt.Print("Would you like to mirror these repositories on Google? (y/n):")
	_, err = fmt.Scanf("%s", &answer)
	if err != nil {
		panic(err)
	}

	switch strings.ToLower(answer) {
	case "y", "yes":
		data := utils.MirrorRepos(nonMirroredRepos)
		err = utils.WriteToFile(fileName, data)
		if err != nil {
			panic(err)
		}
		fmt.Println("The names of the new mirrors are recorded in the text file: " + fileName)
	case "n", "no":
		fmt.Println("No new mirrors were created.")
	default:
		fmt.Println("Command not understood, please answer with 'y' or 'n'.")
	}
}
