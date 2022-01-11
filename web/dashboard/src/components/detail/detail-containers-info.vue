<template>
  <div v-loading="loading">
    <h3 style="display: inline-block;">{{container.name}}</h3>
    <el-form label-position="top">
      <h4 style="display: inline-block;">{{$t('business.pod.image')}}</h4>
      <el-row>
        <el-form-item style="margin-left: 20px;" :label="$t('business.workload.pull_policy')">
          <span>{{container.imagePullPolicy}}</span>
        </el-form-item>
      </el-row>
      <el-row>
        <el-form-item style="margin-left: 20px;" :label="$t('business.workload.container_image')">
          <span>{{container.image}}</span>
        </el-form-item>
      </el-row>

      <div v-if="container.ports">
        <h4 style="display: inline-block;">{{$t('business.network.port')}}</h4>
        <div v-for="(item, index) in container.ports" :key="index">
          <el-row :gutter="20">
            <el-col v-if="item.name" :span="6">
              <el-form-item style="margin-left: 20px;" :label="$t('commons.table.name')">
                <span>{{item.name}}</span>
              </el-form-item>
            </el-col>
            <el-col v-if="item.containerPort" :span="6">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.container_port')">
                <span>{{item.containerPort}}</span>
              </el-form-item>
            </el-col>
            <el-col v-if="item.protocol" :span="6">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.protocol')">
                <span>{{item.protocol}}</span>
              </el-form-item>
            </el-col>
            <el-col v-if="item.hostPort" :span="6">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.public_port')">
                <span>{{item.hostPort}}</span>
              </el-form-item>
            </el-col>
          </el-row>
        </div>
      </div>

      <div v-if="container.env">
        <h4 style="display: inline-block;">{{$t('business.workload.variable')}}</h4>
        <el-form-item style="margin-left: 20px;" :label="$t('business.workload.variable')">
          <el-tag type="success" style="margin-right: 10px;" v-for="(item, index) in container.env" :key="index">{{item.name + " = " + item.value}}</el-tag>
        </el-form-item>
      </div>

      <div v-if="container.workingDir || container.command || container.args">
        <h4 style="display: inline-block;">{{$t('business.workload.command')}}</h4>
        <el-row v-if="container.workingDir">
          <el-form-item style="margin-left: 20px;" :label="$t('business.workload.working_dir')">
            <span>{{container.workingDir}}</span>
          </el-form-item>
        </el-row>
        <el-row v-if="container.command">
          <el-form-item style="margin-left: 20px;" :label="$t('business.workload.commands')">
            <el-tag type="success" style="margin-left: 20px;" v-for="(item, index) in container.command" :key="index">{{item}}</el-tag>
          </el-form-item>
        </el-row>
        <el-row v-if="container.args">
          <el-form-item style="margin-left: 20px;" :label="$t('business.workload.arguments')">
            <el-tag type="success" style="margin-left: 20px;" v-for="(item, index) in container.args" :key="index">{{item}}</el-tag>
          </el-form-item>
        </el-row>
      </div>

      <div v-if="container.volumeMounts">
        <h4 style="display: inline-block;">{{$t('business.workload.mount')}}</h4>
        <div v-for="(item, index) in container.volumeMounts" :key="index">
          <el-row :gutter="20">
            <el-col v-if="item.name" :span="6">
              <el-form-item style="margin-left: 20px;" :label="$t('commons.table.name')">
                <div class="spanInFormStyle"><span :title="item.name">{{item.name}}</span></div>
              </el-form-item>
            </el-col>
            <el-col v-if="item.mountPath" :span="6">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.path')">
                <div class="spanInFormStyle"><span :title="item.mountPath">{{item.mountPath}}</span></div>
              </el-form-item>
            </el-col>
            <el-col v-if="item.readOnly" :span="6">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.read_only')">
                <span>{{item.readOnly}}</span>
              </el-form-item>
            </el-col>
            <el-col v-if="item.subPath" :span="6">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.sub_path_in_volume')">
                <div class="spanInFormStyle"><span :title="item.subPath">{{item.subPath}}</span></div>
              </el-form-item>
            </el-col>
          </el-row>
        </div>
      </div>

      <div v-if="container.resources.limits || container.resources.requests">
        <h4 style="display: inline-block;">{{$t('business.workload.resource')}}</h4>
        <el-row :gutter="20">
          <div v-if="container.resources.requests">
            <el-col v-if="container.resources.requests.cpu" :span="6">
              <el-form-item style="margin-left: 20px;" :label="'CPU ' + $t('business.workload.reservation') + ' (mCPUs)'">
                <span>{{container.resources.requests.cpu}}</span>
              </el-form-item>
            </el-col>
            <el-col v-if="container.resources.requests.cpu" :span="6">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.memory') + $t('business.workload.reservation') + ' (Mi)'">
                <span>{{container.resources.requests.memory}}</span>
              </el-form-item>
            </el-col>
          </div>
          <div v-if="container.resources.limits">
            <el-col v-if="container.resources.limits.cpu" :span="6">
              <el-form-item style="margin-left: 20px;" :label="'CPU ' + $t('business.workload.limit') + ' (mCPUs)'">
                <span>{{container.resources.limits.cpu}}</span>
              </el-form-item>
            </el-col>
            <el-col v-if="container.resources.limits.memory" :span="6">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.memory') + $t('business.workload.limit') + ' (Mi)'">
                <span>{{container.resources.limits.memory}}</span>
              </el-form-item>
            </el-col>
          </div>
        </el-row>
      </div>

      <div v-if="container.livenessProbe || container.readinessProbe || container.startupProbe">
        <h4 style="display: inline-block;">{{$t('business.workload.health_check')}}</h4>
        <div v-for="(item, index) in container.healthCheck" :key="index">
          <h5 style="display: inline-block;margin-left: 10px;">{{item._type}}</h5>
          <div v-if="item._model.exec">
            <el-row :gutter="20">
              <el-col :span="4">
                <el-form-item style="margin-left: 20px;" :label="$t('business.workload.type')">
                  <span>exec</span>
                </el-form-item>
              </el-col>
              <el-col v-if="item._model.exec.command" :span="20">
                <el-form-item style="margin-left: 20px;" :label="$t('business.workload.command')">
                  <div v-for="(item, index) in item._model.exec.command" :key="index">
                    <el-tag v-if="item.length < 200" type="success">{{item}}</el-tag>
                    <div v-else style="background-color: #1F261E;line-height: 20px;"><span class="spanStyle">{{item}}</span></div>
                  </div>
                </el-form-item>
              </el-col>
            </el-row>
          </div>
          <div v-if="item._model.tcpSocket">
            <el-row :gutter="20">
              <el-col :span="4">
                <el-form-item style="margin-left: 20px;" :label="$t('business.workload.type')">
                  <span>tcpSocket</span>
                </el-form-item>
              </el-col>
              <el-col v-if="item._model.exec.port" :span="8">
                <el-form-item style="margin-left: 20px;" :label="$t('business.network.port')">
                  <span>{{item._model.exec.port}}</span>
                </el-form-item>
              </el-col>
            </el-row>
          </div>
          <div v-if="item._model.httpGet">
            <el-row :gutter="20">
              <el-col :span="4">
                <el-form-item style="margin-left: 20px;" :label="$t('business.workload.type')">
                  <span>httpGet</span>
                </el-form-item>
              </el-col>
              <el-col v-if="item._model.httpGet.path" :span="8">
                <el-form-item style="margin-left: 20px;" :label="$t('business.workload.check_path')">
                  <span>{{item._model.httpGet.path}}</span>
                </el-form-item>
              </el-col>
              <el-col v-if="item._model.httpGet.port" :span="8">
                <el-form-item style="margin-left: 20px;" :label="$t('business.workload.check_port')">
                  <span>{{item._model.httpGet.port}}</span>
                </el-form-item>
              </el-col>
            </el-row>
          </div>
          <div v-if="item._model.httpsGet">
            <el-row :gutter="20">
              <el-col :span="4">
                <el-form-item style="margin-left: 20px;" :label="$t('business.workload.type')">
                  <span>httpsGet</span>
                </el-form-item>
              </el-col>
              <el-col v-if="item._model.httpsGet.path" :span="8">
                <el-form-item style="margin-left: 20px;" :label="$t('business.workload.check_path')">
                  <span>{{item._model.httpsGet.path}}</span>
                </el-form-item>
              </el-col>
              <el-col v-if="item._model.httpsGet.port" :span="8">
                <el-form-item style="margin-left: 20px;" :label="$t('business.workload.check_port')">
                  <span>{{item._model.httpsGet.port}}</span>
                </el-form-item>
              </el-col>
            </el-row>
          </div>
          <el-row :gutter="20">
            <el-col v-if="item._model.failureThreshold" :span="4">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.failure_threshold')">
                <span>{{item._model.failureThreshold}}</span>
              </el-form-item>
            </el-col>
            <el-col v-if="item._model.initialDelaySeconds" :span="4">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.initial_delay')">
                <span>{{item._model.initialDelaySeconds}}</span>
              </el-form-item>
            </el-col>
            <el-col v-if="item._model.periodSeconds" :span="4">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.check_interval')">
                <span>{{item._model.periodSeconds}}</span>
              </el-form-item>
            </el-col>
            <el-col v-if="item._model.successThreshold" :span="4">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.seccess_threshold')">
                <span>{{item._model.successThreshold}}</span>
              </el-form-item>
            </el-col>
            <el-col v-if="item._model.timeoutSeconds" :span="4">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.timeout')">
                <span>{{item._model.timeoutSeconds}}</span>
              </el-form-item>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-form>
    <br>
  </div>
</template>

<script>
export default {
  name: "KoDetailContainerInfo",
  props: {
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
    }
  },
}
</script>

<style scoped>
.spanStyle {
  display: block;
  padding-left: 10px;
  margin-top: 5px;
  font-size: 12px;
  color:#67c23a;
  white-space: pre;
}
</style>
