<template>
    <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'LimitRanges'}"  v-loading="loading">
      <yaml-editor :value="yaml" :read-only="true"></yaml-editor>
    </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getLimitRange} from "@/api/limitranges"
import YamlEditor from "@/components/yaml-editor"
export default {
  name: "LimitRangeDetail",
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
      getLimitRange(this.cluster, this.namespace, this.name).then(res => {
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
