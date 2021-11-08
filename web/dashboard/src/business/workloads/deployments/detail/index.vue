<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Deployments'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-row :gutter="20" class="row-box">
        <el-col :span="12">
          <el-card class="el-card">
            <ko-detail-basic :item="form" :yaml-show.sync="yamlShow" />
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="el-card">
            <h3>{{ $t("business.common.conditions") }}</h3>
            <ko-detail-conditions :conditions="form.status.conditions" />
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px" class="row-box">
        <el-col :span="12">
          <el-card class="el-card">
            <ko-detail-service :cluster="clusterName" :namespace="namespace" :name="name" resourceType="Deployment" />
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="el-card">
            <ko-detail-ingress :cluster="clusterName" :namespace="namespace" :name="name" resourceType="Deployment" />
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin-top: 20px" class="row-box">
        <el-col :span="8">
          <el-card class="el-card">
            <ko-detail-pause @refresh="getDetail()" :cluster="clusterName" :yamlInfo="form" resourceType="deployments" />
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card class="el-card">
            <ko-detail-update @refresh="getDetail()" :cluster="clusterName" :yamlInfo="form" resourceType="deployments" />
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card class="el-card">
            <ko-detail-image @refresh="getDetail()" :cluster="clusterName" :yamlInfo="form" resourceType="deployments" />
          </el-card>
        </el-col>
      </el-row>
      <el-row style="margin-top:20px" class="row-box">
        <el-card class="el-card">
          <h3>Pods</h3>
          <ko-detail-pods :cluster="clusterName" :namespace="namespace" :selector="selectors" />
        </el-card>
      </el-row>
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
import KoDetailConditions from "@/components/detail/detail-conditions"
import KoDetailBasic from "@/components/detail/detail-basic"
import KoDetailPods from "@/components/detail/detail-pods"
import KoDetailService from "@/components/detail/detail-service"
import KoDetailIngress from "@/components/detail/detail-ingress"
import KoDetailPause from "@/components/detail/detail-pause"
import KoDetailUpdate from "@/components/detail/detail-update"
import KoDetailImage from "@/components/detail/detail-image"

export default {
  name: "DeploymentDetail",
  components: { KoDetailPods, KoDetailBasic, KoDetailConditions, KoDetailService, KoDetailIngress, KoDetailPause, KoDetailUpdate, KoDetailImage, LayoutContent, YamlEditor },
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
        name: "DeploymentDetail",
        params: { name: this.name, namespace: this.namespace },
        query: { yamlShow: newValue },
      })
    },
  },
  methods: {
    getDetail() {
      this.clusterName = this.$route.query.cluster
      this.loading = true
      getWorkLoadByName(this.clusterName, "deployments", this.namespace, this.name).then((res) => {
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
    this.getDetail()
  },
}
</script>

<style scoped>
</style>
