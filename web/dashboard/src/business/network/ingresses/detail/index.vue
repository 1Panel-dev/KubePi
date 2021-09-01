<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Ingresses'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-row :gutter="20" class="row-box">
        <el-col :span="24">
          <el-card class="el-card">
            <ko-detail-basic :item="item" :yaml-show.sync="yamlShow"></ko-detail-basic>
          </el-card>
        </el-col>
      </el-row>
      <el-row>
        <br>
        <el-col :span="24">
          <el-tabs type="border-card">
            <el-tab-pane :label="$t('business.network.rule')"  v-if="Object.keys(item.spec).length!==0">
              <ko-resource-rule :data="item.spec.rules"></ko-resource-rule>
            </el-tab-pane>
          </el-tabs>
        </el-col>
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
import {getIngress} from "@/api/ingress"
import YamlEditor from "@/components/yaml-editor"
import KoResourceRule from "@/components/detail/detail-rules"
import KoDetailBasic from "@/components/detail/detail-basic"

export default {
  name: "IngressDetail",
  components: { KoDetailBasic, KoResourceRule, LayoutContent, YamlEditor },
  props: {
    name: String,
    namespace: String
  },
  data () {
    return {
      item: {
        metadata: {},
        spec: {}
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
      getIngress(this.cluster, this.namespace, this.name).then(res => {
        this.item = res
        this.yaml = JSON.parse(JSON.stringify(this.item))
      }).finally(() => {
        this.loading = false
      })
    },
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "IngressDetail",
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
