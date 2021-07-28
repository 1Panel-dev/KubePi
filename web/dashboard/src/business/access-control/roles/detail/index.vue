<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Roles'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-card>
        <table style="width: 100%" class="myTable">
          <tr>
            <th scope="col" align="left">
              <h3>{{ $t("business.common.basic") }}</h3>
            </th>
            <th scope="col"></th>
          </tr>
          <tr>
            <td>{{ $t("commons.table.name") }}</td>
            <td>{{ item.metadata.name }}</td>
          </tr>
          <tr>
            <td>{{ $t("commons.table.namespace") }}</td>
            <td>{{ item.metadata.namespace }}</td>
          </tr>
          <tr>
            <td>{{ $t("commons.table.created_time") }}</td>
            <td>{{ item.metadata.creationTimestamp | age }}</td>
          </tr>
          <tr>
            <td>{{ $t("business.common.label") }}</td>
            <td colspan="4">
              <div v-for="(value,key,index) in item.metadata.labels" v-bind:key="index" class="myTag">
                <el-tag type="info" size="small">
                  {{ key }} = {{ value }}
                </el-tag>
              </div>
            </td>
          </tr>
          <tr>
            <td>{{ $t("business.common.annotation") }}</td>
            <td colspan="4">
              <div v-for="(value,key,index) in item.metadata.annotations" v-bind:key="index" class="myTag">
                <el-tag type="info" size="small" v-if="value.length < 100">
                  {{ key }} = {{ value }}
                </el-tag>
                <el-tooltip v-if="value.length > 100" :content="value" placement="top">
                  <el-tag type="info" size="small" v-if="value.length >= 100">
                    {{ key }} = {{ value.substring(0, 100) + "..." }}
                  </el-tag>
                </el-tooltip>
              </div>
            </td>
          </tr>
        </table>
        <div class="bottom-button">
          <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.view_yaml") }}</el-button>
        </div>
      </el-card>
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

export default {
  name: "RoleDetail",
  components: { LayoutContent, YamlEditor, KoDetailRoles },
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
