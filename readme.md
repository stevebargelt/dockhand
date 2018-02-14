# dockhand

A simple tool to help configuring Jenkins [Yet Another Docker Plugin](https://plugins.jenkins.io/yet-another-docker-plugin) without using the Jenkins web UI.

If you are using YAD plugin you know that managing Docker Templates can get out of hand - so many settings in the UI, so easy to miss something important.

I created this tool for two reasons:

1. Configuring Docker templates in the UI is time consuming and error-prone.
2. I wanted to be able to automate the creation of new Cloud Templates as part of a automated build-engineering process. My scripts can use this tool to accomplish that.

For more information on my automated build system see:
[Automated Build System Tutorial Videos](https://www.youtube.com/playlist?list=PLJ3o2ZgH1Q-AEZxfcQ5S4Eat7ggaixuXj)

[Automated Build System Blog Posts](https://www.bargelt.com/2016/10/06/automated-build-system-docker-jenkins-azure-go-intro/)


## Get Labels

```shell
./dockhand getLabels --username jenkinsyaduser --password correcthorsebatterystaple --jenkinsurl https://localhost:8080 --cloudname JenkinsCloud
```

or, if you have a config file with your basic config items (username, password, jenkinsurl) it is as simple as:

```shell
./dockhand getLabels --cloudname JenkinsCloud
```

## Create Yet Another Docker Plugin Template

![dockhand in action](https://media.giphy.com/media/3o7WIFH959CSjU2AbS/giphy.gif)

```shell
./dockhand createDockerTemplate --username jenkinsyaduser --password correcthorsebatterystaple --jenkinsurl https://localhost:8080 --label dotnetcore_2 --image microsoft/dotnet:2.0-sdk
```

or, if you have a config file with your basic config items (username, password, jenkinsurl) it is as simple as:

```shell
./dockhand createDockerTemplate --cloudname EphemeralContainers --label dotnetcore_2 --image microsoft/dotnet:2.0-sdk
```

## It is early

No polish, no glitter. Alright it is downright ugly. 
Works for me at the moment.
