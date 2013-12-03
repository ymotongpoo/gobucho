package bucho

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"launchpad.net/xmlpath"
)

const TorumemoURL = `http://oldriver.org/torumemo/`

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
// and returns status of the operation.
func Torumemo(browser bool, episode int) error {
	var episodePage string
	if episode > 0 {
		episodePage = fmt.Sprintf("%03d.html", episode)
	} else {
		episodePage = "index.html"
	}
	urlStr := TorumemoURL + episodePage

	var err error
	if browser {
		err = openUrl(urlStr)
	} else {
		resp, err := http.Get(urlStr)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		title, body, err := ParseTorumemo(resp.Body)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n\n%v", title, body)
	}
	return err
}

// openUrl opens a browser window to that location.
// This code taken from: http://stackoverflow.com/questions/10377243/how-can-i-launch-a-process-that-is-not-a-file-in-go
func openUrl(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", url, "http://localhost:4001/").Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Cannot open URL %s on this platform", url)
	}
	return err
}

func ParseTorumemo(content io.ReadCloser) (string, string, error) {
	datePath := xmlpath.MustCompile(`//div[@class="date"]`)
	titlePath := xmlpath.MustCompile(`//div[@class="title"]`)
	contentPath := xmlpath.MustCompile(`//div[@class="body"]/p`)
	root, err := xmlpath.ParseHTML(content)
	if err != nil {
		return "", "", err
	}

	date, _ := datePath.String(root)
	title, _ := titlePath.String(root)
	date = strings.TrimSpace(date)
	title = strings.TrimSpace(title)

	iter := contentPath.Iter(root)
	var body string
	for iter.Next() {
		body += iter.Node().String()
	}
	return date + " " + title, body, err
}
