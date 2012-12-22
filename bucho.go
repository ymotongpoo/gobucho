package bucho

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func Show() string {
	// Say show :-)
	return text
}

func getStatuses() ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	r := new(Statuses)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}
	statuses := make([]string, len(*r))
	for i := range *r {
		statuses[i] = (*r)[i].Text
	}
	return statuses, nil
}

// LatestStatus returns bucho's latest tweet
func LatestStatus() (string, error) {
	statuses, err := getStatuses()
	if err != nil {
		return "", err
	}
	return statuses[0], nil
}

// AllStatus returns concatination of bucho's latest 30 tweets
func AllStatus() (string, error) {
	statuses, err := getStatuses()
	if err != nil {
		return "", err
	}
	return strings.Join(statuses, "\n"), nil
}

// Torumemo launches web browser with torumemo, one of the greatest text sites,
// and returns status of the operation. This function assumes default os as Mac or Linux
// TODO (ymotongpoo): test this function on windows
func Torumemo() (string, error) {
	var cmd *exec.Cmd
	var err error
	if os.Getenv("WINVER") != "" {
		cmd = exec.Command("$ie = New-Object -com InternetExplorer.Application;" + 
			`$ie.Navigate("http://oldriver.org/torumemo/")`)
	}
	cmd = exec.Command("open", "http://oldriver.org/torumemo/")
	
	if cmd == nil {
		return "wrong command", err
	} else if err = cmd.Run(); err != nil {
		return "Error", err
	}
	return "OK", nil
}

