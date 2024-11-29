<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{ name: 'Mutatingwebhookconfigurations' }"
    v-loading="loading">
    <div v-if="!yamlShow">
      <el-row class="row-box">
        <el-card class="el-card">
          <ko-detail-basic :item="item" :yaml-show.sync="yamlShow"></ko-detail-basic>
        </el-card>
        <br>

        <h2 v-if="yaml.webhooks && yaml.webhooks.length > 0">Webhooks</h2>

        <el-card class="el-card" v-for="(webhook, index) in (yaml.webhooks|| [])" v-bind:key="index">
        
          <WebHookDetail  :webhook="webhook"/>
          <br/>
        </el-card>
      </el-row>
    </div>
    <div v-if="yamlShow">
      <yaml-editor :value="yaml" :read-only="true"></yaml-editor>
      <div class="bottom-button">
        <el-button @click="yamlShow = !yamlShow">{{ $t("commons.button.back_detail") }}</el-button>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import { getMutatingwebhookconfiguration } from "@/api/mutatingwebhookconfiguration"
import YamlEditor from "@/components/yaml-editor"
import KoDetailBasic from "@/components/detail/detail-basic"
import { isJSON } from "@/utils/data"
import WebHookDetail from "@/components/webhook/webhook-detail.vue"

export default {
  name: "MutatingwebhookconfigurationDetail",
  components: { KoDetailBasic, YamlEditor, LayoutContent, WebHookDetail },
  props: {
    name: String
  },
  data() {
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
    getDetail() {
      this.loading = true
      getMutatingwebhookconfiguration(this.cluster, this.name).then((res) => {
        this.item = res
        this.loading = false
        this.yaml = JSON.parse(JSON.stringify(this.item))

      })
    },
    getContent(value) {
      return JSON.parse(value)
    },
    jsonV(str) {
      return isJSON(str)
    },
    getValue(value) {
      if (value instanceof Object) {
        return JSON.parse(value)
      }
      return value
    }
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "MutatingwebhookconfigurationDetail",
        params: { name: this.name },
        query: { yamlShow: newValue }
      })
      this.getDetail()
    },
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.getDetail()
  },
}
</script>

<style scoped></style>

