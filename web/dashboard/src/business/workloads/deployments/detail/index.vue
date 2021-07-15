<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Deployments'}" v-loading="loading">
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
            <td>{{ form.metadata.name }}</td>
          </tr>
          <tr>
            <td>{{ $t("business.namespace.namespace") }}</td>
            <td>{{ form.metadata.namespace }}</td>
          </tr>
          <tr>
            <td>{{ $t("commons.table.created_time") }}</td>
            <td>{{ form.metadata.creationTimestamp | age }}</td>
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
        <div class="bottom-button">
          <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.view_yaml") }}</el-button>
        </div>
      </el-card>
      <el-tabs style="margin-top:20px" v-model="activeName">
        <el-tab-pane label="Pods" name="Pods">
          <complex-table :data="pods">
            <el-table-column sortable :label="$t('commons.table.status')" prop="status.phase" min-width="30">
              <template v-slot:default="{row}">
                <el-button v-if="row.status.phase ==='Running'" type="success" size="mini" plain round>
                  {{row.status.phase}}
                </el-button>
                <el-button v-if="row.status.phase ==='Terminating'" type="warning" size="mini" plain round>
                  {{row.status.phase}}
                </el-button>
              </template>
            </el-table-column>
            <el-table-column sortable :label="$t('commons.table.name')" prop="metadata.name" min-width="90" />
            <el-table-column sortable :label="$t('business.cluster.nodes')" prop="spec.nodeName" min-width="70" />
            <el-table-column sortable :label="$t('business.pod.image')" min-width="170">
              <template v-slot:default="{row}">
                <div v-for="(item,index) in row.spec.containers" v-bind:key="index" class="myTag">
                  <el-tag type="info" size="small">
                    {{ item.image }}
                  </el-tag>
                </div>
              </template>
            </el-table-column>
          </complex-table>
        </el-tab-pane>
        <el-tab-pane label="Conditions" name="Conditions">
          <complex-table :data="form.status.conditions">
            <el-table-column sortable label="Condition" prop="type" />
            <el-table-column sortable :label="$t('commons.table.status')" prop="status" />
            <el-table-column sortable :label="$t('commons.table.lastUpdateTime')" prop="lastUpdateTime">
              <template v-slot:default="{row}">
                {{ row.lastUpdateTime | age }}
              </template>
            </el-table-column>
            <el-table-column sortable :label="$t('commons.table.message')" min-width="200">
              <template v-slot:default="{row}">
                <span v-if="row.message">[{{ row.reason }} ]: {{ row.message }}</span>
                <span v-if="!row.message">---</span>
              </template>
            </el-table-column>
          </complex-table>
        </el-tab-pane>
      </el-tabs>
    </div>
    <div v-if="yamlShow">
      <yaml-editor :value="form" :read-only="true"></yaml-editor>
      <div class="bottom-button">
        <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.back_detail") }}</el-button>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import { getDeploymentByName } from "@/api/deployments"
import { listPodsWithNsSelector } from "@/api/pods"
import YamlEditor from "@/components/yaml-editor"

import ComplexTable from "@/components/complex-table"

export default {
  name: "DeploymentDetail",
  components: { LayoutContent, ComplexTable, YamlEditor },
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
      pods: [],
    }
  },
  methods: {
    getDetail() {
      this.loading = true
      getDeploymentByName(this.clusterName, this.$route.params.namespace, this.$route.params.name).then((res) => {
        this.form = res
        if (this.form.spec.selector.matchLabels) {
          let selectors = ""
          for (const key in this.form.spec.selector.matchLabels) {
            if (Object.prototype.hasOwnProperty.call(this.form.spec.selector.matchLabels, key)) {
              selectors += key + "=" + this.form.spec.selector.matchLabels[key] + ","
            }
          }
          selectors = selectors.length !== 0 ? selectors.substring(0, selectors.length - 1) : ""
          listPodsWithNsSelector(this.clusterName, this.$route.params.namespace, selectors).then((res) => {
            this.pods = res.items
          })
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
