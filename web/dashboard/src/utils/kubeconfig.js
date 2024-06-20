import FileSaver from "file-saver"
import { getCluster } from "@/api/clusters"
import { getSecret } from "@/api/secrets"
import { getServiceAccount } from "@/api/serviceaccounts"
import { $error } from "../plugins/message"
/*下载kubeconfig文件*/
export async function downloadKubeconfigFileByServiceAccount(clusterName, name, namespace) {
    let res = await getCluster(clusterName);
    let apiUrl = (res.data.spec.connect.forward.apiServer);

    let serviceAccount = await getServiceAccount(clusterName, namespace, name)

    let secretNames = serviceAccount.secrets;

    if (secretNames.length > 0) {
        let secretName = secretNames[0]["name"];
        let secret = await getSecret(clusterName, namespace, secretName)
        let contextName = clusterName + "-" + namespace + "-" + name
        const { token, "ca.crt": caCrt } = secret.data;
        let yaml = require("js-yaml")
        const { Base64 } = require("js-base64")
        let response = yaml.dump({
            apiVersion: "v1",
            kind: "Config",
            clusters: [
                {
                    name: contextName,
                    cluster: {
                        server: apiUrl,
                        "certificate-authority-data": caCrt,
                    },
                },
            ],
            users: [
                {
                    name: name,
                    user: {
                        token: Base64.decode(token),
                    },
                },
            ],
            contexts: [
                {
                    name: `${contextName}-${name}`,
                    context: {
                        user: name,
                        cluster: contextName,
                        namespace: namespace,
                    },
                },
            ],
            "current-context": `${contextName}-${name}`,
        })
        const blob = new Blob([response], { type: "text/yaml" })
        FileSaver.saveAs(blob, `${contextName}-${name}`)
    } else {
        $error("serviceaccount has no secret")
    }



}