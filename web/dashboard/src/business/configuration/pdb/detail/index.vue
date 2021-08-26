<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'PDBs'}"  v-loading="loading">
    <yaml-editor :value="yaml" :read-only="true"></yaml-editor>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import {getPDB} from "@/api/poddisruptionbudgets"
export default {
  name: "PDBDetail",
  components: { YamlEditor, LayoutContent },
  props: {
    name: String,
    namespace: String,
  },
  data () {
    return {
      loading: false,
      cluster: "",
      item: {},
      yaml: {}
    }
  },
  methods: {
    getDetail() {
      this.loading = true
      getPDB(this.cluster, this.namespace, this.name).then(res => {
        this.item = res
        this.yaml = JSON.parse(JSON.stringify(this.item))
        this.loading = false
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.getDetail()
  }
}
</script>

<style scoped>

</style>
