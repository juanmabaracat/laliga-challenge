package scraper

import (
	"errors"
	"github.com/gocolly/colly/v2"
	"os"
	"testing"
)

type mockCollyClient struct {
	*colly.Collector
	err error
}

func (m *mockCollyClient) OnHTML(goquerySelector string, f colly.HTMLCallback) {
}

func (m *mockCollyClient) Visit(URL string) error {
	return m.err
}

func TestNew(t *testing.T) {
	laLigaScraper := New()
	if laLigaScraper == nil {
		t.Error("error creating scraper")
	}
}

func Test_laLiga_ScrapStandings(t *testing.T) {
	tests := []struct {
		name           string
		scraperCli     scraperClient
		outputFileName string
		wantErr        bool
	}{
		{
			name:           "should scrap without error",
			scraperCli:     &mockCollyClient{},
			outputFileName: "test.txt",
			wantErr:        false,
		},
		{
			name:           "should return an error when there a problem connecting to the server",
			scraperCli:     &mockCollyClient{err: errors.New("connection error")},
			outputFileName: "test.txt",
			wantErr:        true,
		},
		{
			name:           "should return an error when it cannot save the file",
			scraperCli:     &mockCollyClient{},
			outputFileName: "",
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			laliga := &laLiga{
				scraperCli: tt.scraperCli,
			}
			if err := laliga.ScrapStandings(tt.outputFileName); (err != nil) != tt.wantErr {
				t.Errorf("ScrapStandings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_saveResults(t *testing.T) {
	tests := []struct {
		name      string
		standings []Standing
		fileName  string
		wantErr   error
	}{
		{
			name: "should save results without error",
			standings: []Standing{{
				Position: "1",
				Team:     "Real Madrid FC",
				Points:   "17",
			}},
			fileName: "test.txt",
			wantErr:  nil,
		},
		{
			name: "should return an error when the file name is empty",
			standings: []Standing{{
				Position: "1",
				Team:     "Real Madrid FC",
				Points:   "17",
			}},
			fileName: "",
			wantErr:  errEmptyFileName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := saveResults(tt.standings, tt.fileName)
			defer os.Remove(tt.fileName)
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("saveResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil && err != nil {
				t.Errorf("GOT ERROR = %v, but expected no error", err)
			}
		})
	}
}
