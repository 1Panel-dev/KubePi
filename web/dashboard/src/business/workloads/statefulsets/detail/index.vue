<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'StatefulSets'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-row class="row-box">
        <el-card class="el-card">
          <ko-detail-basic :item="form" :yaml-show.sync="yamlShow"></ko-detail-basic>
        </el-card>
      </el-row>
      <el-row>
        <el-tabs style="margin-top:20px" v-model="activeName" type="border-card">
          <el-tab-pane label="Pods" name="Pods">
            <ko-detail-pods :cluster="clusterName" :namespace="namespace" :selector="selectors"></ko-detail-pods>
          </el-tab-pane>
          <el-tab-pane :label="$t('business.common.conditions')" name="Conditions">
            <ko-detail-conditions :conditions="form.status.conditions"></ko-detail-conditions>
          </el-tab-pane>
        </el-tabs>
      </el-row>
    </div>
    <div v-if="yamlShow">
      <yaml-editor :value="form" :read-only="true"></yaml-editor>
      <div class="bottom-button">
        <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.back_detail") }}</el-button>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import { getWorkLoadByName } from "@/api/workloads"
import YamlEditor from "@/components/yaml-editor"
import KoDetailBasic from "@/components/detail/detail-basic"
import KoDetailPods from "@/components/detail/detail-pods"
import KoDetailConditions from "@/components/detail/detail-conditions"

export default {
  name: "StatefulSetDetail",
  components: { KoDetailBasic, KoDetailPods, KoDetailConditions, LayoutContent, YamlEditor },
  props: {
    name: String,
    namespace: String,
  },
  data() {
    return {
      form: {
        metadata: {},
        status: {
          conditions: [],
        },
      },
      yamlShow: false,
      activeName: "Pods",
      loading: false,
      clusterName: "",
      selectors: "",
    }
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "StatefulSetDetail",
        params: { name: this.name, namespace: this.namespace },
        query: { yamlShow: newValue },
      })
    },
  },
  methods: {
    getDetail() {
      this.loading = true
      getWorkLoadByName(this.clusterName, "statefulsets", this.namespace, this.name).then((res) => {
        this.form = res
        if (this.form.spec.selector.matchLabels) {
          this.selectors = ""
          for (const key in this.form.spec.selector.matchLabels) {
            if (Object.prototype.hasOwnProperty.call(this.form.spec.selector.matchLabels, key)) {
              this.selectors += key + "=" + this.form.spec.selector.matchLabels[key] + ","
            }
          }
          this.selectors = this.selectors.length !== 0 ? this.selectors.substring(0, this.selectors.length - 1) : ""
        }
        this.loading = false
      })
    },
  },
  created() {
    this.clusterName = this.$route.query.cluster
    this.getDetail()
  },
}
</script>

<style scoped>
</style>
