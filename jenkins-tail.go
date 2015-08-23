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

func getLog(jobUrl string, start int) (string, int, bool, error) {

	values := make(url.Values)
	values.Set("start", strconv.Itoa(start))

	var moreData bool
	moreData = false

	resp, err := http.PostForm(jobUrl, values)

	if err != nil {
		return "", 0, false, err
	}
	defer resp.Body.Close()
	if (resp.StatusCode != 200) {
		return "", 0, false, errors.New("Got " + strconv.Itoa(resp.StatusCode) + " response code from Jenkins!")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, false, err
	}
	moreDataHeader := resp.Header.Get("X-More-Data")
	if (moreDataHeader != "") {
		moreData, err = strconv.ParseBool(moreDataHeader)
		if err != nil {
			return "", 0, false, err
		}
	}

	testSize, err := strconv.Atoi(resp.Header.Get("X-Text-Size"))
	if err != nil {
		return "", 0, false, err
	}
	return string(body), testSize, moreData, nil
}


func getJobUrl(baseUrl string, jobName string, jobNumber string) (string, error) {
	parts := strings.Split(jobName, "/")

	var result string

	for _, name := range parts {
		result = result + "/job/" + name

	}
	url, err := url.Parse(baseUrl + result + "/" + jobNumber + "/logText/progressiveText")
	if err != nil {
		return "", err
	}

	return url.String(), nil;

}


func main() {

	baseUrl := flag.String("url", "", "Jenkins URL")
	job := flag.String("job", "", "Jenkins job name")

	jobNumber := flag.String("build", "lastBuild", "build number")

	flag.Parse()

	jenkinsUrl, err := getJobUrl(*baseUrl, *job, *jobNumber)

	if (err != nil) {
		fmt.Println("Error: " + err.Error())
		os.Exit(3)
	}


	var start int
	var body string
	var moreData bool
	moreData = true
	start = 0
	for (moreData) {
		var newStart int
		var err error
		body, newStart, moreData, err = getLog(jenkinsUrl, start)
		if (err != nil) {
			fmt.Println("Error: " + err.Error())
			os.Exit(3)
		}
		if (newStart > start) {
			fmt.Print(body)
		}
		start = newStart
		time.Sleep(200 * time.Millisecond)
	}
}
