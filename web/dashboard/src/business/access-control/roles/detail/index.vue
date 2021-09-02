<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Roles'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-row class="row-box">
        <el-card class="el-card">
          <ko-detail-basic :item="item" :yaml-show.sync="yamlShow"></ko-detail-basic>
        </el-card>
      </el-row>
      <el-row>
        <el-col :span="24">
          <br>
          <el-tabs type="border-card">
            <el-tab-pane label="Rule" v-if="item.rules.length >0">
              <ko-detail-roles :rules="item.rules"></ko-detail-roles>
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
import {getRole} from "@/api/roles"
import KoDetailRoles from "@/components/detail/detail-roles"
import KoDetailBasic from "@/components/detail/detail-basic"

export default {
  name: "RoleDetail",
  components: { KoDetailBasic, LayoutContent, YamlEditor, KoDetailRoles },
  props: {
    name: String,
    namespace: String
  },
  data () {
    return {
      yaml: {},
      yamlShow: false,
      item: {
        metadata: {},
        rules: []
      },
      loading: false
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getRole(this.cluster,this.namespace, this.name).then(res => {
        this.item = res
        this.yaml = JSON.parse(JSON.stringify(this.item))
        this.loading = false
      })
    }
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "RoleDetail",
        params: {namespace:this.namespace, name: this.name },
        query: { yamlShow: newValue }
      })
    },
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
