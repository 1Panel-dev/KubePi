<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Ingresses'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-row :gutter="20" class="row-box">
        <el-col :span="24">
          <el-card class="el-card">
            <ko-detail-basic :item="item" :yaml-show.sync="yamlShow"></ko-detail-basic>
          </el-card>
        </el-col>
      </el-row>
      <el-row :gutter="20" class="row-box">
      <el-col :span="24">
          <el-card class="el-card" style="width:800px" v-for="c in certificates" v-bind:key="c.key">
             <div class="card_title">
              <h3>{{ $t("business.configuration.certificate_info") }}:{{c.secret}}:{{c.key}}</h3>
              
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
      </el-row>
      <el-row :gutter="20" class="row-box">
      <el-col :span="24" v-if="certificate_errs.length>0">
          <el-card class="el-card" style="width:800px" >
             <div class="card_title">
              <h3>{{ $t("business.configuration.certificate_errors") }}</h3>
              
            </div>
            <table style="width: 100%" class="myTable">
              <tr v-for="c in certificate_errs" v-bind:key="c.secret">
                 <td>Secret: {{ c.secret}}</td>
                 <td style='color:red'>{{ c.err }}</td>
              </tr>
            </table>
          </el-card>
        </el-col>
      </el-row>
      <el-row>
        <br>
        <el-col :span="24">
          <el-tabs type="border-card">
            <el-tab-pane :label="$t('business.network.rule')"  v-if="Object.keys(item.spec).length!==0">
              <ko-resource-rule :data="item.spec.rules"></ko-resource-rule>
            </el-tab-pane>
          </el-tabs>
        </el-col>
      </el-row>
    </div>
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
import {getIngress} from "@/api/ingress"
import YamlEditor from "@/components/yaml-editor"
import KoResourceRule from "@/components/detail/detail-rules"
import KoDetailBasic from "@/components/detail/detail-basic"
import {getSecret} from "@/api/secrets"
import * as x509 from "@peculiar/x509";
import { datetimeFormat } from "fit2cloud-ui/src/filters/time"

export default {
  name: "IngressDetail",
  components: { KoDetailBasic, KoResourceRule, LayoutContent, YamlEditor },
  props: {
    name: String,
    namespace: String
  },
  data () {
    return {
      item: {
        metadata: {},
        spec: {}
      },
      cluster: "",
      yamlShow: false,
      loading: false,
      yaml: {},
      /*证书信息*/
      certificates:[],
      certificate_errs:[],
    }
  },
  methods: {
    async getDetail () {
      this.loading = true
      const res= await getIngress(this.cluster, this.namespace, this.name)
      this.item = res
      this.yaml = JSON.parse(JSON.stringify(this.item))

      //显示tls证书信息
      await this.doWithCertificates()
      
      this.loading = false
     
    },
    //显示tls证书信息
    async doWithCertificates(){
            const tls = this.item.spec.tls
            if(!tls){
              return
            }
            if(tls && tls.length==0){
              this.certificate_errs.push({
                secret: "No Settings",
                err:"The tls attribute is not null but no secrets"
              })
              return ;
            }
            const { Base64 } = require("js-base64")
            for(let i=0,s=tls.length;i<s;i++){
              if(tls && tls.length>0){
                const secretName=tls[i].secretName
                
                if(secretName){
                  try{
                    const secret =  await getSecret(this.cluster, this.namespace, secretName)
                    /*分解证书信息*/
                    for (let key  in secret.data) {
                      const value =secret.data[key]
                      const base64_decode_value =Base64.decode(value)
                      if(base64_decode_value.startsWith("-----BEGIN CERTIFICATE-----")){
                        const cert = new x509.X509Certificate(base64_decode_value)
                        this.certificates.push({
                          secret: secretName,
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
                  }catch(e){
                    console.log(e)
                    this.certificate_errs.push({
                      secret: secretName,
                      err:"not exists"
                    })
                  }
                }
              }
            }
      }
    
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "IngressDetail",
        params: { name: this.name, namespace: this.namespace },
        query: { yamlShow: newValue }
      })
    }
  },
  async created () {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    await this.getDetail()
  }
}
</script>

<style scoped>

</style>
