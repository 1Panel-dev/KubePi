<template>
  <layout-content>
    <complex-table :data="data">
      <template #header>
        <el-button-group>
          <el-button type="primary" size="small" @click="onCreate">
            {{ $t("commons.button.create") }}
          </el-button>
        </el-button-group>
        <br/>
      </template>
      <el-table-column :label="$t('commons.table.name')" min-width="100" fix show-overflow-tooltip> 
        <template v-slot:default="{row}">
          <span class="span-link" @click="onOpenDetail(row)">{{ row.metadata.name }}</span>
        </template>
      </el-table-column>

      <el-table-column :label="$t('commons.table.description')" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ $t(row.metadata.annotations["description"]) }}
        </template>
      </el-table-column>

      <el-table-column :label="$t('commons.table.built_in')" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ $t('commons.bool.' + row.metadata.annotations["builtin"]) }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.age')" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ row.metadata.annotations["created-at"] | ageFormat }}
        </template>
      </el-table-column>
      <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')" fix/>
    </complex-table>

    <el-dialog
        :title="$t(formTitle)+$t('business.cluster.role')"
        :visible.sync="clusterRoleFormDialogOpened"
        width="70%"
        center z-index="10">

      <el-form :model="clusterRoleForm" label-position="left" label-width="144px">
        <el-form-item :label="$t('commons.table.name')">
          <el-input v-model="clusterRoleForm.name" style="width: 80%"></el-input>
        </el-form-item>

        <el-form-item :label="$t('business.cluster.rule')">
          <el-button @click="onRuleCreate"><i class="el-icon-plus "></i></el-button>
          <table border="1" cellspacing="0" style="width: 80%">
            <thead style="background-color: #1d3e4d">
            <tr>
              <th style="width: 30%">{{ $t('business.cluster.api_group') }}</th>
              <th style="width: 30%">{{ $t('business.cluster.resource') }}</th>
              <th style="width: 30%">{{ $t('business.cluster.verb') }}</th>
              <th>{{ $t('commons.table.action') }}</th>
            </tr>
            </thead>
            <tbody>
            <tr v-if="clusterRoleForm.rules.length===0">
              <td style="text-align: center" colspan="4">{{ $t('commons.msg.no_data') }}</td>
            </tr>
            <tr v-for="(item,index) in clusterRoleForm.rules" :key="index">
              <td style="text-align: center">
                <el-select v-model="item.apiGroups" filterable style="width: 100%" multiple>
                  <el-option v-for="(groupName,index) in getApiGroupOptions()"
                             :key="index"
                             :value="groupName">
                    {{ groupName }}
                  </el-option>
                </el-select>
              </td>
              <td>
                <el-select multiple v-model="item.resources" filterable style="width:100%">
                  <el-option v-for="(item,index) in getResourcesByApiGroupNames(item.apiGroups)"
                             :key="index"
                             :value="item">
                    {{ item }}
                  </el-option>
                </el-select>
              </td>
              <td>
                <el-select multiple v-model="item.verbs" filterable style="width:100%">
                  <el-option v-for="(item,index) in verbOptions"
                             :key="index"
                             :value="item">
                    {{ item }}
                  </el-option>
                </el-select>
              </td>
              <td>
                <el-button icon="el-icon-delete" size="mini" @click="onRuleDelete(item)" circle></el-button>
              </td>
            </tr>
            </tbody>
          </table>
        </el-form-item>
      </el-form>

      <span slot="footer" class="dialog-footer">
                <el-button @click="clusterRoleFormDialogOpened = false">{{ $t('commons.button.cancel') }}</el-button>
                <el-button type="primary" @click="onConfirm">{{ $t('commons.button.confirm') }}</el-button>
            </span>
    </el-dialog>

    <el-dialog :title="$t('business.cluster.role')"
               :visible.sync="clusterDetailDialogOpened"
               width="70%"
               center z-index="10">
      <el-form :model="detailForm" label-position="left" label-width="144px">
        <el-form-item :label="$t('commons.table.name')">
          {{ detailForm.name }}
        </el-form-item>

        <el-form-item :label="$t('business.cluster.rule')">
          <table border="1" cellspacing="0" style="width: 80%">
            <thead style="background-color: #1d3e4d">
            <tr>
              <th style="width: 30%">{{ $t('business.cluster.api_group') }}</th>
              <th style="width: 30%">{{ $t('business.cluster.resource') }}</th>
              <th style="width: 30%">{{ $t('business.cluster.verb') }}</th>
            </tr>
            </thead>
            <tbody>
            <tr v-if="detailForm.rules.length===0">
              <td style="text-align: center" colspan="4">{{ $t('commons.msg.no_data') }}</td>
            </tr>
            <tr v-for="(item,index) in detailForm.rules" :key="index">
              <td style="text-align: center">
                {{ item.apiGroups }}
              </td>
              <td style="text-align: center">
                {{ item.resources }}
              </td>
              <td style="text-align: center">
                {{ item.verbs }}
              </td>
            </tr>
            </tbody>
          </table>
        </el-form-item>
      </el-form>


    </el-dialog>


  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {
  listClusterRoles,
  createClusterRole,
  deleteClusterRole,
  updateClusterRole
} from "@/api/clusters";


