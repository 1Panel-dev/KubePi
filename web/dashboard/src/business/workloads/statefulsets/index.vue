<template>
  <layout-content header="StatefulSets">
    <complex-table :selects.sync="selects" :data="data" v-loading="loading" :pagination-config="paginationConfig" :search-config="searchConfig" @search="search">
      <template #header>
        <el-button type="primary" size="small" @click="yamlCreate" v-has-permissions="{scope:'namespace',apiGroup:'apps',resource:'statefulsets',verb:'create'}">
          YAML
        </el-button>
        <el-button type="primary" size="small" @click="onCreate" v-has-permissions="{scope:'namespace',apiGroup:'apps',resource:'statefulsets',verb:'create'}">
          {{ $t("commons.button.create") }}
        </el-button>
        <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()" v-has-permissions="{scope:'namespace',apiGroup:'apps',resource:'statefulsets',verb:'delete'}">
          {{ $t("commons.button.delete") }}
        </el-button>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column sortable :label="$t('commons.table.name')" prop="name" min-width="120" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column sortable :label="$t('business.namespace.namespace')" min-width="80" prop="metadata.namespace" />
      <el-table-column sortable :label="$t('commons.table.status')" min-width="40">
        <template v-slot:default="{row}">
          {{ row.status.readyReplicas || 0 }} / {{ row.status.replicas }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" min-width="60" prop="metadata.creationTimestamp" fix>
        <template v-slot:default="{row}">
          {{ row.metadata.creationTimestamp | age }}
        </template>
      </el-table-column>
      <ko-table-operations :buttons="buttons" :label="$t('commons.table.action')"></ko-table-operations>
    </complex-table>

    <el-dialog :title="$t('commons.button.scale')" width="30%" :close-on-click-modal="false" :visible.sync="dialogScaleVisible">
      <el-form label-position="top" style="margin-left: 30px" ref="form" :model="form">
        <div>
          <el-form-item :label="$t('business.common.basic')">
            <div><span>{{form.namespace}} / {{form.name}}</span></div>
          </el-form-item>
          <el-form-item :label="$t('business.workload.replicas')" prop="replicas">
            <ko-form-item itemType="number" v-model.number="form.replicas" />
          </el-form-item>
        </div>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="dialogScaleVisible = false">{{ $t("commons.button.cancel") }}</el-button>
        <el-button size="small" @click="onScaleSubmit">{{ $t("commons.button.confirm") }}</el-button>
      </div>
    </el-dialog>

    <el-dialog :title="$t('commons.button.scale')" width="70%" :close-on-click-modal="false" :visible.sync="dialogModifyVersionVisible">
      <complex-table :data="form.imagesData" v-loading="loading">
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
        <el-button size="small" @click="onModifyVersionSubmit">{{ $t("commons.button.confirm") }}</el-button>
      </div>
    </el-dialog>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import KoFormItem from "@/components/ko-form-item/index"
import { scaleStatefulset, patchStatefulset } from "@/api/statefulsets"
import { listWorkLoads, deleteWorkLoad } from "@/api/workloads"
import { downloadYaml } from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import ComplexTable from "@/components/complex-table"
import { checkPermissions } from "@/utils/permission"

export default {
  name: "StatefulSets",
  components: { LayoutContent, ComplexTable, KoTableOperations, KoFormItem },
  data() {
    return {
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "StatefulSetEdit",
              params: { operation: "edit", namespace: row.metadata.namespace, name: row.metadata.name },
              query: { yamlShow: false },
            })
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "apps", resource: "statefulsets", verb: "update" })
          },
        },
        {
          label: this.$t("commons.button.scale"),
          icon: "el-icon-copy-document",
          click: (row) => {
            this.onScale(row)
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "apps", resource: "statefulsets", verb: "update" })
          },
        },
        {
          label: this.$t("commons.button.modifying_version"),
          icon: "el-icon-edit-outline",
          click: (row) => {
            this.onModifyVersion(row)
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "apps", resource: "statefulsets", verb: "update" })
          },
        },
        {
          label: this.$t("commons.button.restart"),
          icon: "el-icon-refresh-right",
          click: (row) => {
            this.onRestart(row)
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "apps", resource: "statefulsets", verb: "update" })
          },
        },
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "StatefulSetEdit",
              params: { operation: "edit", namespace: row.metadata.namespace, name: row.metadata.name },
              query: { yamlShow: true },
            })
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "apps", resource: "statefulsets", verb: "update" })
          },
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml", row)
          },
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "apps", resource: "statefulsets", verb: "delete" })
          },
        },
      ],
      loading: false,
      data: [],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      dialogScaleVisible: false,
      form: {
        name: "",
        namespace: "",
        replicas: 1,
        imagesData: [],
      },
      dialogModifyVersionVisible: false,
      searchConfig: {
        keywords: "",
      },
      selects: [],
      clusterName: "",
    }
  },
  methods: {
    onCreate() {
      this.$router.push({ name: "StatefulSetCreate", params: { operation: "create" }, query: { yamlShow: false } })
    },
    openDetail(row) {
      this.$router.push({ name: "StatefulSetDetail", params: { namespace: row.metadata.namespace, name: row.metadata.name }, query: { yamlShow: false } })
    },
    yamlCreate() {
      this.$router.push({ name: "StatefulSetCreate", params: { operation: "create" }, query: { yamlShow: true } })
    },
    onDelete(row) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deleteWorkLoad(this.clusterName, "statefulsets", row.metadata.namespace, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteWorkLoad(this.clusterName, "statefulsets", select.metadata.namespace, select.metadata.name))
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
    onScale(row) {
      this.form.name = row.metadata.name
      this.form.namespace = row.metadata.namespace
      this.form.replicas = row.spec.replicas
      this.dialogScaleVisible = true
    },
    onScaleSubmit() {
      this.loading = true
      let data = { spec: { replicas: this.form.replicas } }
      scaleStatefulset(this.clusterName, this.form.namespace, this.form.name, data)
        .then(() => {
          this.dialogScaleVisible = false
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
    onModifyVersion(row) {
      this.form.imagesData = []
      this.form.name = row.metadata.name
      this.form.namespace = row.metadata.namespace
      if (row.spec.template.spec.initContainers) {
        for (const c of row.spec.template.spec.containers) {
          let index = c.image.lastIndexOf(":")
          this.form.imagesData.push({
            type: this.$t("business.workload.initContainer"),
            name: c.name,
            image: index !== -1 ? c.image.substring(0, index) : "",
            version: index !== -1 ? c.image.substring(index + 1, c.image.length) : "",
            newVersion: "",
          })
        }
      }
      for (const c of row.spec.template.spec.containers) {
        let index = c.image.lastIndexOf(":")
        this.form.imagesData.push({
          type: this.$t("business.workload.standardContainer"),
          name: c.name,
          image: index !== -1 ? c.image.substring(0, index) : "",
          version: index !== -1 ? c.image.substring(index + 1, c.image.length) : "",
          newVersion: "",
        })
      }
      this.dialogModifyVersionVisible = true
    },
    onModifyVersionSubmit() {
      this.loading = true
      let containers = []
      let initContainers = []
      for (const c of this.form.imagesData) {
        if (c.type === this.$t("business.workload.initContainer")) {
          initContainers.push({
            name: c.name,
            image: c.image + ":" + (c.newVersion !== "" ? c.newVersion : c.version),
          })
        } else {
          containers.push({
            name: c.name,
            image: c.image + ":" + (c.newVersion !== "" ? c.newVersion : c.version),
          })
        }
      }
      let data = { spec: { template: { spec: { containers: containers, initContainers: initContainers } } } }
      patchStatefulset(this.clusterName, this.form.namespace, this.form.name, data)
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
    onRestart(row) {
      let data = { spec: { template: { metadata: { annotations: { "kubectl.kubernetes.io/restartedAt": new Date() } } } } }
      patchStatefulset(this.clusterName, row.metadata.namespace, row.metadata.name, data)
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
    search(resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      listWorkLoads(this.clusterName, "statefulsets", true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize)
        .then((res) => {
          this.data = res.items
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
