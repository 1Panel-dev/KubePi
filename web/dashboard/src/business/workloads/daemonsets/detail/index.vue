<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'DaemonSets'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-card>
        <ko-detail-basic :item="form" :yaml-show.sync="yamlShow" />
      </el-card>
      <el-tabs style="margin-top:20px" v-model="activeName" type="border-card">
        <el-tab-pane label="Pods" name="Pods">
          <ko-detail-pods :cluster="clusterName" :namespace="namespace" :selector="selectors" />
        </el-tab-pane>
        <el-tab-pane label="Conditions" name="Conditions">
          <ko-detail-conditions :conditions="form.status.conditions" />
        </el-tab-pane>
      </el-tabs>
    </div>
    <div v-if="yamlShow">
      <yaml-editor :value="form" :read-only="true" />
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
  name: "DaemonSetDetail",
  components: { KoDetailConditions, KoDetailPods, KoDetailBasic, LayoutContent, YamlEditor },
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
        name: "DaemonSetDetail",
        params: { name: this.name, namespace: this.namespace },
        query: { yamlShow: newValue },
      })
    },
  },
  methods: {
    getDetail() {
      this.loading = true
      getWorkLoadByName(this.clusterName, "daemonsets", this.namespace, this.name).then((res) => {
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
