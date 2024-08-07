<template>
  <div v-loading="loading">
    <h3 style="display: inline-block;">{{container.name}}</h3>
    <el-form>
      <el-row :gutter="20" style="margin-top: 20px" class="row-box">
        <el-col :span="8">
          <el-card class="el-card" :title="$t('business.pod.image')">
            <el-form-item label="State" v-if="containerStatusStates[container.name]!=''">
              <span><pre>{{containerStatusStates[container.name]}}</pre></span>
            </el-form-item>
            <el-form-item label="LastState" v-if="containerStatusLastStates[container.name]!=''">
              <span></span><pre>{{containerStatusLastStates[container.name]}}</pre></span>
            </el-form-item>
            <el-form-item :label="$t('business.workload.pull_policy')">
              <span>{{container.imagePullPolicy}}</span>
            </el-form-item>
            <el-form-item :label="$t('business.workload.container_image')">
              <span>{{container.image}}</span>
            </el-form-item>

            <div v-if="container.workingDir || container.command || container.args">
              <h4 style="display: inline-block;">{{$t('business.workload.command')}}</h4>
              <el-form-item v-if="container.workingDir" :label="$t('business.workload.working_dir')">
                <span>{{container.workingDir}}</span>
              </el-form-item>
              <el-form-item v-if="container.command" :label="$t('business.workload.commands')">
                <el-tag type="success" v-for="(item, index) in container.command" :key="index">{{item}}</el-tag>
              </el-form-item>
              <el-form-item v-if="container.args" :label="$t('business.workload.arguments')">
                <el-tag type="success" v-for="(item, index) in container.args" :key="index">{{item}}</el-tag>
              </el-form-item>
            </div>

            <div v-if="container.ports">
              <h4 style="display: inline-block;">{{$t('business.workload.port')}}</h4>
              <complex-table :header-cell-style="{background: '#19191c'}" :data="container.ports">
                <el-table-column :label="$t('commons.table.name')" prop="name" show-overflow-tooltip fix />
                <el-table-column :label="$t('business.workload.container_port')" prop="containerPort" show-overflow-tooltip fix />
                <el-table-column :label="$t('business.workload.protocol')" prop="protocol" show-overflow-tooltip fix />
                <el-table-column :label="$t('business.workload.host_port')" prop="hostPort" show-overflow-tooltip fix />
              </complex-table>
            </div>
          </el-card>
        </el-col>

        <el-col :span="16" v-if="container.volumeMounts">
          <el-card class="el-card" :title="$t('business.network.port')">
            <div v-if="container.env">
              <h4 style="display: inline-block;">{{$t('business.workload.variable')}}</h4>
              <el-form-item :label="$t('business.workload.variable')">
                <el-tag type="success" style="margin-right: 10px;" v-for="(item, index) in container.env" :key="index">{{item.name + " = " + item.value}}</el-tag>
              </el-form-item>
            </div>

            <div v-if="container.volumeMounts">
              <h4 style="display: inline-block;">{{$t('business.workload.mount')}}</h4>
              <complex-table :header-cell-style="{background: '#19191c'}" :data="container.volumeMounts">
                <el-table-column :label="$t('commons.table.name')" prop="name" show-overflow-tooltip fix />
                <el-table-column :label="$t('business.workload.path')" prop="mountPath" show-overflow-tooltip fix />
                <el-table-column :label="$t('business.workload.read_only')" prop="readOnly" show-overflow-tooltip fix />
                <el-table-column :label="$t('business.workload.sub_path_in_volume')" prop="subPath" show-overflow-tooltip fix />
              </complex-table>
            </div>

            <div v-if="container.resources && (container.resources.limits || container.resources.requests)">
              <h4 style="display: inline-block;">{{$t('business.workload.resource')}}</h4>
              <el-row :gutter="20">
                <div v-if="container.resources && container.resources.requests">
                  <el-col v-if="container.resources.requests.cpu" :span="12">
                    <el-form-item :label="'CPU ' + $t('business.workload.reservation') + ' (mCPUs)'">
                      <span>{{container.resources.requests.cpu}}</span>
                    </el-form-item>
                  </el-col>
                  <el-col v-if="container.resources.requests.cpu" :span="12">
                    <el-form-item :label="$t('business.workload.memory') + $t('business.workload.reservation') + ' (Mi)'">
                      <span>{{container.resources.requests.memory}}</span>
                    </el-form-item>
                  </el-col>
                </div>
                <div v-if="container.resources && container.resources.limits">
                  <el-col v-if="container.resources.limits.cpu" :span="12">
                    <el-form-item :label="'CPU ' + $t('business.workload.limit') + ' (mCPUs)'">
                      <span>{{container.resources.limits.cpu}}</span>
                    </el-form-item>
                  </el-col>
                  <el-col v-if="container.resources.limits.memory" :span="12">
                    <el-form-item :label="$t('business.workload.memory') + $t('business.workload.limit') + ' (Mi)'">
                      <span>{{container.resources.limits.memory}}</span>
                    </el-form-item>
                  </el-col>
                </div>
              </el-row>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <div v-if="container.livenessProbe || container.readinessProbe || container.startupProbe">
        <h4>{{$t('business.workload.health_check')}}</h4>
        <el-row :gutter="20" class="row-box">
          <el-col v-for="(item, index) in container.healthCheck" :key="index" :span="6">
            <el-card class="el-card">
              <h4>{{item._type}}</h4>
              <div v-if="item._model.exec">
                <el-form-item :label="$t('business.workload.type')">
                  <span>exec</span>
                </el-form-item>
                <el-form-item v-if="item._model.exec && item._model.exec.command" :label="$t('business.workload.command')">
                  <div v-for="(item, index) in item._model.exec.command" :key="index">
                    <el-tag v-if="item.length < 200" type="success">{{item}}</el-tag>
                    <div v-else style="background-color: #1F261E;line-height: 20px;"><span class="spanStyle">{{item}}</span></div>
                  </div>
                </el-form-item>
              </div>
              <div v-if="item._model.tcpSocket">
                <el-form-item :label="$t('business.workload.type')">
                  <span>tcpSocket</span>
                </el-form-item>
                <el-form-item v-if="item._model.exec && item._model.exec.port" :label="$t('business.network.port')">
                  <span>{{item._model.exec.port}}</span>
                </el-form-item>
              </div>
              <div v-if="item._model.httpGet">
                <el-form-item :label="$t('business.workload.type')">
                  <span>httpGet</span>
                </el-form-item>
                <el-form-item v-if="item._model.httpGet && item._model.httpGet.path" :label="$t('business.workload.check_path')">
                  <span>{{item._model.httpGet.path}}</span>
                </el-form-item>
                <el-form-item v-if="item._model.httpGet && item._model.httpGet.port" :label="$t('business.workload.check_port')">
                  <span>{{item._model.httpGet.port}}</span>
                </el-form-item>
              </div>
              <div v-if="item._model.httpsGet">
                <el-form-item :label="$t('business.workload.type')">
                  <span>httpsGet</span>
                </el-form-item>
                <el-form-item v-if="item._model.httpGet &&  item._model.httpsGet.path" :label="$t('business.workload.check_path')">
                  <span>{{item._model.httpsGet.path}}</span>
                </el-form-item>
                <el-form-item v-if="item._model.httpsGet && item._model.httpsGet.port" :label="$t('business.workload.check_port')">
                  <span>{{item._model.httpsGet.port}}</span>
                </el-form-item>
              </div>
              <el-form-item v-if="item._model.failureThreshold" :label="$t('business.workload.failure_threshold')">
                <span>{{item._model.failureThreshold}}</span>
              </el-form-item>
              <el-form-item v-if="item._model.initialDelaySeconds" :label="$t('business.workload.initial_delay')">
                <span>{{item._model.initialDelaySeconds}}</span>
              </el-form-item>
              <el-form-item v-if="item._model.periodSeconds" :label="$t('business.workload.check_interval')">
                <span>{{item._model.periodSeconds}}</span>
              </el-form-item>
              <el-form-item v-if="item._model.successThreshold" :label="$t('business.workload.seccess_threshold')">
                <span>{{item._model.successThreshold}}</span>
              </el-form-item>
              <el-form-item v-if="item._model.timeoutSeconds" :label="$t('business.workload.timeout')">
                <span>{{item._model.timeoutSeconds}}</span>
              </el-form-item>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-form>
    <br>
    <el-row style="margin-top: 20px">
          <KoDetailContainersMetrics :containerInfo="containerInfo" :cluster="cluster"  :yamlInfo="yamlInfo"/>
    </el-row>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"
