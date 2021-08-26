<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'PSPs'}"  v-loading="loading">
    <yaml-editor :value="yaml" :read-only="true"></yaml-editor>
    <div class="bottom-button">
      <el-button @click="onCancel()">{{ $t("commons.button.back_detail") }}</el-button>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import {getPSP} from "@/api/podsecuritypolicies"
export default {
  name: "PSPDetail",
  components: { YamlEditor, LayoutContent },
  props: {
    name: String,
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
      getPSP(this.cluster, this.name).then(res => {
        this.item = res
        this.yaml = JSON.parse(JSON.stringify(this.item))
        this.loading = false
      })
    },
    onCancel() {
      this.$router.push({ name: "PSPs" })
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
