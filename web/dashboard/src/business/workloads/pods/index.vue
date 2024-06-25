<template>
  <div>
  <layout-content header="Pods">
    <div style="float: left">
      <el-button type="primary" size="small" @click="onCreate"
                   v-has-permissions="{scope:'namespace',apiGroup:'',resource:'pods',verb:'create'}">
        YAML
      </el-button>
      <el-button type="primary" size="small" @click="onTop"
                  v-has-permissions="{scope:'namespace',apiGroup:'',resource:'pods',verb:'list'}">
        Top
      </el-button>
      <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()"
                  v-has-permissions="{scope:'namespace',apiGroup:'',resource:'pods',verb:'delete'}">
        {{ $t("commons.button.delete") }}
      </el-button>
      <el-button type="primary" size="small"
                   @click="exportToXlsx()" icon="el-icon-download">
          export to excel
      </el-button>
    </div>
    <complex-table :selects.sync="selects" :data="data" v-loading="loading" :pagination-config="paginationConfig"
                   :search-config="searchConfig" @search="search" @sort-change='sortTableFun'>
      <el-table-column type="selection" fix ></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="80" show-overflow-tooltip fix sortable="name">
        <template v-slot:default="{row}">
          <span class="span-link" @click="openDetail(row)">{{ row.metadata.name }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" min-width="45" prop="namespace" sortable="namespace"/>
      <el-table-column :label="$t('commons.table.status')" min-width="30" sortable="status_phase" prop="status_phase">
        <template v-slot:default="{row}">
          <div v-if="row.status.phase ==='Running'">
            <i class="el-icon-check" />&nbsp; &nbsp; &nbsp;
            {{ $t("commons.status.Running") }}
          </div>
          <div v-if="row.status.phase ==='Failed'">
            <i class="el-icon-close" />&nbsp; &nbsp; &nbsp;
            {{ $t("commons.status.Failed") }}
          </div>
          <div v-if="row.status.phase ==='Pending'">
            <i class="el-icon-loading" />&nbsp; &nbsp; &nbsp;
            {{ $t("commons.status.Pending") }}
          </div>
          <div v-if="row.status.phase ==='Succeeded'">
            <i class="el-icon-finished" />&nbsp; &nbsp; &nbsp;
            {{ $t("commons.status.Succeeded") }}
          </div>
          <div v-if="row.status.phase ==='Unknown'">
            <i class="iconfont iconwenhao" />&nbsp; &nbsp; &nbsp;
            {{ $t("commons.status.Unknown") }}
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.ready')" min-width="30" prop="pod_status" sortable="pod_status">ready
        <template v-slot:default="{row}">
          {{ getPodStatus(row) }}
        </template>
      </el-table-column>
      <el-table-column label="IP" min-width="40" prop="podIP" sortable="podIP"/>
      <el-table-column :label="$t('business.cluster.nodes')" min-width="45" show-overflow-tooltip prop="nodeName" sortable="nodeName"/>
      <el-table-column :label="$t('commons.table.created_time')" show-overflow-tooltip min-width="35" prop="creationTimestamp" fix sortable="creationTimestamp">
        <template v-slot:default="{row}">
          {{ row.metadata.creationTimestamp | age }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.pod.restart_count')" show-overflow-tooltip min-width="35"  fix sortable="restart_count" prop="restart_count">
        <template v-slot:default="{row}">
          {{ getRestartTimes(row) }}
        </template>
      </el-table-column>
      <el-table-column width="90px" :label="$t('commons.table.action')">
        <template v-slot:default="{row}">
          <el-button circle @click="onEdit(row)" size="mini" icon="el-icon-edit"
                     v-has-permissions="{scope:'namespace',apiGroup:'',resource:'pods',verb:'update'}"/>
          <el-dropdown style="margin-left: 10px" @command="handleClick($event,row)" :hide-on-click="false">
            <el-button circle icon="el-icon-more" size="mini"/>
            <el-dropdown-menu slot="dropdown">
              <div v-if="row.containers.length > 1">
                <el-popover placement="left" trigger="hover">
                  <div v-for="c in row.containers" :key="c">
                    <p style="margin: 0">
                      <el-button :disabled="!checkExecPermissions()" @click="openTerminal(row, c)" type="text">{{ c }}</el-button>
                    </p>
                  </div>
                  <el-dropdown-item slot="reference" :disabled="!checkExecPermissions()" icon="el-icon-date" command="terminal">
                    {{ $t("commons.button.terminal") }}
                    <i class="el-icon-arrow-right"/>
                  </el-dropdown-item>
                </el-popover>
                <el-popover placement="left" trigger="hover">
                  <div v-for="c in row.containers" :key="c">
                    <p style="margin: 0">
                      <el-button :disabled="!checkLogPermissions()" @click="openTerminalLogs(row, c)" type="text">{{ c }}</el-button>
                    </p>
                  </div>
                  <el-dropdown-item slot="reference" :disabled="!checkLogPermissions()" icon="el-icon-notebook-2" command="logs">
                    {{ $t("commons.button.logs") }}
                    <i class="el-icon-arrow-right"/>
                  </el-dropdown-item>
                </el-popover>
                <el-popover placement="left" trigger="hover">
                  <div v-for="c in row.containers" :key="c">
                    <p style="margin: 0">
                      <el-button :disabled="!checkExecPermissions()" @click="openPodFiles(row, c)" type="text">{{ c }}</el-button>
                    </p>
                  </div>
                  <el-dropdown-item slot="reference" icon="el-icon-files" command="files">
                    {{ $t("business.pod.pod_file") }}
                    <i class="el-icon-arrow-right"/>
                  </el-dropdown-item>
                </el-popover>
              </div>
              <div v-if="row.containers.length == 1">
                <el-dropdown-item :disabled="!checkExecPermissions()" icon="iconfont iconline-terminalzhongduan" command="terminal">
                  {{ $t("commons.button.terminal") }}
                </el-dropdown-item>
                <el-dropdown-item :disabled="!checkLogPermissions()" icon="el-icon-tickets" command="logs">{{ $t("commons.button.logs") }}
                </el-dropdown-item>
                <el-dropdown-item :disabled="!checkExecPermissions()" icon="el-icon-files" command="files">{{ $t("business.pod.pod_file") }}
                </el-dropdown-item>
              </div>
              <el-dropdown-item icon="el-icon-download" command="download">{{ $t("commons.button.download_yaml") }}</el-dropdown-item>
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
import writeXlsxFile from "write-excel-file";
import { cpuUnitConvert, memoryUnitConvert } from "@/utils/unitConvert"
import { listPodMetrics } from "@/api/apis"

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
      orderField: null,
      orderMethod: null
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
    checkExecPermissions () {
      return checkPermissions({ scope: 'namespace', apiGroup: '', resource: 'pods/exec', verb: 'create' })
    },
    checkLogPermissions () {
      return checkPermissions({ scope: 'namespace', apiGroup: '', resource: 'pods/log', verb: 'get' })
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
        case "files":
          this.openPodFiles(row)
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
    openPodFiles(row,container) {
      let c
      if (container) {
        c = container
      } else {
        c = row.containers[0]
      }
      this.$router.push({ name: "PodFile", params: { namespace: row.metadata.namespace, name: row.metadata.name },query:{
          container: c
        }})
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
      this.$router.push({ name: "PodCreateYaml", query: { type: "pods" } })
    },
    onTop () {
      this.$router.push({ name: "PodTop" })
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
    doWithPodList(items){
          let result=[]
          if(items && items.length>0)
          for (const item of items) {
            let container = []
            for (const c of item.spec.containers) {
              container.push(c.name)
            }
            
            item.containers = container
            item.namespace=item.metadata.namespace
            item.status_phase=item.status.phase
            item.pod_status=this.getPodStatus(item)
            item.podIP=item.status.podIP
            item.nodeName=item.spec.nodeName
            item.creationTimestamp=item.metadata.creationTimestamp
            item.restart_count=Number(this.getRestartTimes(item))
            result.push(item)
          }
          if(!this.orderField || !this.orderMethod){
            return result
          } else{
            let orderMethod=this.orderMethod
             let orderField=this.orderField
            result=result.sort(function(a,b){
              if(orderMethod=='asc'){
                
                return (a[orderField]>b[orderField])?1:-1
              }else{
                return (a[orderField]<b[orderField])?1:-1
              }
            })
            return result
          }
    },
    search (resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      if(!this.orderField || !this.orderMethod){
      listWorkLoads(this.clusterName, "pods", true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize)
        .then((res) => {
          this.data =this.doWithPodList( res.items )
          this.paginationConfig.total = res.total
        }).finally(() => {
          this.loading = false
        })
      } else {
        let currentPage=this.paginationConfig.currentPage
        let pageSize=this.paginationConfig.pageSize
        listWorkLoads(this.clusterName, "pods", true, this.searchConfig.keywords)
        .then((res) => {
          this.data = this.doWithPodList( res.items ).slice(currentPage*pageSize-pageSize,currentPage*pageSize)
          this.paginationConfig.total = res.total
        }).finally(() => {
          this.loading = false
        })
      }  
    },
    sortTableFun(val){
       this.orderField=null
       this.orderMethod=null
       if (val.order) {
          this.orderField = val.prop
          this.orderMethod = (val.order == "descending" ? "desc" : "asc");
          this.search(false);
       }
    },
    /*导出配额信息为excel*/
    async exportToXlsx(){
      const schema = [
       {
        column: this.$t("commons.table.name"),
        type: String,
        value: (row) => row.metadata.name,
       },
       {
        column: this.$t("business.namespace.namespace"),
        type: String,
        value: (row) => row.metadata.namespace,
       },
       {
        column: "Pod IP",
        type: String,
        value: (row) => row.status.podIP,
       },
       {
        column: "NodeName",
        type: String,
        value: (row) => row.spec.nodeName,
       },
       {
        column: "requests cpu(m)",
        type: Number,
        value: (row) => {
             let result=0
             for(let i=0,s=row.spec.containers.length;i<s;i++){
               result=result+cpuUnitConvert(row.spec.containers[i].resources.requests.cpu)
             }
             return result
        },
       },
       {
        column: "requests memory(Mi)",
        type: Number,
        value: (row) => {
             let result=0
             for(let i=0,s=row.spec.containers.length;i<s;i++){
               result=result+memoryUnitConvert(row.spec.containers[i].resources.requests.memory)
             }
             return result
        },
       },
       {
        column: "limits cpu(m)",
        type: Number,
        value: (row) => {
             let result=0
             for(let i=0,s=row.spec.containers.length;i<s;i++){
               result=result+cpuUnitConvert(row.spec.containers[i].resources.limits.cpu)
             }
             return result
        },
       },
       {
        column: "limits memory(Mi)",
        type: Number,
        value: (row) => {
             let result=0
             for(let i=0,s=row.spec.containers.length;i<s;i++){
               result=result+memoryUnitConvert(row.spec.containers[i].resources.limits.memory)
             }
             return result
        },
       },
       {
        column: "use memory(Mi)",
        type: Number,
        value: (row) => row.memoryUsage,
       },
       {
        column: "use cpu(m)",
        type: Number,
        value: (row) => row.cpuUsage,
       },
      ];
      const data=await listWorkLoads(this.clusterName, "pods", true, this.searchConfig.keywords)
      
      const PodMetrics=await listPodMetrics(this.clusterName)
      const PodMetricsItems=  PodMetrics.items
      const PodMetricsMap={}
          for (const item of PodMetricsItems) {
            let cpu = 0
            let memory = 0
            for (const c of item.containers) {
              cpu += cpuUnitConvert(c.usage.cpu)
              memory += memoryUnitConvert(c.usage.memory)
            }
            PodMetricsMap[item.metadata.namespace+"|"+item.metadata.name] = {
              cpu: cpu,
              memory: memory,
            }
          }
      const pods =data.items;
      for(let i=0,s=pods.length;i<s;i++){
        if(PodMetricsMap[pods[i].metadata.namespace+"|"+pods[i].metadata.name]){
          const m=PodMetricsMap[pods[i].metadata.namespace+"|"+pods[i].metadata.name]
          pods[i].cpuUsage= m.cpu || 0
          pods[i].memoryUsage= m.memory || 0
        }
      }
      await writeXlsxFile((data.items||[]), {
         schema,
         fileName: "pods.xlsx",
      });
    }
  },
  mounted () {
    this.clusterName = this.$route.query.cluster
    this.search()
  },
}
</script>