import KoDetailContainersMetrics from "@/components/detail/detail-containers-metrics.vue"
export default {
  name: "KoDetailContainerInfo",
  components: { ComplexTable, KoDetailContainersMetrics },
  props: {
    cluster: String,
    yamlInfo: Object,
    containerInfo: Object,
  },
  watch: {
    yamlInfo: {
      handler(newYamlInfo) {
        if (newYamlInfo) {
          if (newYamlInfo.spec.containers) {
            if (this.containerInfo.type === "init-container") {
              for (const item of newYamlInfo.spec.initContainers) {
                if (item.name === this.containerInfo.name) {
                  this.container = item
                  break
                }
              }
            } else {
              for (const item of newYamlInfo.spec.containers) {
                if (item.name === this.containerInfo.name) {
                  this.container = item
                  break
                }
              }
            }
            this.container.healthCheck = []
            if (this.container.livenessProbe) {
              this.container.healthCheck.push({ _type: this.$t("business.workload.liveness_check"), _model: this.container.livenessProbe })
            }
            if (this.container.readinessProbe) {
              this.container.healthCheck.push({ _type: this.$t("business.workload.readiness_check"), _model: this.container.readinessProbe })
            }
            if (this.container.startupProbe) {
              this.container.healthCheck.push({ _type: this.$t("business.workload.startup_check"), _model: this.container.startupProbe })
            }
            if (this.container.healthCheck.length === 0) {
              delete this.container.healthCheck
            }
            this.loading = false
          }
          this.containerStatuses={}
         
          for(let i=0,s=newYamlInfo.status.containerStatuses.length;i<s;i++){
            const item=newYamlInfo.status.containerStatuses[i]
            this.containerStatuses[item.name] = item
            this.containerStatusStates[item.name] =this.getContainerStatusState(item.name)
            this.containerStatusLastStates[item.name] =this.getContainerStatusLastState(item.name)
          }
        }
      },
      immediate: true,
      deep: true,
    },
    containerInfo: {
      handler(newValue) {
        if (newValue) {
          if (this.yamlInfo.spec.containers) {
            if (newValue.type === "init-container") {
              for (const item of this.yamlInfo.spec.initContainers) {
                if (item.name === newValue.name) {
                  this.container = item
                  break
                }
              }
            } else {
              for (const item of this.yamlInfo.spec.containers) {
                if (item.name === newValue.name) {
                  this.container = item
                  break
                }
              }
            }
            this.container.healthCheck = []
            if (this.container.livenessProbe) {
              this.container.healthCheck.push({ _type: this.$t("business.workload.liveness_check"), _model: this.container.livenessProbe })
            }
            if (this.container.readinessProbe) {
              this.container.healthCheck.push({ _type: this.$t("business.workload.readiness_check"), _model: this.container.readinessProbe })
            }
            if (this.container.startupProbe) {
              this.container.healthCheck.push({ _type: this.$t("business.workload.startup_check"), _model: this.container.startupProbe })
            }
            if (this.container.healthCheck.length === 0) {
              delete this.container.healthCheck
            }
            this.loading = false
          }
        }
      },
      immediate: true,
      deep: true,
    },
  },
  data() {
    return {
      loading: true,
      activeNames: 0,
      container: {
        image: "",
        resources: {},
      },
      containerStatuses:{},
      containerStatusStates:{},
      containerStatusLastStates:{}
    }
  },
  methods:{
     getContainerStatusState(container_name){
       
       const status= this.containerStatuses[container_name]

       let state ="\n"
       if(status.state){
        for(let key in status.state) {
          state = state +key +"\n" 
          for(let attr in  status.state[key]){
            state = state +attr+":"+status.state[key][attr] +"\n" 
          }
        }
       }
       return state
     },
     getContainerStatusLastState(container_name){
       
       const status= this.containerStatuses[container_name]

       let state ="\n"
       if(status.lastState){
        for(let key in status.lastState) {
          state = state +key +"\n" 
          for(let attr in  status.lastState[key]){
            state = state +attr+":"+status.lastState[key][attr] +"\n" 
          }
        }
       }
       return state
     }
  }
}
</script>

<style scoped>
.spanStyle {
  display: block;
  padding-left: 10px;
  margin-top: 5px;
  font-size: 12px;
  color: #67c23a;
  white-space: pre;
}
</style>
