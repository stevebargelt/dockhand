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
	"bytes"
	"fmt"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stevebargelt/dockhand/jenkins"
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

	data := struct {
		Cloudname string
	}{viper.GetString("cloudname")}

	t, err := template.ParseFiles("scripts/getLabels.groovy")
	var tpl bytes.Buffer
	err = t.Execute(&tpl, data)
	if err != nil {
		panic(err)
	}

	body, err := jenkins.RunScript(viper.GetString("jenkinsurl"), viper.GetString("username"), viper.GetString("password"), tpl.String())

	fmt.Println(body)

	return nil
}
