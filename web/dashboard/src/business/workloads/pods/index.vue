<template>
  <div>
  <layout-content header="Pods">
    <el-alert
            v-if="showText"
            :title="$t('business.pod.metric_server_tip')"
            type="warning">
    </el-alert>
    <br>
    <complex-table :selects.sync="selects" :data="data" v-loading="loading" :pagination-config="paginationConfig"
                   :search-config="searchConfig" @search="search">
      <template #header>
        <el-button type="primary" size="small" @click="onCreate"
                   v-has-permissions="{scope:'namespace',apiGroup:'',resource:'pods',verb:'create'}">
          YAML
        </el-button>
        <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()"
                   v-has-permissions="{scope:'namespace',apiGroup:'',resource:'pods',verb:'delete'}">
          {{ $t("commons.button.delete") }}
        </el-button>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="80" show-overflow-tooltip fix>
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" min-width="45" prop="metadata.namespace"/>
      <el-table-column :label="$t('commons.table.status')" min-width="30">
        <template v-slot:default="{row}">
          {{ getPodStatus(row) }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.workload.restarts')" min-width="35">
        <template v-slot:default="{row}">
          {{ getRestartTimes(row) }}
        </template>
      </el-table-column>
      <el-table-column label="IP" min-width="40" prop="status.podIP"/>
      <el-table-column :label="$t('business.cluster.nodes')" min-width="45" show-overflow-tooltip prop="spec.nodeName"/>
      <el-table-column :label="'Cpu'"  min-width="45">
        <template v-slot:default="{row}">
          {{ getPodUsage(row.metadata.name, "cpu") }}
        </template>
      </el-table-column>
      <el-table-column :label="'Memory'" min-width="45">
        <template v-slot:default="{row}">
          {{ getPodUsage(row.metadata.name, "memory") }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" show-overflow-tooltip min-width="35" prop="metadata.creationTimestamp" fix>
        <template v-slot:default="{row}">
          {{ row.metadata.creationTimestamp | age }}
        </template>
      </el-table-column>
      <el-table-column min-width="45" :label="$t('commons.table.action')">
        <template v-slot:default="{row}">
          <el-button circle @click="onEdit(row)" size="mini" icon="el-icon-edit"
                     v-has-permissions="{scope:'namespace',apiGroup:'',resource:'pods',verb:'update'}"/>
          <el-dropdown style="margin-left: 10px" @command="handleClick($event,row)" :hide-on-click="false">
            <el-button circle icon="el-icon-more" size="mini"/>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item icon="el-icon-download" command="download">{{ $t("commons.button.download_yaml") }}
              </el-dropdown-item>
              <div v-if="row.containers.length > 1">
                <el-popover placement="left" trigger="hover">
                  <div v-for="c in row.containers" :key="c">
                    <p style="margin: 0">
                      <el-button @click="openTerminal(row, c)" type="text">{{ c }}</el-button>
                    </p>
                  </div>
                  <el-dropdown-item slot="reference" icon="el-icon-date" command="terminal">
                    {{ $t("commons.button.terminal") }}
                    <i class="el-icon-arrow-right"/>
                  </el-dropdown-item>
                </el-popover>
                <el-popover placement="left" trigger="hover">
                  <div v-for="c in row.containers" :key="c">
                    <p style="margin: 0">
                      <el-button @click="openTerminalLogs(row, c)" type="text">{{ c }}</el-button>
                    </p>
                  </div>
                  <el-dropdown-item slot="reference" icon="el-icon-notebook-2" command="logs">
                    {{ $t("commons.button.logs") }}
                    <i class="el-icon-arrow-right"/>
                  </el-dropdown-item>
                </el-popover>
              </div>
              <div v-if="row.containers.length == 1">
                <el-dropdown-item icon="iconfont iconline-terminalzhongduan" command="terminal">
                  {{ $t("commons.button.terminal") }}
                </el-dropdown-item>
                <el-dropdown-item icon="el-icon-tickets" command="logs">{{ $t("commons.button.logs") }}
                </el-dropdown-item>
              </div>
              <el-dropdown-item icon="el-icon-delete" :disabled="!onCheckPermissions()" command="delete">
                {{ $t("commons.button.delete") }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </complex-table>
  </layout-content>
  </div>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {listWorkLoads, deleteWorkLoad, getWorkLoadByName} from "@/api/workloads"
import {downloadYaml} from "@/utils/actions"
import ComplexTable from "@/components/complex-table"
import {checkPermissions} from "@/utils/permission"
import {listPodMetrics} from "@/api/apis"

export default {
  name: "Pods",
  components: { LayoutContent, ComplexTable },
  data () {
    return {
      loading: false,
      data: [],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: "",
      },
      selects: [],
      clusterName: "",
      podUsage: [],
      showText: false
    }
  },
  methods: {
    openDetail (row) {
      this.$router.push({
        name: "PodDetail",
        params: { namespace: row.metadata.namespace, name: row.metadata.name },
        query: { yamlShow: false }
      })
    },
    onCheckPermissions () {
      return checkPermissions({ scope: "namespace", apiGroup: "", resource: "pods", verb: "delete" })
    },
    handleClick (btn, row) {
      switch (btn) {
        case "download":
          downloadYaml(row.metadata.name + ".yml", getWorkLoadByName(this.clusterName, "pods", row.metadata.namespace, row.metadata.name))
          break
        case "terminal":
          this.openTerminal(row)
          break
        case "logs":
          this.openTerminalLogs(row)
          break
        case "delete":
          this.onDelete(row)
          break
      }
    },
    onEdit (row) {
      this.$router.push({ name: "PodEdit", params: { namespace: row.metadata.namespace, name: row.metadata.name } })
    },
    openTerminal (row, container) {
      let c
      if (container) {
        c = container
      } else {
        c = row.containers[0]
      }
      let routeUrl = this.$router.resolve({
        path: "/terminal",
        query: {
          cluster: this.clusterName,
          namespace: row.metadata.namespace,
          pod: row.metadata.name,
          container: c,
          type: "terminal"
        }
      })
      window.open(routeUrl.href, "_blank")
    },
    openTerminalLogs (row, container) {
      let c
      if (container) {
        c = container
      } else {
        c = row.containers[0]
      }
      let routeUrl = this.$router.resolve({
        path: "/terminal",
        query: {
          cluster: this.clusterName,
          namespace: row.metadata.namespace,
          pod: row.metadata.name,
          container: c,
          type: "log"
        }
      })
      window.open(routeUrl.href, "_blank")
    },
    onDelete (row) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deleteWorkLoad(this.clusterName, "pods", row.metadata.namespace, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteWorkLoad(this.clusterName, "pods", select.metadata.namespace, select.metadata.name))
            }
          }
        }
        if (this.ps.length !== 0) {
          Promise.all(this.ps)
            .then(() => {
              this.search(true)
              this.$message({
                type: "success",
                message: this.$t("commons.msg.delete_success"),
              })
            })
            .catch(() => {
              this.search(true)
            })
        }
      })
    },
    onCreate () {
      this.$router.push({ name: "PodCreate" })
    },
    getPodStatus (row) {
      if (row.status.containerStatuses) {
        let readyCount = 0
        for (const c of row.status.containerStatuses) {
          if (c.ready) {
            readyCount++
          }
        }
        return readyCount + "/" + row.status.containerStatuses.length
      }
      return
    },
    getRestartTimes (row) {
      if (row.status.containerStatuses) {
        let restartCount = 0
        for (const c of row.status.containerStatuses) {
          restartCount += c.restartCount
        }
        return restartCount
      }
      return 0
    },
    search (resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      listWorkLoads(this.clusterName, "pods", true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize)
        .then((res) => {
          console.log(res)
          this.data = res.items
          for (const item of this.data) {
            let container = []
            for (const c of item.spec.containers) {
              container.push(c.name)
            }
            item.containers = container
          }
          this.paginationConfig.total = res.total
          this.listPodMetric()
        })
        .catch(error => {
          console.log(error.message)
        })
        .finally(() => {
          this.loading = false
        })
    },
    listPodMetric () {
      const namespace = sessionStorage.getItem("namespace")
      listPodMetrics(this.clusterName, namespace).then(res => {
        this.podUsage = res.items
      }).catch(error => {
        this.showText = true
        console.log(error.message)
      })
    },
    getPodUsage (name, type) {
      let result = "0 m"
      if (this.podUsage.length > 0) {
        for (let item of this.podUsage) {
          if (item.metadata.name === name) {
            let usage = 0
            for (let container of item.containers) {
              if (type === "cpu") {
                if (container.usage.cpu.indexOf("n") > -1) {
                  usage = usage + parseInt(container.usage.cpu)
                }
                if (container.usage.cpu.indexOf("m") > -1) {
                  usage = usage + parseInt(container.usage.cpu) * 1000 * 1000
                }
              }
              if (type === "memory") {
                if (container.usage.memory.indexOf("Ki") > -1) {
                  usage = usage + parseInt(container.usage.memory)
                }
                if (container.usage.memory.indexOf("Mi") > -1) {
                  usage = usage + parseInt(container.usage.memory) * 1000
                }
              }
            }
            const unit = type === "cpu" ? "m" : "Mi"
            if (type === "cpu") {
              result = (usage / 1000000).toFixed(2)
            } else {
              result = (usage / 1000).toFixed(2)
            }
            result = result + unit
          }
        }
      }
      return result
    }
  },
  mounted () {
    this.clusterName = this.$route.query.cluster
    this.search()
  },
}
</script>
