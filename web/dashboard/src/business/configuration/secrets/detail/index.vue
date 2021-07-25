<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Secrets'}" v-loading="loading">
    <el-row :gutter="20">
      <div v-if="!yamlShow">
        <el-col :span="24">
          <el-card>
            <table style="width: 100%" class="myTable">
              <tr>
                <th scope="col" align="left">
                  <h3>{{ $t("business.common.basic") }}</h3>
                </th>
                <th scope="col"></th>
              </tr>
              <tr>
                <td>{{ $t("commons.table.name") }}</td>
                <td>{{ item.metadata.name }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.configuration.type") }}</td>
                <td>{{ item.type }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.namespace.namespace") }}</td>
                <td>{{ item.metadata.namespace }}</td>
              </tr>
              <tr>
                <td>{{ $t("commons.table.created_time") }}</td>
                <td>{{ item.metadata.creationTimestamp | age }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.common.label") }}</td>
                <td colspan="4">
                  <div v-for="(value,key,index) in item.metadata.labels" v-bind:key="index" class="myTag">
                    <el-tag type="info" size="small">
                      {{ key }} = {{ value }}
                    </el-tag>
                  </div>
                </td>
              </tr>
              <tr>
                <td>{{ $t("business.common.annotation") }}</td>
                <td colspan="4">
                  <div v-for="(value,key,index) in item.metadata.annotations" v-bind:key="index" class="myTag">
                    <el-tag type="info" size="small" v-if="value.length < 100">
                      {{ key }} = {{ value }}
                    </el-tag>
                    <el-tooltip v-if="value.length > 100" :content="value" placement="top">
                      <el-tag type="info" size="small" v-if="value.length >= 100">
                        {{ key }} = {{ value.substring(0, 100) + "..." }}
                      </el-tag>
                    </el-tooltip>
                  </div>
                </td>
              </tr>
            </table>
            <div class="bottom-button">
              <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.view_yaml") }}</el-button>
            </div>
          </el-card>
        </el-col>
        <el-col :span="24">
          <br>
          <el-card>
            <div class="card_title">
              <h3>{{ $t("business.configuration.data") }}</h3>
            </div>
            <div v-for="(value,key) in item.data" v-bind:key="key">
              <ko-data :title="key">
                <json-viewer v-if="jsonV(value)" :value="getContent(value)" :copyable=true
                             theme="jv-dark" :expanded="true" :expand-depth="3"></json-viewer>
                <el-card v-else style="background: #112234;border: 0;">
                  <div style="white-space: pre-line;">
                    <span>{{ getValue(value) }} </span>
                  </div>
                </el-card>
              </ko-data>
            </div>
          </el-card>
        </el-col>
      </div>
      <div v-if="yamlShow">
        <yaml-editor :value="yaml" :read-only="true"></yaml-editor>
        <div class="bottom-button">
          <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.back_detail") }}</el-button>
        </div>
      </div>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getSecret} from "@/api/secrets"
import {isJSON} from "@/utils/data"
import KoData from "@/components/ko-data"
import YamlEditor from "@/components/yaml-editor"

export default {
  name: "SecretDetail",
  components: { YamlEditor, KoData, LayoutContent },
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
      yaml: {}
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getSecret(this.cluster, this.namespace, this.name).then(res => {
        this.loading = false
        this.item = res
        this.yaml = JSON.parse(JSON.stringify(this.item))
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
