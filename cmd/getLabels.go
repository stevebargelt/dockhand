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
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stevebargelt/dockhand/jenkins"
)

var getLabelsCmd = &cobra.Command{
	Use:   "getLabels",
	Short: "Get the labels from the YAD Docker Templates in Jenkins",
	Long: `Returns the labels used in all of the Docker Templates in 
				 all of the clouds of type Jenkins Yet Another Docker Plugin .`,
}

func init() {

	RootCmd.AddCommand(getLabelsCmd)
	getLabelsCmd.RunE = getLabels
}

// GetLabels returns the labels used in a YAD Cloud
func GetLabels() ([]string, error) {

	script, err := ioutil.ReadFile("scripts/getLabels.groovy")
	if err != nil {
		return nil, err
	}

	body, err := jenkins.RunScript(viper.GetString("jenkinsurl"), viper.GetString("username"), viper.GetString("password"), string(script))
	if err != nil {
		return nil, err
	}

	body = strings.Replace(body, "Result: [", "", -1)
	body = strings.Replace(body, "]", "", -1)
	body = strings.TrimSuffix(body, "\n")
	labels := strings.Split(body, ", ")

	return labels, nil

}

func getLabels(cmd *cobra.Command, args []string) error {

	labels, err := GetLabels()
	if err != nil {
		return err
	}

	for _, label := range labels {
		fmt.Println(label)
	}

	return nil
}
