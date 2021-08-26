<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'ResourceQuotas'}" v-loading="loading">
    <yaml-editor :value="yaml" :read-only="true"></yaml-editor>
    <div class="bottom-button">
      <el-button @click="onCancel()">{{ $t("commons.button.back_detail") }}</el-button>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import {getResourceQuota} from "@/api/resourcequota"

export default {
  name: "ResourceQuotaDetail",
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
    getDetail () {
      this.loading = true
      getResourceQuota(this.cluster, this.namespace, this.name).then(res => {
        this.item = res
        this.yaml = JSON.parse(JSON.stringify(this.item))
        this.loading = false
      })
    },
    onCancel () {
      this.$router.push({ name: "ResourceQuotas" })
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
