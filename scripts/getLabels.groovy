def uniqueLabels = []

Jenkins.instance.clouds.each { cloud ->
  
  if (cloud.getClass() == com.github.kostyasha.yad.DockerCloud) {
    def templates = cloud.getTemplates();

	templates.each { template ->
          uniqueLabels.add(template.labelString	)
	}
	uniqueLabels.unique()
  }

}

return uniqueLabels