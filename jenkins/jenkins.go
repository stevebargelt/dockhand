package jenkins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type crumb struct {
	Class             string `json:"_class"`
	Crumb             string `json:"crumb"`
	CrumbRequestField string `json:"crumbRequestField"`
}

func GetCrumb(url, username, password string) string {

	crumbAPI := fmt.Sprintf("%s%s", url, "/crumbIssuer/api/json")

	req, err := http.NewRequest("GET", crumbAPI, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(username, password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var jenkinsCrumb crumb
	var crumbHeader string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(bodyBytes, &jenkinsCrumb)
		if err != nil {
			log.Fatal(err)
		}
		crumbHeader = fmt.Sprintf("%s=%s", jenkinsCrumb.CrumbRequestField, jenkinsCrumb.Crumb)
	}
	return crumbHeader

}

func RunScript(url string, username string, password string, script string) (string, error) {

	crumbHeader := GetCrumb(url, username, password)
	query := fmt.Sprintf("%s&script=%s", crumbHeader, script)

	scriptURL := fmt.Sprintf("%s%s", url, "/scriptText")
	body := strings.NewReader(query)
	req, err := http.NewRequest("POST", scriptURL, body)
	if err != nil {
		return "", nil
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var bodyString string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString = string(bodyBytes)
	}
	return bodyString, nil

}
