<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'ConfigMaps'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-row class="row-box">
        <el-card class="el-card">
          <ko-detail-basic :item="item" :yaml-show.sync="yamlShow"></ko-detail-basic>
        </el-card>
        <br>
        <el-card>
          <div class="card_title">
            <h3>{{ $t("business.configuration.data") }}</h3>
          </div>
          <div>
            <div v-if="item.data">
              <div v-for="(value,key) in item.data" v-bind:key="key">
                <ko-data :title="key">
                  <json-viewer v-if="jsonV(value)" :value="getContent(value)" :copyable=true
                               theme="jv-dark" :expanded="true" :expand-depth="4"></json-viewer>
                  <el-card v-else style="background: #202124;border: 0;">
                    <div style="white-space: pre-line;background: #202124;">
                      <span>{{ getValue(value) }} </span>
                    </div>
                  </el-card>
                </ko-data>
              </div>
            </div>
            <div v-else-if="item.binaryData">
              <span> Binary Data: {{ bystesLength(item.binaryData.content) }} bytes</span>
            </div>
            <div v-else>
              <span>{{ $t("business.configuration.no_data") }}</span>
            </div>
          </div>
        </el-card>
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
import {getConfigMap} from "@/api/configmaps"
import YamlEditor from "@/components/yaml-editor"
import KoDetailBasic from "@/components/detail/detail-basic"
import KoData from "@/components/ko-data"
import {isJSON} from "@/utils/data"

export default {
  name: "ConfigMapDetail",
  components: { KoDetailBasic, YamlEditor, LayoutContent, KoData },
  props: {
    name: String,
    namespace: String,
  },
  data () {
    return {
      item: {
        metadata: {},
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
      getConfigMap(this.cluster, this.namespace, this.name).then((res) => {
        this.item = res
        this.loading = false
        if (this.yamlShow) {
          this.yaml = JSON.parse(JSON.stringify(this.item))
        }
      })
    },
    bystesLength (str) {
      let count = 0
      for (let i = 0; i < str.length; i++) {
        if (str.charCodeAt(i) > 255) {
          count += 2
        } else {
          count++
        }
      }
      return count
    },
    getContent (value) {
      return JSON.parse(value)
    },
    jsonV (str) {
      return isJSON(str)
    },
    getValue (value) {
      if (value instanceof Object) {
        return JSON.parse(value)
      }
      return value
    }
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "ConfigMapDetail",
        params: { name: this.name, namespace: this.namespace },
        query: { yamlShow: newValue }
      })
      this.getDetail()
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
