<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'ConfigMaps'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-card>
        <ko-detail-basic :item="item" :yaml-show.sync="yamlShow"></ko-detail-basic>
      </el-card>
      <br>
      <el-card>
        <div class="card_title">
          <h3>{{ $t("business.configuration.data") }}</h3>
        </div>
        <div>
          <div v-if="item.data">
            <json-editor :value="item.data">
            </json-editor>
          </div>
          <div v-else-if="item.binaryData">
            <span> Binary Data: {{ bystesLength(item.binaryData.content) }} bytes</span>
          </div>
          <div v-else>
            <span>{{ $t("business.configuration.no_data") }}</span>
          </div>
        </div>
      </el-card>
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
import JsonEditor from "@/components/json-editor"
import YamlEditor from "@/components/yaml-editor"
import KoDetailBasic from "@/components/detail/detail-basic"

export default {
  name: "ConfigMapDetail",
  components: { KoDetailBasic, YamlEditor, JsonEditor, LayoutContent },
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
