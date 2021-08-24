<template>
  <layout-content header="Pods">
    <complex-table :selects.sync="selects" :data="data" v-loading="loading" :pagination-config="paginationConfig" :search-config="searchConfig" @search="search">
      <template #header>
        <el-button-group>
          <el-button type="primary" size="small" @click="onCreate" v-has-permissions="{apiGroup:'core',resource:'pods',verb:'create'}">
            {{ $t("commons.button.create") }}
          </el-button>
          <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()" v-has-permissions="{apiGroup:'core',resource:'pods',verb:'delete'}">
            {{ $t("commons.button.delete") }}
          </el-button>
        </el-button-group>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column sortable :label="$t('commons.table.name')" prop="name" min-width="80" show-overflow-tooltip fix>
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column sortable :label="$t('business.namespace.namespace')" min-width="45" prop="metadata.namespace" />
      <el-table-column sortable :label="$t('commons.table.status')" min-width="30">
        <template v-slot:default="{row}">
          {{getPodStatus(row)}}
        </template>
      </el-table-column>
      <el-table-column sortable :label="$t('business.workload.restarts')" min-width="35">
        <template v-slot:default="{row}">
          {{getRestartTimes(row)}}
        </template>
      </el-table-column>
      <el-table-column sortable label="IP" min-width="40" prop="status.podIP" />
      <el-table-column sortable :label="$t('business.cluster.nodes')" min-width="45" prop="spec.nodeName" />
      <el-table-column :label="$t('commons.table.created_time')" min-width="35" prop="metadata.creationTimestamp" fix>
        <template v-slot:default="{row}">
          {{ row.metadata.creationTimestamp | age }}
        </template>
      </el-table-column>
      <el-table-column min-width="45" :label="$t('commons.table.action')">
        <template v-slot:default="{row}">
          <el-button circle @click="onEdit(row)" size="mini" icon="el-icon-edit" v-has-permissions="{apiGroup:'core',resource:'pods',verb:'update'}" />
          <el-dropdown style="margin-left: 10px" @command="handleClick($event,row)" :hide-on-click="false">
            <el-button circle icon="el-icon-more" size="mini" />
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item icon="el-icon-download" command="download">{{$t("commons.button.download_yaml")}}</el-dropdown-item>
              <div v-if="row.containers.length > 1">
                <el-popover placement="left" trigger="hover">
                  <div v-for="c in row.containers" :key="c">
                    <p style="margin: 0">
                      <el-button @click="openTerminal(row, c)" type="text">{{c}}</el-button>
                    </p>
                  </div>
                  <el-dropdown-item slot="reference" icon="el-icon-date" command="terminal">
                    {{$t("commons.button.terminal")}}
                    <i class="el-icon-arrow-right" />
                  </el-dropdown-item>
                </el-popover>
                <el-popover placement="left" trigger="hover">
                  <div v-for="c in row.containers" :key="c">
                    <p style="margin: 0">
                      <el-button @click="openTerminalLogs(row, c)" type="text">{{c}}</el-button>
                    </p>
                  </div>
                  <el-dropdown-item slot="reference" icon="el-icon-notebook-2" command="logs">
                    {{$t("commons.button.logs")}}
                    <i class="el-icon-arrow-right" />
                  </el-dropdown-item>
                </el-popover>
              </div>
              <div v-if="row.containers.length == 1">
                <el-dropdown-item icon="el-icon-date" command="terminal">{{$t("commons.button.terminal")}}</el-dropdown-item>
                <el-dropdown-item icon="el-icon-notebook-2" command="logs">{{$t("commons.button.logs")}}</el-dropdown-item>
              </div>
              <el-dropdown-item icon="el-icon-delete" :disabled="!onCheckPermissions()" command="delete">{{$t("commons.button.delete")}}</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import { listWorkLoads, deleteWorkLoad } from "@/api/workloads"
import { downloadYaml } from "@/utils/actions"
import ComplexTable from "@/components/complex-table"
import { checkPermissions } from "@/utils/permission"
import { randomNum } from "@/utils/randomNum"

export default {
  name: "Pods",
  components: { LayoutContent, ComplexTable },
  data() {
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
    }
  },
  methods: {
    openDetail(row) {
      this.$router.push({ name: "PodDetail", params: { namespace: row.metadata.namespace, name: row.metadata.name }, query: { yamlShow: false } })
    },
    onCheckPermissions() {
      checkPermissions({ apiGroup: "core", resource: "jobs", verb: "delete" })
    },
    handleClick(btn, row) {
      switch (btn) {
        case "download":
          downloadYaml(row.metadata.name + ".yml", row)
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
    onEdit(row) {
      this.$router.push({name: "PodEdit", params: {namespace: row.metadata.namespace, name: row.metadata.name}})
    },
    openTerminal(row, container) {
      let c
      if (container) {
        c = container
      } else {
        c = row.containers[0]
      }
      let existTerminals = this.$store.getters.terminals
      const item = {
        type: "terminal",
        key: randomNum(8),
        name: row.metadata.name,
        cluster: this.clusterName,
        namespace: row.metadata.namespace,
        pod: row.metadata.name,
        container: c,
        containers: row.containers,
      }
      existTerminals.push(item)
      this.$store.commit("terminal/TERMINALS", existTerminals)
    },
    openTerminalLogs(row, container) {
      let c
      if (container) {
        c = container
      } else {
        c = row.containers[0]
      }
      let existTerminals = this.$store.getters.terminals
      const item = {
        type: "logs",
        key: randomNum(8),
        name: row.metadata.name,
        cluster: this.clusterName,
        namespace: row.metadata.namespace,
        pod: row.metadata.name,
        container: c,
        containers: row.containers,
      }
      existTerminals.push(item)
      this.$store.commit("terminal/TERMINALS", existTerminals)
    },
    onDelete(row) {
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
    onCreate() {
      this.$router.push({ name: "PodCreate" })
    },
    getPodStatus(row) {
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
    getRestartTimes(row) {
      if (row.status.containerStatuses) {
        let restartCount = 0
        for (const c of row.status.containerStatuses) {
          restartCount += c.restartCount
        }
        return restartCount
      }
      return 0
    },
    search(resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      listWorkLoads(this.clusterName, "pods", true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize)
        .then((res) => {
          this.data = res.items
          for (const item of this.data) {
            let container = []
            for (const c of item.spec.containers) {
              container.push(c.name)
            }
            item.containers = container
          }
          this.paginationConfig.total = res.total
        })
        .catch((error) => {
          console.log(error)
        })
        .finally(() => {
          this.loading = false
        })
    },
  },
  mounted() {
    this.clusterName = this.$route.query.cluster
    this.search()
  },
}
</script>
