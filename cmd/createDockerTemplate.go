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

var (
	jenkinsUsername string
	jenkinsPassword string
	cloudName       string
	label           string
	image           string
)

type crumb struct {
	Class             string `json:"_class"`
	Crumb             string `json:"crumb"`
	CrumbRequestField string `json:"crumbRequestField"`
}

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

	createDockerTemplateCmd.Flags().StringVarP(&jenkinsUsername, "username", "u", "", "Jenkins username")
	createDockerTemplateCmd.MarkFlagRequired("username")
	createDockerTemplateCmd.Flags().StringVarP(&jenkinsPassword, "password", "p", "", "Jenkins password")
	createDockerTemplateCmd.MarkFlagRequired("password")
	createDockerTemplateCmd.Flags().StringVarP(&cloudName, "cloudname", "c", "", "The Jenkins Yet Another Docker 'Cloud Name' to add Docker Template to")
	createDockerTemplateCmd.Flags().StringVarP(&label, "label", "l", "", "The unique label to use for this Docker Template")
	createDockerTemplateCmd.Flags().StringVarP(&image, "image", "i", "", "the docker image for this template")

	viper.BindPFlag("username", createDockerTemplateCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", createDockerTemplateCmd.PersistentFlags().Lookup("password"))

	createDockerTemplateCmd.RunE = createDockerTemplate
}

func createDockerTemplate(cmd *cobra.Command, args []string) error {

	crumbAPI := fmt.Sprintf("%s%s", jenkinsURL, "/crumbIssuer/api/json")

	req, err := http.NewRequest("GET", crumbAPI, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(jenkinsUsername, jenkinsPassword)

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
import com.github.kostyasha.yad.commons.*;
import com.github.kostyasha.yad.DockerCloud;
import com.github.kostyasha.yad.DockerContainerLifecycle;
import com.github.kostyasha.yad.DockerSlaveTemplate;
import com.github.kostyasha.yad.launcher.DockerComputerJNLPLauncher;
import com.github.kostyasha.yad.strategy.DockerOnceRetentionStrategy;

// Let's find the cloud!
def myCloud = Jenkins.instance.getInstance().getCloud("` + cloudName + `");
if (!myCloud) {
  println("Cloud not found, aborting.") 
  return false
}

def label = "` + label + `"
def image = "` + image + `"

def launcher = new DockerComputerJNLPLauncher();
launcher.setUser("jenkins");
launcher.setLaunchTimeout(60);

def pullImage = new DockerPullImage();
pullImage.setPullStrategy(DockerImagePullStrategy.PULL_NEVER);

//remove
def removeContainer = new DockerRemoveContainer();
removeContainer.setRemoveVolumes(true);
removeContainer.setForce(true);

def createContainer = new DockerCreateContainer();

//allows Slaves to reference the host Docker to run Docker in Docker
//Inception. Nuff said.
def volumeList = ["/var/run/docker.sock:/var/run/docker.sock"]
createContainer.setVolumes(volumeList);

//lifecycle
def containerLifecycle = new DockerContainerLifecycle();
containerLifecycle.setImage(image);
containerLifecycle.setPullImage(pullImage);
containerLifecycle.setRemoveContainer(removeContainer);
containerLifecycle.setCreateContainer(createContainer);

//Node Properties (environment variables)
def nodeProperties = new ArrayList<>();

def slaveTemplate = new DockerSlaveTemplate();
slaveTemplate.setLabelString(label);
slaveTemplate.setLauncher(launcher);
slaveTemplate.setMode(Node.Mode.EXCLUSIVE);
slaveTemplate.setRetentionStrategy(new DockerOnceRetentionStrategy(5));
slaveTemplate.setDockerContainerLifecycle(containerLifecycle);
slaveTemplate.setNodeProperties(nodeProperties);

def templates = myCloud.getTemplates();
def newTemplates = new ArrayList<DockerSlaveTemplate>();
newTemplates.addAll(templates);
newTemplates.add(slaveTemplate);

myCloud.setTemplates(newTemplates);
Jenkins.getActiveInstance().save();

return true`

	query := fmt.Sprintf("%s&script=%s", crumbHeader, script)

	scriptURL := fmt.Sprintf("%s%s", jenkinsURL, "/scriptText")
	body := strings.NewReader(query)
	req, err = http.NewRequest("POST", scriptURL, body)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(jenkinsUsername, jenkinsPassword)
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
