<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{ name: 'Namespaces' }" v-loading="loading">
    <el-row class="row-box">
      <el-card v-if="!yamlShow" class="el-card">
        <el-col :span="24">
          <ko-detail-basic v-if="item.metadata.name" :item="item" :yaml-show.sync="yamlShow"></ko-detail-basic>
        </el-col>
      </el-card>
    </el-row>
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
import {getNamespace} from "@/api/namespaces"
import YamlEditor from "@/components/yaml-editor"
import KoDetailBasic from "@/components/detail/detail-basic"

export default {
  name: "NamespaceDetail",
  components: { YamlEditor, LayoutContent, KoDetailBasic },
  props: {
    name: String
  },
  data () {
    return {
      item: {
        metadata: {},
        status: {}
      },
      loading: false,
      yamlShow: false,
      yaml: {}
    }
  },
  methods: {
    getNamespaceByName () {
      this.loading = true
      getNamespace(this.cluster, this.name).then(res => {
        this.loading = false
        this.item = res
        this.yaml = JSON.parse(JSON.stringify(this.item))
      })
    },
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        path: "/namespaces/detail/" + this.name,
        query: { yamlShow: newValue }
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.getNamespaceByName()
  },
}
</script>

<style scoped>

</style>
