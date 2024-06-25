<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Secrets'}" v-loading="loading">
    <el-row :gutter="20" class="row-box">
      <div v-if="!yamlShow">
        <el-col :span="24">
          <el-card class="el-card">
            <ko-detail-basic :item="item" :yaml-show.sync="yamlShow"></ko-detail-basic>
          </el-card>
        </el-col>
        <el-col :span="24">
          <el-card class="el-card" style="width:800px" v-for="c in certificates" v-bind:key="c.key">
             <div class="card_title">
              <h3>{{ $t("business.configuration.certificate_info") }}:{{c.key}}</h3>
              
            </div>
            <table style="width: 100%" class="myTable">
              <tr>
                 <td>{{ $t("business.configuration.certificate_signatureAlgorithm") }}</td>
                 <td>{{ c.signatureAlgorithm }}</td>
              </tr>
              <tr>
                 <td>subject</td>
                 <td>{{ c.subject }}</td>
              </tr>
              <tr>
                 <td>issuer</td>
                 <td>{{ c.issuer }}</td>
              </tr>
              <tr>
                 <td>{{ $t("business.configuration.certificate_validity") }}</td>
                 <td>{{ c.validity }}<font style='color:red' v-if="c.notAfter< new Date()">({{ $t("business.configuration.certificate_outdate") }})</font></td>
              </tr>
            </table>
          </el-card>
        </el-col>
        <el-col :span="24">
          <br>
          <el-card class="el-card" style="width:800px">
            <div class="card_title">
              <h3>{{ $t("business.configuration.data") }}</h3>
            </div>
                        <div v-for="(value,key) in item.data" v-bind:key="key">
              <ko-data :title="key">
                <json-viewer v-if="jsonV(value)&&key!=='token'" :value="getContent(value)" :copyable="copyable"
                             theme="jv-dark" :expanded="true" :expand-depth="3"></json-viewer>
                <el-card v-else style="background: #202124;border: 0;">
                  <div style="white-space: pre-line;background: #202124;width: 100%">
                    <span style="word-break:break-all;">{{ getValue(value) }} </span>
                  </div>
                </el-card>
              </ko-data>
            </div>
          </el-card>
        </el-col>
      </div>
    </el-row>
    <div v-if="yamlShow">
      <yaml-editor :value="yaml" :read-only="true"></yaml-editor>
      <div class="bottom-button">
        <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.back_detail") }}</el-button>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getSecret} from "@/api/secrets"
import {isJSON} from "@/utils/data"
import KoData from "@/components/ko-data"
import YamlEditor from "@/components/yaml-editor"
import KoDetailBasic from "@/components/detail/detail-basic"
import * as x509 from "@peculiar/x509";
import { datetimeFormat } from "fit2cloud-ui/src/filters/time"

export default {
  name: "SecretDetail",
  components: { KoDetailBasic, YamlEditor, KoData, LayoutContent },
  props: {
    name: String,
    namespace: String
  },
  data () {
    return {
      item: {
        metadata: {},
        type: ""
      },
      cluster: "",
      yamlShow: false,
      loading: false,
      yaml: {},
      copyable: { copyText: this.$t("commons.actions.copy") },
      /*证书信息*/
      certificates:[]
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getSecret(this.cluster, this.namespace, this.name).then(res => {
        this.loading = false
        this.item = res
        this.yaml = JSON.parse(JSON.stringify(this.item))

        /*分解证书信息*/
        for (let key  in this.item.data) {
          const value =this.item.data[key]
          const base64_decode_value =this.getValue(value)
          if(base64_decode_value.startsWith("-----BEGIN CERTIFICATE-----")){
            const cert = new x509.X509Certificate(base64_decode_value)
            this.certificates.push({
              key:key,
              signatureAlgorithm:cert.signatureAlgorithm.name +" "+cert.signatureAlgorithm.hash.name,
              subject:cert.subject,
              issuer:cert.issuer,
              validity:datetimeFormat(cert.notBefore) + " ~ " + datetimeFormat(cert.notAfter),
              notBefore:cert.notBefore,
              notAfter:cert.notAfter,
            })
          }
        }
        //console.log(this.certificates)
      })
      
    },
    getContent (value) {
      const { Base64 } = require("js-base64")
      const content = Base64.decode(value)
      return JSON.parse(content)
    },
    jsonV (str) {
      const { Base64 } = require("js-base64")
      const content = Base64.decode(str)
      return isJSON(content)
    },
    getValue (value) {
      const { Base64 } = require("js-base64")
      return Base64.decode(value)
    }
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "SecretDetail",
        params: { name: this.name, namespace: this.namespace },
        query: { yamlShow: newValue }
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.getDetail()
  }
}
</script>

<style scoped>

</style>
