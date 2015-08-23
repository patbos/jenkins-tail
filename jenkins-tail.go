package main

import (
	"flag"
	"strings"
	"errors"
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
	"strconv"
	"net/url"
	"time"
)

func getText(jobUrl string, start int) (string, int, bool, error) {

	values := make(url.Values)
	values.Set("start", strconv.Itoa(start))

	resp, err := http.PostForm(jobUrl, values)

	if err != nil {
		return "", 0, false, err
	}
	defer resp.Body.Close()
	if (resp.StatusCode != 200) {
		return "", 0, false, errors.New("Got " + strconv.Itoa(resp.StatusCode) + " response code from Jenkins!")
	}

	body, err := ioutil.ReadAll(resp.Body)
	moreData, err := strconv.ParseBool(resp.Header.Get("X-More-Data"))
	testSize, err := strconv.Atoi(resp.Header.Get("X-Text-Size"))
	return string(body), testSize, moreData, nil
}


func getJobUrl(jobName string) string {
	parts := strings.Split(jobName, "/")

	var result string

	for _, name := range parts {
		result = result + "/job/" + strings.Replace(name, " ", "%20", -1)

	}

	return result

}


func main() {

	baseUrl := flag.String("url", "", "Jenkins URL")
	job := flag.String("job", "", "Jenkins job name")

	jobNumber := flag.String("build", "lastBuild", "build number")

	flag.Parse()

	jobName := getJobUrl(*job)
	jenkinsUrl := *baseUrl + jobName + "/" + *jobNumber + "/logText/progressiveText"

	var start int
	var body string
	var moreData bool
	moreData = true
	start = 0
	for (moreData) {
		var newStart int
		var err error
		body, newStart, moreData, err = getText(jenkinsUrl, start)
		if (err != nil) {
			fmt.Println(err)
			os.Exit(3)
		}
		if (newStart > start) {
			fmt.Print(string(body))
		}
		start = newStart
		time.Sleep(200 * time.Millisecond)
	}
}
