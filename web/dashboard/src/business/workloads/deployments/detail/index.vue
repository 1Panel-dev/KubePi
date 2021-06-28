<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Deployments'}" v-loading="loading">
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
              <td>{{ form.metadata.name }}</td>
            </tr>
            <tr>
              <td>{{ $t("business.namespace.namespace") }}</td>
              <td>{{ form.metadata.namespace }}</td>
            </tr>
            <tr>
              <td>{{ $t("commons.table.created_time") }}</td>
              <td>{{ form.metadata.creationTimestamp | datetimeFormat }}</td>
            </tr>
            <tr>
              <td>{{ $t("business.common.label") }}</td>
              <td colspan="4">
                <div v-for="(value,key,index) in form.metadata.labels" v-bind:key="index" class="myTag">
                  <el-tag type="info" size="small">
                    {{ key }} = {{ value }}
                  </el-tag>
                </div>
              </td>
            </tr>
            <tr>
              <td>{{ $t("business.common.annotation") }}</td>
              <td colspan="4">
                <div v-for="(value,key,index) in form.metadata.annotations" v-bind:key="index" class="myTag">
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
            <json-editor :value="form">
            </json-editor>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import { getDeploymentByName } from "@/api/workloads"
import JsonEditor from "@/components/json-editor"

export default {
  name: "DeploymentDetail",
  components: { JsonEditor, LayoutContent },
  data() {
    return {
      form: {
        metadata: {},
      },
      loading: false,
    }
  },
  methods: {
    getDetail() {
      this.loading = true
      getDeploymentByName(this.$route.params.cluster, this.$route.params.namespace, this.$route.params.name).then((res) => {
        this.form = res
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