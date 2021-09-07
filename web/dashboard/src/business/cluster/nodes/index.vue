<template>
  <layout-content header="Nodes">
    <complex-table :data="data" v-loading="loading" :selects.sync="selects" :pagination-config="paginationConfig" @search="search()"
                   :search-config="searchConfig">
      <template #header>
        <el-button type="primary" size="small" :disabled="selects.length !== 1 ||  selects[0] === null || selects[0].spec.unschedulable"
                   @click="cordon(selects[0],true)" icon="el-icon-video-pause">
           Cordon
        </el-button>
        <el-button type="primary" size="small" :disabled="selects.length !== 1 ||  selects[0] === null || !selects[0].spec.unschedulable"
                   @click="cordon(selects[0],false)" icon="el-icon-video-play">
          Uncordon
        </el-button>
        <el-button type="primary" size="small" :disabled="selects.length !== 1"
                   @click="drain(selects[0])" icon="el-icon-refresh-right">
          Drain
        </el-button>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" fix max-width="30px">
        <template v-slot:default="{row}">
          <el-link @click="onDetail(row)"> {{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column label="Internal IP" prop="metadata.name" max-width="30px">
        <template v-slot:default="{row}">
          <div v-for="(address,index) in row.status.addresses" v-bind:key="index">
            <span v-if="address.type === 'InternalIP'">{{ address.address }}</span>
          </div>
        </template>
      </el-table-column>
<!--      <el-table-column :label="$t('business.node.ready')" prop="status.conditions" max-width="30px">-->
<!--        <template v-slot:default="{row}">-->
<!--          <div v-for="(condition,index) in row.status.conditions" v-bind:key="index">-->
<!--            <span v-if="condition.type === 'Ready'">{{ condition.status }}</span>-->
<!--          </div>-->
<!--        </template>-->
<!--      </el-table-column>-->
      <el-table-column :label="$t('commons.table.status')">
        <template v-slot:default="{row}">
          <el-button v-if="row.spec.unschedulable === true" type="warning" size="mini" plain round>
            Cordoned
          </el-button>
          <el-button v-else type="success" size="mini" plain round>
            Active
          </el-button>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.node.role')" prop="metadata.labels" min-width="180px" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <div v-for="(item,name,index) in row.metadata.labels" :key="index" style="display:inline-block">
            <span v-if="item"></span>
            <el-tag v-if="name.indexOf('node-role.kubernetes.io') > -1">{{ name.substring(name.indexOf(".io") + 4) }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <ko-table-operations :buttons="buttons" :label="$t('commons.table.action')"></ko-table-operations>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {cordonNode, listNodes} from "@/api/nodes"
import KoTableOperations from "@/components/ko-table-operations"
import {evictionPod, listPods, listPodsWithNsSelector} from "@/api/pods"
import {checkPermissions} from "@/utils/permission"
import {deleteNamespace} from "@/api/namespaces"

export default {
  name: "NodeList",
  components: { KoTableOperations, ComplexTable, LayoutContent },
  data () {
    return {
      data: [],
      loading: false,
      keywords: "",
      clusterName: "",
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: ""
      },
      selects: [],
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              path: "/nodes/edit/" + row.metadata.name,
              query: { yamlShow: false }
            })
          },
          disabled: () => {
            return !checkPermissions({ scope: "cluster", apiGroup: "", resource: "nodes", verb: "update" })
          }
        },
        {
          label: this.$t("commons.button.view_yaml"),
          icon: "el-icon-view",
          click: (row) => {
            this.$router.push({ path: "/nodes/detail/" + row.metadata.name, query: { yamlShow: "true" } })
          }
        },
      ]
    }
  },
  methods: {
    search () {
      this.loading = true
      const { currentPage, pageSize } = this.paginationConfig
      listNodes(this.clusterName, true, this.searchConfig.keywords, currentPage, pageSize).then(res => {
        this.loading = false
        this.data = res.items
        this.paginationConfig.total = res.total
      })
    },
    onDetail (row) {
      this.$router.push({ path: "/nodes/detail/" + row.metadata.name, query: { yamlShow: "false" } })
    },


    cordon(row,isCordon) {
      this.$confirm(
        this.$t("commons.confirm_message.delete"),
        this.$t("commons.message_box.prompt"), {
          confirmButtonText: this.$t("commons.button.confirm"),
          cancelButtonText: this.$t("commons.button.cancel"),
          type: "warning",
        }).then(() => {
        let data = { spec: { unschedulable: isCordon } }
        cordonNode(this.clusterName, row.metadata.name, data).then(() => {
          this.$message({ type: "success", message: this.$t("commons.msg.operation_success") })
          this.search()
        })
      })
    },
    drain(row) {
      this.$confirm(
        this.$t("commons.confirm_message.delete"),
        this.$t("commons.message_box.prompt"), {
          confirmButtonText: this.$t("commons.button.confirm"),
          cancelButtonText: this.$t("commons.button.cancel"),
          type: "warning",
        }).then(() => {
        let data = { spec: { unschedulable: true } }
        cordonNode(this.clusterName, row.metadata.name, data).then(() => {
          listPodsWithNsSelector(this.clusterName,"","","spec.nodeName="+row.metadata.name).then(res => {
            this.batchDeletePod(res.items, row.metadata.name)
          })
        })
      })
    },
    batchDeletePod (items, nodeName) {
      const ps = []
      for (const pod of items) {
        if (pod.spec.nodeName === nodeName && pod.metadata.ownerReferences.kind !== "daemonset") {
          const rmPod = {
            apiVersion: "policy/v1beta1",
            kind: "Eviction",
            metadata: {
              name: pod.metadata.name,
              namespace: pod.metadata.namespace,
              creationTimestamp: null,
            },
            deleteOptions: {},
          }
          console.log(pod.spec.nodeName)
          ps.push(evictionPod(this.clusterName, pod.metadata.namespace, pod.metadata.name, rmPod))
        }
      }
      Promise.all(ps)
        .then(() => {
          this.search()
          this.$message({
            type: "success",
            message: this.$t("business.node.drain_success"),
          })
        })
        .catch(() => {
          this.search()
        })
    }

  },
  created () {
    this.clusterName = this.$route.query.cluster
    this.search()
  }
}
</script>

<style scoped>

</style>
