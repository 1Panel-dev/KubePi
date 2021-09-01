<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Services'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-row :gutter="20" class="row-box">
        <el-col :span="15">
          <el-card class="el-card" >
            <ko-detail-basic :item="item" :yaml-show.sync="yamlShow"></ko-detail-basic>
          </el-card>
        </el-col>
        <el-col :span="9">
          <el-card class="el-card" >
            <table style="width: 100%" class="myTable">
              <tr>
                <th scope="col" width="30%" align="left">
                  <h3>{{ $t("business.common.resource") }}</h3>
                </th>
                <th scope="col"></th>
              </tr>
              <tr>
                <td>{{ $t("business.configuration.type") }}</td>
                <td colspan="4">{{ item.spec.type }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.configuration.cluster_ip") }}</td>
                <td colspan="4">{{ item.spec.clusterIP }}</td>
              </tr>
              <tr>
                <td>Session Affinity</td>
                <td colspan="4">{{ item.spec.sessionAffinity }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.configuration.selector") }}</td>
                <td colspan="4">
                  <div v-for="(value,key,index) in item.spec.selector" v-bind:key="index" class="myTag">
                    <el-tag type="info" size="small">
                      {{ key }} = {{ value }}
                    </el-tag>
                  </div>
                </td>
              </tr>
            </table>
          </el-card>
        </el-col>
      </el-row>
      <el-row >
        <el-col :span="24">
          <br>
          <el-tabs type="border-card">
            <el-tab-pane label="Pods" v-if="Object.keys(item.metadata).length!==0">
              <resource-pod :cluster="cluster" :namespace="item.metadata.namespace"
                            :selector="getSelector(item.spec.selector)"></resource-pod>
            </el-tab-pane>
            <el-tab-pane label="Ports" v-if="Object.keys(item.spec).length!==0">
              <resource-ports :data="item.spec.ports"></resource-ports>
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
import YamlEditor from "@/components/yaml-editor"
import {getService} from "@/api/services"
import ResourcePod from "@/components/detail/detail-pods"
import ResourcePorts from "@/components/detail/detail-ports"
import KoDetailBasic from "@/components/detail/detail-basic"

export default {
  name: "ServiceDetail",
  components: { KoDetailBasic, ResourcePorts, ResourcePod, YamlEditor, LayoutContent },
  props: {
    name: String,
    namespace: String,
  },
  data () {
    return {
      item: {
        metadata: {},
        spec: {}
      },
      loading: false,
      yamlShow: false,
      cluster: "",
      yaml: {},
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getService(this.cluster, this.namespace, this.name).then((res) => {
        this.item = res
        this.yaml = JSON.parse(JSON.stringify(this.item))
      }).finally(() => {
        this.loading = false
      })
    },
    getSelector (selector) {
      let result = ""
      for (const key in selector) {
        if (Object.prototype.hasOwnProperty.call(selector, key)) {
          result = result + key + "=" + selector[key] + ","
        }
      }
      return result.substr(0, result.length - 1)
    },
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "ServiceDetail",
        params: { name: this.name, namespace: this.namespace },
        query: { yamlShow: newValue }
      })
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.getDetail()
  },
}
</script>

<style scoped>
</style>
