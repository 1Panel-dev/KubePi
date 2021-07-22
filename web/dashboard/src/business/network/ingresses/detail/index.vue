<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Ingresses'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-row :gutter="20">
        <el-col :span="24">
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
                <td>{{ $t("business.namespace.namespace") }}</td>
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
import KoResourceRule from "@/components/resources/resource-rules"

export default {
  name: "IngressDetail",
  components: { KoResourceRule, LayoutContent, YamlEditor },
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
