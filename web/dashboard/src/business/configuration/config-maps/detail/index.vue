<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'ConfigMaps'}" v-loading="loading">
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
            <td>{{ $t("business.namespace.namespace") }}</td>
            <td>{{ item.metadata.namespace }}</td>
          </tr>
          <tr>
            <td>{{ $t("commons.table.created_time") }}</td>
            <td>{{ item.metadata.creationTimestamp | datetimeFormat }}</td>
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
                  {{ key }} = {{value}}
                </el-tag>
                <el-tooltip v-if="value.length > 100" :content="value" placement="top">
                  <el-tag type="info" size="small" v-if="value.length >= 100">
                    {{ key }} = {{value.substring(0, 100) + "..." }}
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
      <br>
      <el-card>
        <div class="card_title">
          <h3>{{ $t("business.configuration.data") }}</h3>
        </div>
        <div>
          <div v-if="item.data">
            <json-editor :value="item.data">
            </json-editor>
          </div>
          <div v-else-if="item.binaryData">
            <span> Binary Data: {{ bystesLength(item.binaryData.content) }} bytes</span>
          </div>
          <div v-else>
            <span>{{ $t("business.configuration.no_data") }}</span>
          </div>
        </div>
      </el-card>
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
import { getConfigMap } from "@/api/configmaps"
import JsonEditor from "@/components/json-editor"
import YamlEditor from "@/components/yaml-editor"

export default {
  name: "ConfigMapDetail",
  components: { YamlEditor, JsonEditor, LayoutContent },
  props: {
    name: String,
    namespace: String,
  },
  data() {
    return {
      item: {
        metadata: {},
      },
      loading: false,
      yamlShow: false,
      cluster: "",
      yaml: {},
    }
  },
  methods: {
    getDetail() {
      this.loading = true
      getConfigMap(this.cluster, this.namespace, this.name).then((res) => {
        this.item = res
        this.loading = false
        if (this.yamlShow) {
          this.yaml = JSON.parse(JSON.stringify(this.item))
        }
      })
    },
    bystesLength(str) {
      let count = 0
      for (let i = 0; i < str.length; i++) {
        if (str.charCodeAt(i) > 255) {
          count += 2
        } else {
          count++
        }
      }
      return count
    },
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "ConfigMapDetail",
        params: { name: this.name, namespace: this.namespace },
        query: { yamlShow: newValue }
      })
      this.getDetail()
    },
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.getDetail()
  },
}
</script>

<style scoped>
</style>