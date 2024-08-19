<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'events'}" v-loading="loading">
   
    <el-row>
      <div >
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
import { isJSON } from "@/utils/data"
import { getEventsWithNs } from "@/api/events"
import YamlEditor from "@/components/yaml-editor"

export default {
  name: "EventDetail",
  components: { YamlEditor, LayoutContent },
  props: {
    name: String,
    namespace: String,
    cluster: String
  },
  data() {
    return {
      item: {
        metadata: {},
        spec: {
        },
        status: {},
      },
      cluster: "",
      yamlShow: false,
      loading: false,
      yaml: {},
    }
  },
  methods: {
    getDetail() {
      this.loading = true
      getEventsWithNs(this.cluster, this.namespace,this.name).then((res) => {
        this.loading = false
        this.item = res
        this.yaml = JSON.parse(JSON.stringify(this.item))
      })
    },
    getContent(value) {
      const { Base64 } = require("js-base64")
      const content = Base64.decode(value)
      return JSON.parse(content)
    },
    jsonV(str) {
      const { Base64 } = require("js-base64")
      const content = Base64.decode(str)
      return isJSON(content)
    },
    getValue(value) {
      const { Base64 } = require("js-base64")
      return Base64.decode(value)
    },
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "EventDetail",
        params: { name: this.name,namespace: this.namespace, cluster: this.cluster},
        query: { yamlShow: newValue },
      })
      this.getDetail()
    },
  },
  created() {
    this.getDetail()
  },
}
</script>

<style scoped>
</style>
