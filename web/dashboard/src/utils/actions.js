import {download} from "fit2cloud-ui/src/tools/utils"
window.jsyaml = require("js-yaml")

export function downloadYaml (name, content) {
  let yaml = require("js-yaml")
  return download(name,yaml.dump(content))
}