<template>
  <div style="margin-top: 20px">
    <ko-card title="Target">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Target Reference">
              <el-select filterable clearable v-model="form.scaleTargetRef" value-key="name" style="width:100%"
                         @change="transform" required>
                <el-option v-for="d in deploymentItem.items"
                           :key="d.metadata.name"
                           :label="d.metadata.name"
                           :value="{
                              kind: deploymentItem.kind,
                              apiVersion: deploymentItem.apiVersion,
                              name: d.metadata.name
                           }">
                  <span style="float: left">{{ d.metadata.name }}</span>
                  <span style="float: right; color: #8492a6; font-size: 13px">{{ deploymentItem.kind }}</span>
                </el-option>
                <el-option v-for="d in replicaSetItem.items"
                           :key="d.metadata.name"
                           :label="d.metadata.name"
                           :value="{
                              kind: replicaSetItem.kind,
                              apiVersion: replicaSetItem.apiVersion,
                              name: d.metadata.name
                           }">
                  <span style="float: left">{{ d.metadata.name }}</span>
                  <span style="float: right; color: #8492a6; font-size: 13px">{{ replicaSetItem.kind }}</span>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Minimum Replicas">
              <el-input type="number"  v-model.number="form.minReplicas" required></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Maximum Replicas">
              <el-input type="number"  v-model.number="form.maxReplicas" required></el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </ko-card>
  </div>
</template>

<script>
import KoCard from "@/components/ko-card"
// import KoFormItem from "@/components/ko-form-item"
import {listDeploymentsByNs} from "@/api/deployments"
import {listNsReplicaSets} from "@/api/replicasets"

export default {
  name: "KoHpaTarget",
  components: { KoCard },
  props: {
    namespace: String,
    cluster: String,
    specObj: Object
  },
  data () {
    return {
      form: {
        scaleTargetRef: {}
      },
      data: [{
        metadata: {}
      }],
      deploymentItem: {
        apiVersion: "",
        items: [],
        kind: "Deployment"
      },
      replicaSetItem: {
        apiVersion: "",
        items: [],
        kind: "ReplicaSet"
      },
    }
  },
  watch: {
    namespace (oldValue, newValue) {
      if (oldValue !== newValue) {
        this.form = {
          scaleTargetRef: {}
        }
      }
      this.getReferences()
    }
  },
  methods: {
    getReferences () {
      listDeploymentsByNs(this.cluster, this.namespace).then(res => {
        this.deploymentItem = {
          items: res.items,
          apiVersion: res.apiVersion,
          kind: "Deployment"
        }
      })
      listNsReplicaSets(this.cluster, this.namespace).then(res => {
        this.replicaSetItem = {
          items: res.items,
          apiVersion: res.apiVersion,
          kind: "ReplicaSet"
        }
      })
    },
    transform () {
      this.$emit("update:specObj", this.form)
    }
  },
  mounted () {
    this.getReferences()
    if (this.specObj) {
      this.form = this.specObj
    }
  }
}
</script>

<style scoped>

</style>
