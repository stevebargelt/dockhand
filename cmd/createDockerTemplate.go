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
	Short: "Creates a docker template in your Jenkins instance",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	RootCmd.AddCommand(createDockerTemplateCmd)
	createDockerTemplateCmd.Flags().StringVarP(&label, "label", "l", "", "The unique label to use for this Docker Template.")
	createDockerTemplateCmd.MarkFlagRequired("label")
	createDockerTemplateCmd.Flags().StringVarP(&image, "image", "i", "", "The docker image that this template will use for builds.")
	createDockerTemplateCmd.MarkFlagRequired("image")
	createDockerTemplateCmd.RunE = createDockerTemplate
}

func createDockerTemplate(cmd *cobra.Command, args []string) error {

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
