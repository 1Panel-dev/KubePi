<template>
  <div>
    <h3 style="display: inline-block;">{{ $t("business.workload.container_image") }}</h3>
    <el-button style="margin-left: 10px;" type="text" icon="el-icon-edit" @click="dialogModifyVersionVisible = true">{{$t('commons.button.click_to_edit')}}</el-button>
    <li style="margin-left: 20px;font-size: 16px;" v-for="(item, index) in imagesData" :key="index">{{item.name}}:{{item.version}}</li>

    <el-dialog :title="$t('commons.button.scale')" width="70%" :close-on-click-modal="false" :visible.sync="dialogModifyVersionVisible">
      <complex-table :data="imagesData" v-loading="loading">
        <el-table-column :label="$t('business.workload.container_type')" prop="type" min-width="10" />
        <el-table-column :label="$t('business.workload.name')" prop="name" min-width="20" show-overflow-tooltip />
        <el-table-column :label="$t('business.workload.container_image')" prop="image" min-width="40" show-overflow-tooltip />
        <el-table-column :label="$t('business.workload.current_version')" prop="version" min-width="15" />
        <el-table-column :label="$t('business.workload.new_version')" prop="newVersion" min-width="20">
          <template v-slot:default="{row}">
            <ko-form-item itemType="input" v-model="row.newVersion" />
          </template>
        </el-table-column>
      </complex-table>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="dialogModifyVersionVisible = false">{{ $t("commons.button.cancel") }}</el-button>
        <el-button size="small" @click="updateForm">{{ $t("commons.button.confirm") }}</el-button>
      </div>
    </el-dialog>
  </div>

</template>

<script>
import { updateWorkLoad } from "@/api/workloads"
import KoFormItem from "@/components/ko-form-item/index"
import ComplexTable from "@/components/complex-table"

export default {
  name: "KoDetailImage",
  components: { ComplexTable, KoFormItem },
  props: {
    cluster: String,
    yamlInfo: Object,
    resourceType: String,
  },
  watch: {
    yamlInfo: {
      handler(newYamlInfo) {
        if (newYamlInfo) {
          if (newYamlInfo.spec) {
            this.imagesData = []
            this.form = newYamlInfo
            this.form.namespace = newYamlInfo.metadata.namespace
            if (newYamlInfo.spec.template.spec.initContainers) {
              for (const c of newYamlInfo.spec.template.spec.containers) {
                let index = c.image.lastIndexOf(":")
                this.imagesData.push({
                  type: this.$t("business.workload.initContainer"),
                  name: c.name,
                  image: index !== -1 ? c.image.substring(0, index) : "",
                  version: index !== -1 ? c.image.substring(index + 1, c.image.length) : "",
                  newVersion: index !== -1 ? c.image.substring(index + 1, c.image.length) : "",
                })
              }
            }
            for (const c of newYamlInfo.spec.template.spec.containers) {
              let index = c.image.lastIndexOf(":")
              this.imagesData.push({
                type: this.$t("business.workload.standardContainer"),
                name: c.name,
                image: index !== -1 ? c.image.substring(0, index) : "",
                version: index !== -1 ? c.image.substring(index + 1, c.image.length) : "",
                newVersion: index !== -1 ? c.image.substring(index + 1, c.image.length) : "",
              })
            }
          }
          this.loading = false
        }
      },
      immediate: true,
      deep: true,
    },
  },
  data() {
    return {
      loading: true,
      dialogModifyVersionVisible: false,
      form: {},
      imagesData: [],
    }
  },
  methods: {
    updateForm() {
      this.loading = true
      for (const c of this.form.spec.template.spec.containers) {
        for (const i of this.imagesData) {
          if (c.name == i.name && i.type === this.$t("business.workload.standardContainer")) {
            c.image = i.image + ":" + (i.newVersion !== "" ? i.newVersion : i.version)
            break
          }
        }
      }
      if (this.form.spec.template.spec.initContainers) {
        for (const c of this.form.spec.template.spec.initContainers) {
          for (const i of this.imagesData) {
            if (c.name == i.name && i.type === this.$t("business.workload.initContainer")) {
              c.image = i.image + ":" + (i.newVersion !== "" ? i.newVersion : i.version)
              break
            }
          }
        }
      }
      updateWorkLoad(this.cluster, this.resourceType, this.form.metadata.namespace, this.form.metadata.name, this.form)
        .then(() => {
          this.dialogModifyVersionVisible = false
          this.loading = true
          this.$message({
            type: "success",
            message: this.$t("commons.msg.operation_success"),
          })
        })
        .finally(() => {
          this.loading = false
        })
    },
  },
}
</script>

<style scoped>
</style>
