# Coding Challenge
ChatGPT doesn’t have access to current information. As an AI Developer, you need ChatGPT to answer the following questions -
1. Who is the current leader of La Liga EA Sports?
2. Which teams have more than 6 points?

## Technical Solution
Part A - Write a scrapper that can navigate to https://www.laliga.com/en-GB/laliga-easports/standing, uses some kind of “browser simulation” like Selenium to get all the textual information of the webpage. 
Once it obtains the text information, it can store the information in the filesystem as a flat file.
Part B - Write a script that uses the flat file and adds that information to the context window of ChatGPT using OpenAI API and then asks the questions that are mentioned above. 
Once you get the response, print it on the screen and save it as a file as well.

## Technical details
- [Go 1.21.0](https://go.dev/): Go version 1.21.0
- [Scraper](https://github.com/gocolly/colly/tree/master): Scraping Tool
- [go-openai](https://github.com/sashabaranov/go-openai): OpenAI API Client

## Requirements
- Go version 1.21.0
- OpenAI API key

### How to install GO
Mac using Homebrew:
```
brew install go
```
Other OS or installation options see: https://go.dev/doc/manage-install
### How to run the application:
Clone the repository and move to the project folder:
```
git clone https://github.com/juanmabaracat/laliga-challenge.git
cd laliga-challenge
```
Create the env variable OPENAI_API_KEY with your openai API key:
```
export OPENAI_API_KEY=[YOUR API KEY]
```
Run the application:
```
go run main.go
```
Run all tests (root folder):
```
go test ./...
```
## Notes
- The scraping result is saved in "scraping_results.txt".
- The conversation with the AI is printed in the console and also saved in conversation.txt.
- Every time the application is run, those files are overwritten.