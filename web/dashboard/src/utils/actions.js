import {download} from "fit2cloud-ui/src/tools/utils"
window.jsyaml = require("js-yaml")

export function downloadYaml (name, getDetail) {
  let yaml = require("js-yaml")
  getDetail.then(res => {
    return download(name,yaml.dump(res))
  })
}
