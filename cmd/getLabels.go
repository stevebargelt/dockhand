// Copyright © 2018 Steve Bargelt <steve@bargelt.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getLabelsCmd represents the getLabels command
var getLabelsCmd = &cobra.Command{
	Use:   "getLabels",
	Short: "Get the labels from the Docker Templates in Jenkins",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	RootCmd.AddCommand(getLabelsCmd)
	createDockerTemplateCmd.Flags().StringVarP(&cloudName, "cloudname", "c", "", "The Jenkins Yet Another Docker 'Cloud Name' to add Docker Template to")
	viper.BindPFlag("cloudname", createDockerTemplateCmd.PersistentFlags().Lookup("cloudname"))

	getLabelsCmd.RunE = getLabels
}

func getLabels(cmd *cobra.Command, args []string) error {

	crumbAPI := fmt.Sprintf("%s%s", viper.GetString("jenkinsurl"), "/crumbIssuer/api/json")

	req, err := http.NewRequest("GET", crumbAPI, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(viper.GetString("username"), viper.GetString("password"))

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

	script := `
	def myCloud = Jenkins.instance.getInstance().getCloud("` + viper.GetString("cloudname") + `");

if (!myCloud) {
  println("Cloud not found, aborting.") 
  return false
}

def templates = myCloud.getTemplates();

def uniqueLabels = []
templates.each { template ->
 words = template.labelString.split()
 def labelListForSlave = []
 words.each() {
          uniqueLabels.add(it)
 }
}
uniqueLabels.unique()

return uniqueLabels`

	query := fmt.Sprintf("%s&script=%s", crumbHeader, script)

	scriptURL := fmt.Sprintf("%s%s", viper.GetString("jenkinsurl"), "/scriptText")
	body := strings.NewReader(query)
	req, err = http.NewRequest("POST", scriptURL, body)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(viper.GetString("username"), viper.GetString("password"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}

	return nil
}
