package main

import (
	"fmt"
	"laliga-challenge/artificialintelligence"
	"laliga-challenge/scraper"
	"log"
	"os"
)

const (
	conversationFile    = "conversation.txt"
	scrapingResultsFile = "scraping_results.txt"
)

func main() {
	superAI, cliErr := artificialintelligence.NewClient()
	if cliErr != nil {
		log.Fatal("Error creating AI client: ", cliErr)
	}

	laLigaScraper := scraper.New()
	scrapErr := laLigaScraper.ScrapStandings(scrapingResultsFile)
	if scrapErr != nil {
		log.Fatal("error scraping la liga", scrapErr)
	}

	conversation, errConversationFile := os.Create(conversationFile)
	if errConversationFile != nil {
		log.Fatal(errConversationFile)
	}
	defer conversation.Close()

	standingData, openErr := os.ReadFile(scrapingResultsFile)
	if openErr != nil {
		log.Fatal("error opening scraping file", openErr)
	}

	superAI.AddMessage(string(standingData))

	messages := []string{
		"Given the information I gave you about the current situation of La Liga EA Sports, I'll ask you some questions.",
		"Who is the current leader of La Liga EA Sports?",
		"Which teams have more than 6 points?",
	}

	for _, message := range messages {
		fmt.Println(message)
		fmt.Fprintln(conversation, message)
		response, err := superAI.Ask(message)
		if err != nil {
			log.Fatal(err)
		}

		response = "AI RESPONSE: " + response
		fmt.Println(response)
		fmt.Fprintln(conversation, response)
	}

	fmt.Println("\nEnd of execution, goodbye!")
}
