<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'ClusterRoleBindings'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-row :gutter="20" class="row-box">
        <el-col :span="15">
          <el-card class="el-card">
            <ko-detail-basic :item="item" :yaml-show.sync="yamlShow"></ko-detail-basic>
          </el-card>
        </el-col>
        <el-col :span="9">
          <el-card class="el-card">
            <table style="width: 100%" class="myTable">
              <tr>
                <th scope="col" width="30%" align="left">
                  <h3>{{ $t("business.common.resource") }}</h3>
                </th>
                <th scope="col"></th>
              </tr>
              <tr>
                <td>Kind</td>
                <td>{{ item.roleRef.kind }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.access_control.role") }}</td>
                <td>{{ item.roleRef.name }}</td>
              </tr>
            </table>
          </el-card>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="24">
          <br>
          <el-tabs type="border-card" v-if="Object.keys(item.metadata).length!==0">
            <el-tab-pane :label="$t('business.access_control.authorized')">
              <ko-detail-role-binding :subjects="item.subjects"></ko-detail-role-binding>
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
import YamlEditor from "@/components/yaml-editor"
import {getClusterRoleBinding} from "@/api/clusterrolebindings"
import KoDetailBasic from "@/components/detail/detail-basic"
import KoDetailRoleBinding from "@/components/detail/detail-role-binding"

export default {
  name: "ClusterRoleBindingDetail",
  components: { KoDetailRoleBinding, YamlEditor, LayoutContent, KoDetailBasic },
  props: {
    name: String,
  },
  data () {
    return {
      loading: false,
      cluster: "",
      item: {
        metadata: {},
        subjects: [],
        roleRef: {}
      },
      yaml: {},
      yamlShow: false
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getClusterRoleBinding(this.cluster, this.name).then(res => {
        this.item = res
        this.yaml = JSON.parse(JSON.stringify(this.item))
        this.loading = false
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
