package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aegershman/pivotal-tracker-resource/models"

	"github.com/salsita/go-pivotaltracker/v5/pivotal"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Expected path to sources as first arg")
	}
	sourceDir := os.Args[1]
	if err := os.Chdir(sourceDir); err != nil {
		log.Fatalf("Failed to access source dir '%s': %s", sourceDir, err)
	}

	var outRequest models.OutRequest
	if err := json.NewDecoder(os.Stdin).Decode(&outRequest); err != nil {
		log.Fatalln(err)
	}

	err := outRequest.Params.ReplaceName(sourceDir)
	if err != nil {
		log.Fatalln(err)
	}

	client := pivotal.NewClient(outRequest.Source.Token)
	if outRequest.Source.BaseURL != "" {
		if err := client.SetBaseURL(outRequest.Source.BaseURL); err != nil {
			log.Fatalln(err)
		}
	}

	// Filter works like this
	// name:"some story name goes here"
	filter := fmt.Sprintf("name:\"%s\"", outRequest.Params.Name)
	log.Println("Searching using filter: ", filter)
	stories, err := client.Stories.List(
		outRequest.Source.ProjectID,
		filter,
	)
	if err != nil {
		log.Fatalln(err)
	}

	// If story already exists, update it. Otherwise, create it.
	if len(stories) == 0 {
		log.Println("No matches. Creating new story.")

		createdStory, _, err := client.Stories.Create(
			outRequest.Source.ProjectID,
			&outRequest.Params.StoryRequest,
		)

		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Created", createdStory.Id, createdStory.Name)

	} else {
		matchingStory := stories[0]
		log.Printf("Updating %d %s", matchingStory.Id, matchingStory.Name)

		if len(stories) > 1 {
			log.Printf("But found %d similar matching stories:", len(stories))
			for _, s := range stories {
				fmt.Println(s.Id, s.Name)
			}
		}

		if _, _, err := client.Stories.Update(
			outRequest.Source.ProjectID,
			matchingStory.Id,
			&outRequest.Params.StoryRequest,
		); err != nil {
			log.Fatalln(err)
		}

	}

	fmt.Print("{\"version\":{\"ref\":\"none\"}}")
}
