<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'PersistentVolumes'}" v-loading="loading">
    <div class="grid-content bg-purple-light">
      <div v-if="yamlShow">
        <yaml-editor :value="yaml" :is-edit="true" ref="yaml_editor"></yaml-editor>
      </div>
      <div>
        <div style="float: right;margin-top: 10px">
          <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
          <el-button v-if="!yamlShow" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
          <el-button v-loading="loading" @click="onSubmit" type="primary">
            {{ $t("commons.button.submit") }}
          </el-button>
        </div>
      </div>
    </div>
  </layout-content>

</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import {listStorageClasses} from "@/api/storageclass";
import {getPvc, updatePvc} from "@/api/pvc";

export default {
  name: "PersistentVolumeClaimEdit",
  props: {
    name: String,
    namespace: String,
  },
  components: { YamlEditor, LayoutContent },
  data () {
    return {
      loading: false,
      yamlShow: true,
      form: {
        apiVersion: "v1",
        kind: "PersistentVolumeClaim",
        metadata: {
          name: "",
          namespace: "",
        },
        spec: {
          accessModes: [],
          resources: {
            requests: {
              storage: ""
            }
          },
        },
        status: {
          phase: "",
          accessModes: [],
          capacity: []
        }
      },
      yaml: {},
      cluster: "",
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getPvc(this.cluster, this.namespace, this.name).then(res => {
        this.form = res
        this.yaml = JSON.parse(JSON.stringify(this.form))
        this.currentStorageCapacity = parseInt(this.form.spec.resources.requests.storage)
      }).finally(() => {
        this.loading = false
      })
      if (this.yamlShow) {
        this.yaml = this.transformYaml()
      }
    },
    handleClick (tab) {
      this.activeName = tab.index
    },
    onCancel () {
      this.$router.push({ name: "PersistentVolumeClaim" })
    },
    onEditYaml () {
      this.yamlShow = true
      this.yaml = this.transformYaml()
    },
    onSubmit () {
      let data = {}
      if (this.yamlShow) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        data = this.transformYaml()
      }
      this.loading = true
      updatePvc(this.cluster, this.namespace, data,this.name ).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "PersistentVolumeClaim" })
      }).finally(() => {
        this.loading = false
      })
    },
    loadStorageClasses () {
      this.storageClasses = []
      listStorageClasses(this.cluster).then((res) => {
        this.storageClasses.push('None')
        for (const sc of res.items){
          this.storageClasses.push(sc.metadata.name)
        }
      })
    },
    transformYaml () {
      if (this.form.spec.nodeAffinity) {
        this.$refs.ko_node_scheduling.transformation(this.form.spec)
      }
      return JSON.parse(JSON.stringify(this.form))
    },
    setFormAccessModel() {
      this.form.spec.accessModes = []
      this.form.spec.accessModes.push(this.currentAccessModes)
    },
    setStorageCapacity() {
      this.form.spec.capacity.storage = this.currentStorageCapacity.toString() + 'Gi'
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.getDetail()
  }
}
</script>

<style scoped>

</style>
