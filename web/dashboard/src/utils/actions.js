//import {download} from "fit2cloud-ui/src/tools/utils"
window.jsyaml = require("js-yaml")
import FileSaver from "file-saver"
export function downloadYaml (name, getDetail) {
  let yaml = require("js-yaml")
  getDetail.then(res => {
    //return download(name,yaml.dump(res))
    const blob = new Blob([yaml.dump(res)], { type: "text/yaml" })
    FileSaver.saveAs(blob, name)
  })
}
