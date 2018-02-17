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
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stevebargelt/dockhand/jenkins"
)

var (
	label string
	image string
)

// createDockerTemplateCmd represents the createDockerTemplate command
var createDockerTemplateCmd = &cobra.Command{
	Use:   "createDockerTemplate",
	Short: "Creates a docker template in a Jenkins YAD cloud",
	Long: `Creates a 'Docker Template' in a Jenkins Yet Another Docker Plugin Cloud given 
	the cloud name, a label, and a docker image to pull. The label must be unique across
	all of the Docker Templates otherwise there is no guarantee a job will run on the 
	container that it was intended to run on.`,
}

func init() {
	RootCmd.AddCommand(createDockerTemplateCmd)
	createDockerTemplateCmd.Flags().StringVarP(&cloudName, "cloudname", "c", "", "The Jenkins Yet Another Docker 'Cloud Name' to add Docker Template to")
	createDockerTemplateCmd.MarkFlagRequired("cloudname")
	createDockerTemplateCmd.Flags().StringVarP(&label, "label", "l", "", "The unique label to use for this Docker Template")
	createDockerTemplateCmd.MarkFlagRequired("label")
	createDockerTemplateCmd.Flags().StringVarP(&image, "image", "i", "", "The docker image that this template will use")
	createDockerTemplateCmd.MarkFlagRequired("image")
	createDockerTemplateCmd.RunE = createDockerTemplate
}

func createDockerTemplate(cmd *cobra.Command, args []string) error {

	usedLabels, err := GetLabels()
	if err != nil {
		return nil
	}

	for _, usedLabel := range usedLabels {
		if label == usedLabel {
			return errors.New("cannot create a Docker Template with a duplicate label")
		}
	}

	data := struct {
		Cloudname string
		Label     string
		Image     string
	}{
		viper.GetString("cloudname"),
		label,
		image,
	}

	t, err := template.ParseFiles("scripts/createDockerTemplate.groovy")
	var tpl bytes.Buffer
	err = t.Execute(&tpl, data)
	if err != nil {
		return err
	}

	body, err := jenkins.RunScript(viper.GetString("jenkinsurl"),
		viper.GetString("username"), viper.GetString("password"), tpl.String())
	if err != nil {
		return err
	}

	body = strings.Replace(body, "Result: ", "", -1)
	fmt.Println(body)

	return nil
}
