def myCloud = Jenkins.instance.getInstance().getCloud("{{.Cloudname}}");

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

return uniqueLabels