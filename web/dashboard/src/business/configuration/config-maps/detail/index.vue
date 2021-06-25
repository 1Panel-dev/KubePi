<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'ConfigMaps'}" v-loading="loading">
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card>
          <table style="width: 100%" class="myTable">
            <tr>
              <th scope="col"  align="left">
                <h3>{{ $t("business.common.basic") }}</h3></th>
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
                  <el-tag type="info" size="small">
                    {{ key }} = {{ value.length > 100 ? value.substring(0, 100) + "..." : value }}
                  </el-tag>
                </div>
              </td>
            </tr>
          </table>
        </el-card>
      </el-col>
      <el-col :span="24">
        <br>
        <el-card>
          <div class="card_title">
            <h3>{{ $t("business.configuration.data") }}</h3>
          </div>
          <div>
            <json-editor :value="item.data">

            </json-editor>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getConfigMap} from "@/api/configmaps"
import JsonEditor from "@/components/json-editor"

export default {
  name: "ConfigMapDetail",
  components: { JsonEditor, LayoutContent },
  props: {
    name: String,
    namespace: String,
    yamlShow: {
      type: Boolean,
      default: function () {
        return false
      }
    },
  },
  data () {
    return {
      item: {
        metadata: {}
      },
      loading: false
    }
  },
  methods: {
    getDetail () {
      getConfigMap(this.$store.state.user.cluster, this.namespace, this.name).then(res => {
        this.item = res
      })
    }
  },
  created () {
    this.getDetail()
  }
}
</script>

<style scoped>

</style>