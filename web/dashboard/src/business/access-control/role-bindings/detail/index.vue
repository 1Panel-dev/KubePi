<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'RoleBindings'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-row :gutter="20" class="row-box">
        <el-col :span="15">
          <el-card class="el-card">
            <ko-detail-basic :item="item" :yaml-show.sync="yamlShow"></ko-detail-basic>
          </el-card>
        </el-col>
        <el-col :span="9">
          <el-card>
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
      <yaml-editor :value="yaml" ref="yaml_editor" :read-only="true"></yaml-editor>
      <div class="bottom-button">
        <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.back_detail") }}</el-button>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import {getRoleBinding} from "@/api/rolebings"
import KoDetailRoleBinding from "@/components/detail/detail-role-binding"
import KoDetailBasic from "@/components/detail/detail-basic"

export default {
  name: "RoleBindingEdit",
  components: { YamlEditor, LayoutContent, KoDetailRoleBinding, KoDetailBasic },
  props: {
    name: String,
    namespace: String
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
    onCancel () {
      this.$router.push({ name: "RoleBindings" })
    },
    getDetail () {
      this.loading = true
      getRoleBinding(this.cluster, this.namespace, this.name).then(res => {
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