export default {
  name: "ClusterRoles",
  props: ["name"],
  components: {LayoutContent, ComplexTable},
  data() {
    return {
      clusterRoleFormDialogOpened: false,
      clusterDetailDialogOpened: false,
      apiGroupResources: {
        "core": ["events", "services", "endpoints",
          "configmaps", "secrets", "resourcequotas", "limitranges", "persistentvolumeclaims",
          "pods", "containers", "serviceaccounts"
        ],
        "apps": ["deployments", "daemonsets", "replicasets", "statefulsets"],
        "batch": ["jobs", "cronjobs"],
        "autoscaling": ["horizontalpodautoscalers"],
        "networking.k8s.io": ["ingresses", "networkpolicies"],
        "rbac.authorization.k8s.io": ["role", "rolebinding"]
      },
      operation: "create",
      verbOptions: ["*", "create", "delete", "deletecollection", "get", "list", "patch", "update", "watch"],
      clusterRoleForm: {
        name: "",
        rules: [],
      },
      detailForm: {
        name: "",
        rules: []
      },
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.onEdit(row)
          },
          disabled: (row) => {
            return row.metadata.annotations['builtin'] === 'true'
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
          disabled: (row) => {
            return row.metadata.annotations['builtin'] === 'true'
          }
        },
      ],
      data: [],
    }
  },
  computed: {
    formTitle() {
      return `commons.button.${this.operation}`
    }

  },
  methods: {
    list() {
      this.loading = false
      listClusterRoles(this.name, "namespace").then(data => {
        this.loading = false
        this.data = data.data
      })
    },

    onOpenDetail(item) {
      this.detailForm = {
        name: item.metadata.name,
        rules: item.rules
      }
      for (const rule of this.detailForm.rules) {
        for (let i = 0; i < rule.apiGroups.length; i++) {
          console.log(rule)
          if (rule.apiGroups[i] === "") {
            rule.apiGroups[i] = "core"
          }
        }
      }
      this.clusterDetailDialogOpened = true
    },

    onCreate() {
      this.operation = "create"
      this.clusterRoleForm = {
        name: "",
        rules: [],
      }
      this.clusterRoleFormDialogOpened = true
    },

    getApiGroupOptions() {
      return ["*"].concat(Object.keys(this.apiGroupResources))
    },
    getResourcesByApiGroupNames(groupNames) {
      let res = []
      for (const gn of groupNames) {
        if (this.apiGroupResources[gn]) {
          res = res.concat(this.apiGroupResources[gn])
        }
      }
      return ["*"].concat(res)
    },
    onRuleDelete(index) {
      this.clusterRoleForm.rules.splice(index, 1)
    },
    onRuleCreate() {
      for (const nr of this.clusterRoleForm.rules) {
        if (nr.apiGroups.length === 0 || nr.resources.length === 0 || nr.verbs.length === 0) {
          this.$message.error(this.$t('business.cluster.namespace_role_form_check_msg'))
          return
        }
      }
      this.clusterRoleForm.rules.push({
        apiGroups: [],
        resources: [],
        verbs: []
      })
    },
    onEdit(row) {
      this.operation = "edit"
      this.clusterRoleForm.name = row.metadata.name
      this.clusterRoleForm.rules = []
      for (const rule of row.rules) {
        const r = {
          apiGroups: [],
          resources: rule.resources,
          verbs: rule.verbs,
        }
        for (const g of rule.apiGroups) {
          if (g === "") {
            r.apiGroups.push("core")
          } else {
            r.apiGroups.push(g)
          }
        }
        this.clusterRoleForm.rules.push(r)
      }
      this.clusterRoleFormDialogOpened = true
    }
    ,
    onDelete(row) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: 'warning'
      }).then(() => {
        deleteClusterRole(this.name, row.metadata.name).then(() => {
          this.list()
        })
      });
    }
    ,
    onConfirm() {
      const req = {
        metadata: {
          name: this.clusterRoleForm.name,
          annotations: {
            "description": this.clusterRoleForm.description
          },
          labels: {
            "kubepi.org/role-type": "namespace",
          }
        },
        rules: []
      }
      for (const rule of this.clusterRoleForm.rules) {
        req.rules.push({
          apiGroups: rule.apiGroups,
          resources: rule.resources,
          verbs: rule.verbs,
        })
      }
      switch (this.operation) {
        case "create":
          createClusterRole(this.name, req).then(() => {
            this.list()
            this.clusterRoleFormDialogOpened = false
            this.$message.success(this.$t('commons.msg.create_success'))
          })
          break
        case "edit":
          updateClusterRole(this.name, this.clusterRoleForm.name, req).then(() => {
            this.list()
            this.clusterRoleFormDialogOpened = false
            this.$message.success(this.$t('commons.msg.update_success'))
          })
          break
      }
    }
    ,
  },
  created() {
    this.list()
  }
}
</script>

<style scoped>

</style>