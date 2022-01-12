<template>
  <layout-content>
    <div>
      <el-row>
        <el-col :span="6">
          <el-scrollbar style="height:900px">
            <el-tree
                    v-loading="loading"
                    :empty-text="$t('commons.table.empty_text')"
                    class="filter-tree"
                    :default-expand-all="true"
                    :data="data"
                    :props="props"
                    node-key="id"
                    ref="tree"
                    :highlight-current="true"
                    @node-click="handleNodeClick"
                    style="margin: 5px">
            </el-tree>
          </el-scrollbar>
        </el-col>
        <el-col :span="18">
            <c-r-d-list v-if="resourceType==='Definitions'"></c-r-d-list>
            <c-r-list :key="key" v-if="resourceType==='CustomResources'" :version="crItem.version" :names="crItem.names" :group="crItem.group" :scope="crItem.scope"></c-r-list>
        </el-col>
      </el-row>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {downloadYaml} from "@/utils/actions"
import {deleteCustomResourceDefinition, getCustomResourceDefinition, listCustomResourceDefinitions,} from "@/api/customresourcedefinitions"
import {checkPermissions} from "@/utils/permission"
import CRDList from "@/business/custom-resource/crd"
import CRList from "@/business/custom-resource/cr/cr"

export default {
  name: "CustomResourceDefinitions",
  components: { CRList, CRDList, LayoutContent },
  data () {
    return {
      data: [],
      selects: [],
      loading: false,
      cluster: "",
      props: {
        children: "children",
        label: "label"
      },
      resources:[],
      resourceType: "Definitions",
      crItem:{},
      key:0,
      buttons: [
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "CustomResourceDefinitionEdit",
              params: { name: row.metadata.name }
            })
          },
          disabled: () => {
            return !checkPermissions({
              scope: "namespace",
              apiGroup: "apiextensions.k8s.io",
              resource: "customresourcedefinitions",
              verb: "update"
            })
          }
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml", getCustomResourceDefinition(this.cluster, row.metadata.name))
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
          disabled: () => {
            return !checkPermissions({
              scope: "namespace",
              apiGroup: "apiextensions.k8s.io",
              resource: "customresourcedefinitions",
              verb: "delete"
            })
          }
        },
      ],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: ""
      }
    }
  },
  methods: {
    search () {
      this.loading = true
      listCustomResourceDefinitions(this.cluster).then(res => {
        this.resources = res.items
        let groups = new Map()
        for (const item of res.items) {
          const group = item.spec.group
          let children = groups.get(group)
          if (children === undefined) {
            groups.set(group, [item])
          } else {
            children.push(item)
            groups.set(group, children)
          }
        }
        this.data.push({label:"Definitions"})
        for (let [key, value] of groups.entries()) {
          let array = []
          for (let children of value) {
            array.push({ label: children.spec.names.kind, group: key,name:children.spec.names.plural,version: children.spec.versions[0].name,scope:children.spec.scope})
          }
          this.data.push({
            label: key,
            children: array
          })
        }

      }).finally(() => {
        this.loading = false
      })
    },
    handleNodeClick (data) {
      if (data.label === "Definitions") {
        this.resourceType = "Definitions"
      }else if (data.children === undefined) {
        this.key ++
        this.resourceType = "CustomResources"
        this.crItem = {
          version: data.version,
          group: data.group,
          names:data.name,
          scope: data.scope,
        }
      }
    },
    onDelete (row) {
      this.$confirm(
        this.$t("commons.confirm_message.delete"),
        this.$t("commons.message_box.prompt"), {
          confirmButtonText: this.$t("commons.button.confirm"),
          cancelButtonText: this.$t("commons.button.cancel"),
          type: "warning",
        }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deleteCustomResourceDefinition(this.cluster, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteCustomResourceDefinition(this.cluster, select.metadata.name))
            }
          }
        }
        if (this.ps.length !== 0) {
          Promise.all(this.ps)
            .then(() => {
              this.search()
              this.$message({
                type: "success",
                message: this.$t("commons.msg.delete_success"),
              })
            })
            .catch(() => {
              this.search()
            })
        }
      })
    },
    openDetail (row) {
      this.$router.push({
        name: "CustomResourceDefinitionDetail",
        params: { name: row.metadata.name },
        query: { yamlShow: false }
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.search()
  }
}
</script>

<style scoped>

</style>
