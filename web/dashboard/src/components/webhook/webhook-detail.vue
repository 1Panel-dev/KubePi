<template>
    <table style="width: 100%" class="myTable">
            <tr>
              <td>{{ $t("commons.table.name") }}</td>
              <td >{{ webhook.name }}</td>
            </tr>
            <tr>
              <td>Client Config</td>
              <td >
                  <p>name:{{ webhook.clientConfig.service.name }}</p>
                  <p>namespace:{{ webhook.clientConfig.service.namespace }}</p>
                  <p>CA:</p>
                  <p v-for="(c,i) in this.getCA()" v-bind:key="i">{{c}}</p>
              </td>
            </tr>
            <tr>
              <td>Match Policy</td>
              <td>{{webhook.matchPolicy}}</td>
            </tr>
            <tr>
              <td>Failure Policy</td>
              <td>{{webhook.failurePolicy}}</td>
            </tr>
            <tr>
              <td>Admission Review Versions</td>
              <td>{{webhook.admissionReviewVersions.join(",")}}</td>
            </tr>
            <tr>
              <td>Reinvocation Policy</td>
              <td>{{webhook.reinvocationPolicy}}</td>
            </tr>
            <tr>
              <td>Side Effects</td>
              <td>{{webhook.sideEffects}}</td>
            </tr>
            <tr>
              <td>Timeout Seconds</td>
              <td>{{webhook.timeoutSeconds}}</td>
            </tr>
            <tr>
              <td>Namespace Selector</td>
              <td>
                <p v-for="(key, s1) in Object.keys(webhook.namespaceSelector)"  v-bind:key="s1">
                  {{key}}: <ko-detail-key-value v-if="webhook.namespaceSelector[key] && Object.keys(webhook.namespaceSelector[key]).length > 0" :valueObj="webhook.namespaceSelector[key]"></ko-detail-key-value>
                </p>
              </td>
            </tr>
            <tr>
              <td>Object Selector</td>
              <td>
                <p v-for="(key, s2) in Object.keys(webhook.objectSelector)"  v-bind:key="s2">
                  {{key}}: <ko-detail-key-value v-if="webhook.objectSelector[key] && Object.keys(webhook.objectSelector[key]).length > 0" :valueObj="webhook.objectSelector[key]"></ko-detail-key-value>
                </p>
              </td>
            </tr>
            <tr>
              <td>Rules</td>
              <td>
                <p v-for="(rule, s3) in webhook.rules"  v-bind:key="s3">
                  <ko-detail-key-value  :valueObj="rule"></ko-detail-key-value>
                </p>
              </td>
            </tr>
          </table>
</template>

<script>
import KoDetailKeyValue from "@/components/detail/detail-key-value"
import * as x509 from "@peculiar/x509";
import moment from "moment";

export default {
  name: "WebHookDetail",
  components: {  KoDetailKeyValue },
  props: {
    webhook: Object
  },
  data() {
    return {

    }
  },
  methods: {
      getCA(){
        const { Base64 } = require("js-base64")
        const certificates= [];
        const certificateString = Base64.decode( this.webhook["clientConfig"]["caBundle"]);
        if (!certificateString.startsWith("-----BEGIN CERTIFICATE-----")) {
            certificates.push ("Invalid certificate")
        }
        const cert = new x509.X509Certificate(certificateString)
        certificates.push("Subject:"+cert.subject)
        certificates.push("SignatureAlgorithm:"+cert.signatureAlgorithm.name+" "+cert.signatureAlgorithm.hash.name)
        certificates.push("Issuer:"+cert.issuer)
        certificates.push("Not before:"+moment(cert.notBefore).format())
        certificates.push("Expires:"+moment(cert.notAfter).format())
        return certificates;
      }
  }
}

</script>
<style scoped></style>
