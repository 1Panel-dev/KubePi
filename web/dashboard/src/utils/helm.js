import FileSaver from "file-saver"
import { listSecretsWithNsAndLabelSelector } from "@/api/secrets"
import {
    createTarGzip
  } from "nanotar";
import JSONL from "@saiansh2525/jsonlines";

import pako from 'pako'
import {Base64} from "js-base64";
import yaml from "js-yaml";


function Uint8ArrayToString(fileData){
    var dataString = "";
    for (var i = 0; i < fileData.length; i++) {
      dataString += String.fromCharCode(fileData[i]);
    }
   
    return dataString
  
  }
// b64Data-->传入加密的数据进行解密
function gunzip (b64Data) {
    let strData = atob(b64Data)
    // Convert binary string to character-number array
    const charData = strData.split('').map(function (x) { return x.charCodeAt(0) })
    // Turn number array into byte-array
    const binData = new Uint8Array(charData)
    // // unzip
    const data = pako.inflate(binData)
    // Convert gunzipped byteArray back to ascii string:
    strData = Uint8ArrayToString(data)
    return strData
}
/*下载已安装的helm release 原始chart包*/
export async function downloadHelmReleases(clusterName, releases) {

    //用于生成tgz压缩文件
    let tgz=[];
    if(releases && releases.length>0){
     for (let i = 0; i < releases.length; i++) {
        const release = releases[i];
        const { name, namespace } = release;
        //获取对应secrets
        const secrets_res = await listSecretsWithNsAndLabelSelector(clusterName, namespace,"name="+name+",status=deployed,owner=helm");
        const secrets=secrets_res.items.filter(item => item.metadata.name.startsWith("sh.helm.release."))
        //加入文件到tgz
        await helmSecretToTgz(clusterName,namespace,secrets,tgz)
     }
     const tgzData=await createTarGzip(tgz)
     const blob = new Blob([tgzData], { type: "application/gzip" })
     FileSaver.saveAs(blob, `export.tar.gz`)
    }

}
//加入secret中的加密文件到tgz
async function helmSecretToTgz(clusterName,namespace,secrets,tgz) {
    for (let ii = 0; ii < secrets.length; ii++) {
        const secret=secrets[ii];
        //base64解密
        const gzip_encode_data=Base64.decode(secret.data.release)
        //console.log(gzip_encode_data)
        //gzip解密
        const decompressedData = gunzip(gzip_encode_data);
        //console.log ("Decompressed data:", decompressedData);

        const chartJSONS = JSONL.parse(decompressedData); // 根据换行或者回车进行识别

        for (let i = 0; i < chartJSONS.length; i++) {
            const chartJSON= chartJSONS[i];
            if(chartJSON["info"]["status"]=="deployed"){//已部署的chart release
            //加入加密后的chartjson数据到tgz
               await helmChartJSONToTgz(clusterName,namespace,chartJSON,tgz)
            }

        }
        
    }
}
//加入加密后的chartjson数据到tgz
async function helmChartJSONToTgz(clusterName,namespace,chartJSON,tgz) {

    //用户自定义values
    const user_provided_values_yaml = yaml.dump(chartJSON["config"]);

    //包名
    const chart_name = chartJSON["chart"]["metadata"]["name"];
    //版本
    const chart_version = chartJSON["chart"]["metadata"]["version"];

    
    //用户自定义values 文件
    tgz.push({name:clusterName+"/"+namespace+"/"+chart_name+"/"+chart_version+"/user_provided_values.yaml",data:user_provided_values_yaml});

    //处理Chart.yaml和values.yaml
    const chart_yaml = yaml.dump(chartJSON["chart"]["metadata"]);
    tgz.push({name:clusterName+"/"+namespace+"/"+chart_name+"/"+chart_version+"/"+chart_name+"/Chart.yaml",data:chart_yaml});

    //console.log(chart_yaml);
    const values_yaml = yaml.dump(chartJSON["chart"]["values"]);
    tgz.push({name:clusterName+"/"+namespace+"/"+chart_name+"/"+chart_version+"/"+chart_name+"/values.yaml",data:values_yaml});
    //console.log(values_yaml);


    for (let i = 0, s = chartJSON["chart"]["templates"].length; i < s; i++) {
        const template = chartJSON["chart"]["templates"][i];
        //模板文件路径
        const template_name = template["name"];
        //模板文件base64加密后
        const template_data = template["data"];
  
        const template_data_raw = Base64.decode(template_data);
  
        tgz.push({name:clusterName+"/"+namespace+"/"+chart_name+"/"+chart_version+"/"+chart_name+"/"+template_name,data:template_data_raw});
    
      }
}
