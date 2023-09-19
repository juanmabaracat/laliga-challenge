package scraper

import (
	"bufio"
	"errors"
	"github.com/gocolly/colly/v2"
	"log/slog"
	"os"
	"strings"
)

var (
	errVisitingPage  = errors.New("error visiting the page")
	errEmptyFileName = errors.New("fileName cannot be empty")
)

const (
	targetURL = "https://www.laliga.com/en-GB/laliga-easports/standing"
)

type scraperClient interface {
	OnHTML(goquerySelector string, f colly.HTMLCallback)
	Visit(URL string) error
}

type Standing struct {
	Position string
	Team     string
	Points   string
}

type LaLiga interface {
	ScrapStandings(outputFileName string) error
}

func New() LaLiga {
	return &laLiga{scraperCli: colly.NewCollector()}
}

type laLiga struct {
	scraperCli scraperClient
}

// ScrapStandings scrap La Liga EA Sports standing page and save the results in outputFileName
func (l *laLiga) ScrapStandings(outputFileName string) error {
	standings := make([]Standing, 0)

	l.scraperCli.OnHTML("div.show div.styled__StandingTabBody-sc-e89col-0", func(e *colly.HTMLElement) {
		st := Standing{
			Position: e.ChildText("div:first-child > p:first-child"),
			Team:     e.ChildText("div:nth-child(2) > p"),
			Points:   e.ChildText("div:nth-child(3)"),
		}
		standings = append(standings, st)
	})

	slog.Info("scraping is about to start...")
	err := l.scraperCli.Visit(targetURL)
	if err != nil {
		return errVisitingPage
	}

	slog.Info("Standings scraped successfully", "total", len(standings))
	if saveErr := saveResults(standings, outputFileName); saveErr != nil {
		return saveErr
	}

	return nil
}

func saveResults(standings []Standing, fileName string) error {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		return errEmptyFileName
	}

	file, err := os.Create(fileName)
	if err != nil {
		slog.Error("Error creating file:", "error", err, "fileName", fileName)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	writer.WriteString("Position, Team, Points\n")

	for _, standing := range standings {
		record := []string{
			standing.Position,
			standing.Team,
			standing.Points,
		}

		_, writeErr := writer.WriteString(strings.Join(record, ", ") + "\n")
		if writeErr != nil {
			slog.Error("Error saving the results:", "error", writeErr)
			return writeErr
		}
	}

	slog.Info("La Liga EA Sports standings saved successfully", "file name", fileName)

	return nil
}
