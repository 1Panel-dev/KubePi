<template>
  <div v-loading="loading">
    <h3>{{$t('business.workload.storage')}}</h3>
    <el-form label-position="top">
      <el-collapse v-model="activeNames">
        <el-collapse-item style="margin-top: 10px" :title="item._type" :name="index" v-for="(item, index) in form.spec.volumes" :key="index">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.name')">
                <div class="spanInFormStyle"><span :title="item._name">{{item._name}}</span></div>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.type')">
                <span>{{item._type}}</span>
              </el-form-item>
            </el-col>
          </el-row>

          <div v-if="item._type === 'NFS'">
            <el-row :gutter="20">
              <el-col v-if="item._path" :span="12">
                <el-form-item style="margin-left: 20px;" label="Path">
                  <span>{{item._path}}</span>
                </el-form-item>
              </el-col>
              <el-col v-if="item._server" :span="12">
                <el-form-item style="margin-left: 20px;" label="Server">
                  <span>{{item._server}}</span>
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <div v-if="item._type === 'PVC'">
            <el-row :gutter="20">
              <el-col v-if="item._detail" :span="12">
                <el-form-item style="margin-left: 20px;" label="claimName">
                  <div class="spanInFormStyle"><span :title="item._detail">{{item._detail}}</span></div>
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <div v-if="item._type === 'HostPath'">
            <el-row :gutter="20">
              <el-col v-if="item._path" :span="12">
                <el-form-item style="margin-left: 20px;" label="Path">
                  <div class="spanInFormStyle"><span :title="item._path">{{item._path}}</span></div>
                </el-form-item>
              </el-col>
              <el-col v-if="item._hostType" :span="12">
                <el-form-item style="margin-left: 20px;" label="HostType">
                  <div class="spanInFormStyle"><span :title="item._hostType">{{item._hostType}}</span></div>
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <div v-if="item._type === 'ConfigMap'">
            <el-row :gutter="20">
              <el-col v-if="item._detail" :span="12">
                <el-form-item style="margin-left: 20px;" label="ConfigMapName">
                  <div class="spanInFormStyle"><span :title="item._detail">{{item._detail}}</span></div>
                </el-form-item>
              </el-col>
              <el-col v-if="item._defaultMode" :span="12">
                <el-form-item style="margin-left: 20px;" label="DefaultMode">
                  <span>{{item._defaultMode}}</span>
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <div v-if="item._type === 'Secret'">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item style="margin-left: 20px;" label="SecretName">
                  <div class="spanInFormStyle"><span :title="item._detail">{{item._detail}}</span></div>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item style="margin-left: 20px;" label="DefaultMode">
                  <span>{{item._defaultMode}}</span>
                </el-form-item>
              </el-col>
            </el-row>
          </div>
        </el-collapse-item>
      </el-collapse>
    </el-form>
    <br>
  </div>
</template>

<script>
export default {
  name: "KoDetailUpdate",
  props: {
    yamlInfo: Object,
  },
  watch: {
    yamlInfo: {
      handler(newYamlInfo) {
        if (newYamlInfo) {
          if (newYamlInfo.spec?.volumes) {
            this.form = newYamlInfo
            for (const item of this.form.spec.volumes) {
              item._name = item.name
              if (item.persistentVolumeClaim) {
                item._type = "PVC"
                item._detail = item.persistentVolumeClaim.claimName
                item._defaultMode = item.persistentVolumeClaim.defaultMode
                continue
              }
              if (item.secret) {
                item._type = "Secret"
                item._detail = item.secret.secretName
                item._defaultMode = item.secret.defaultMode
                continue
              }
              if (item.configMap) {
                item._type = "ConfigMap"
                item._detail = item.configMap.name
                item._defaultMode = item.configMap.defaultMode
                continue
              }
              if (item.emptyDir) {
                item._type = "EmptyDir"
                continue
              }
              if (item.nfs) {
                item._type = "NFS"
                item._path = item.nfs.path
                item._server = item.nfs.server
                continue
              }
              if (item.hostPath) {
                item._type = "HostPath"
                item._path = item.hostPath.path
                item._hostType = item.hostPath.type
                continue
              }
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
      form: {
        spec: {
          volumes: [],
        },
      },
    }
  },
}
</script>

<style scoped>
</style>
