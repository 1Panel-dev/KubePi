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
      <el-table-column :label="$t('commons.table.name')" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ row.metadata.name }}
        </template>
      </el-table-column>

      <el-table-column :label="$t('business.cluster.scope')" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ $t('business.cluster.' + row.metadata.labels["kubeoperator.io/role-type"]) }}
        </template>
      </el-table-column>

      <el-table-column :label="$t('commons.table.built_in')" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ $t('commons.bool.' + row.metadata.annotations["builtin"]) }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ row.metadata.annotations["created-at"] }}
        </template>
      </el-table-column>
      <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')" fix/>
    </complex-table>

    <el-dialog
        :title="$t('commons.button.create')+$t('business.cluster.rule')"
        :visible.sync="ruleDialogOpened"
        width="40%"
        center z-index="20">
      <el-form :model="ruleForm" label-position="left" label- width="144px">
        <el-form-item label="API Group">
          <el-select v-model="ruleForm.apiGroup" style="width:100%" @change="onAPIGroupChange">
            <el-option v-for="(item,index) in apiGroupsOptions" :key="index"
                       :value="item.name">
                            <span v-if="item.name">
                                {{ item.name }}
                            </span>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="API Resource">
          <el-select multiple v-model="ruleForm.resources" style="width:100%" :disabled="resourcesDisable"
                     @change="onResourcesChange">
            <el-option v-for="(item,index) in apiResourceOptions"
                       :key="index"
                       :value="item">
              {{ item }}
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="Verbs">
          <el-select multiple v-model="ruleForm.verbs" style="width:100%" :disabled="verbsDisable"
                     @change="onVerbsChange">
            <el-option v-for="(item,index) in verbOptions"
                       :key="index"
                       :value="item">
              {{ item }}
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>


      <span slot="footer" class="dialog-footer">
                <el-button @click="ruleDialogOpened=false">{{ $t('commons.button.cancel') }}</el-button>
                <el-button type="primary" @click="onRuleConfirm">{{ $t('commons.button.confirm') }}</el-button>
            </span>
    </el-dialog>

    <el-dialog
        :title="$t('commons.button.create')+$t('business.cluster.role')"
        :visible.sync="createDialogOpened"
        width="50%"
        center z-index="10">

      <el-form :model="clusterRoleForm" label-position="left" label-width="144px">
        <el-form-item :label="$t('commons.table.name')">
          <el-input v-model="clusterRoleForm.name" style="width: 80%"></el-input>
        </el-form-item>
        <el-form-item :label="$t('business.cluster.scope')">
          <el-radio v-model="clusterRoleForm.type" label="cluster">{{ $t('business.cluster.cluster') }}</el-radio>
          <el-radio v-model="clusterRoleForm.type" label="namespace">{{ $t('business.cluster.namespace') }}</el-radio>
        </el-form-item>

        <el-form-item :label="$t('business.cluster.rule')">
          <el-button @click="onRuleCreate"><i class="el-icon-plus "></i></el-button>
          <el-table
              :data="clusterRoleForm.rules"
              border
              style="width: 80%">
            <el-table-column
                prop="apiGroup"
                label="API Group"
            >
            </el-table-column>
            <el-table-column
                prop="resource"
                label="API Resource"
            >
              <template v-slot:default="{row}">
                <span v-for="v in row.resources" :key="v">{{ v }}<br/></span>
              </template>
            </el-table-column>
            <el-table-column
                prop="verbs"
                label="Verbs"
            >
              <template v-slot:default="{row}">
                <span v-for="v in row.verbs" :key="v">{{ v }}<br/></span>
              </template>
            </el-table-column>
            <el-table-column>
              <template v-slot:default="{row}">
                <el-button
                    size="mini" circle @click="onRuleDelete(row)">
                  <i class="el-icon-delete"></i>
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
      </el-form>


      <span slot="footer" class="dialog-footer">
                <el-button @click="createDialogOpened = false">{{ $t('commons.button.cancel') }}</el-button>
                <el-button type="primary" @click="onConfirm">{{ $t('commons.button.confirm') }}</el-button>
            </span>
    </el-dialog>


    <el-dialog
        :title="$t('commons.button.edit')+$t('business.cluster.role')"
        :visible.sync="editDialogOpened"
        width="50%"
        center z-index="10">

      <el-form :model="editForm" label-position="left" label-width="144px">
        <el-form-item :label="$t('commons.table.name')">
          <el-input v-model="editForm.name" style="width: 80%" disabled></el-input>
        </el-form-item>
        <el-form-item :label="$t('business.cluster.scope')" disabled>
          <el-radio v-model="editForm.type" label="cluster">{{ $t('business.cluster.cluster') }}</el-radio>
          <el-radio v-model="editForm.type" label="namespace">{{ $t('business.cluster.namespace') }}</el-radio>
        </el-form-item>

        <el-form-item :label="$t('business.cluster.rule')">
          <el-button @click="onRuleCreate"><i class="el-icon-plus "></i></el-button>
          <el-table
              :data="editForm.rules"
              border
              style="width: 80%">
            <el-table-column
                prop="apiGroup"
                label="API Group"
            >
            </el-table-column>
            <el-table-column
                prop="resource"
                label="API Resource"
            >
              <template v-slot:default="{row}">
                <span v-for="v in row.resources" :key="v">{{ v }}<br/></span>
              </template>
            </el-table-column>
            <el-table-column
                prop="verbs"
                label="Verbs"
            >
              <template v-slot:default="{row}">
                <span v-for="v in row.verbs" :key="v">{{ v }}<br/></span>
              </template>
            </el-table-column>
            <el-table-column>
              <template v-slot:default="{row}">
                <el-button
                    size="mini" circle @click="onRuleDelete(row)">
                  <i class="el-icon-delete"></i>
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
      </el-form>


      <span slot="footer" class="dialog-footer">
                <el-button @click="editDialogOpened = false">{{ $t('commons.button.cancel') }}</el-button>
                <el-button type="primary" @click="onConfirm">{{ $t('commons.button.confirm') }}</el-button>
            </span>
    </el-dialog>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {
  listClusterRoles,
  listClusterApiGroups,
  listClusterResourceByGroupVersion,
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
      createDialogOpened: false,
      editDialogOpened: false,
      ruleDialogOpened: false,
      apiGroupsOptions: [],
      apiResourceOptions: [],
      operation: "",
      verbOptions: ["*", "create", "delete", "deletecollection", "get", "list", "patch", "update", "watch"],
      clusterRoleForm: {
        name: "",
        type: "cluster",
        rules: [],
      },
      editForm: {
        name: "",
        rules: [],
      },
      resourcesDisable: false,
      verbsDisable: false,
      ruleForm: {
        apiGroup: "",
        resources: [],
        verbs: []
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
  methods: {
    list() {
      this.loading = false
      listClusterRoles(this.name).then(data => {
        this.loading = false
        this.data = data.data
      })
    },
    onCreate() {
      this.operation = "create"
      this.clusterRoleForm = {
        name: "",
        rules: [],
        type: "cluster",
      }
      this.createDialogOpened = true
      listClusterApiGroups(this.name).then(data => {
        this.apiGroupsOptions.push({name: "*"})
        this.apiGroupsOptions = this.apiGroupsOptions.concat(data.data)
        for (const group of this.apiGroupsOptions) {
          if (group.name === "") {
            group.name = "core"
          }
        }
      })
    },

    onEdit(row) {
      this.operation = "update"
      this.editForm.name = row.metadata.name
      this.editForm.rules = []
      this.editForm.type = row.metadata.labels["kubeoperator.io/role-type"]
      for (const rule of row.rules) {
        if (rule.apiGroups) {
          const r = {
            apiGroup: [],
            resources: rule.resources,
            verbs: rule.verbs,
          }
          for (const g of rule.apiGroups) {
            if (g === "") {
              r.apiGroup = "core"
            } else {
              r.apiGroup = g
            }
          }
          this.editForm.rules.push(r)
        }
      }
      this.editDialogOpened = true
      listClusterApiGroups(this.name).then(data => {
        this.apiGroupsOptions.push({name: "*"})
        this.apiGroupsOptions = this.apiGroupsOptions.concat(data.data)
        for (const group of this.apiGroupsOptions) {
          if (group.name === "") {
            group.name = "core"
          }
        }
      })
    },
    onRuleCreate() {
      this.ruleForm = {
        apiGroup: "",
        resources: [],
        verbs: [],
      }
      this.ruleDialogOpened = true
      this.resourcesDisable = false
      this.verbsDisable = false
    },
    onAPIGroupChange() {
      this.apiResourceOptions = []
      this.ruleForm.resources = []
      this.ruleForm.verbs = []
      this.resourcesDisable = false
      this.verbsDisable = false
      if (this.ruleForm.apiGroup === "*") {
        this.apiResourceOptions = ["*"]
        this.ruleForm.resources = ["*"]
        this.resourcesDisable = true
        return
      }
      listClusterResourceByGroupVersion(this.name, this.ruleForm.apiGroup).then(data => {
        this.apiResourceOptions.push("*")
        this.apiResourceOptions = this.apiResourceOptions.concat(data.data);
      })
    },
    onResourcesChange() {
      if (this.ruleForm.resources.indexOf("*") > -1) {
        this.ruleForm.resources = ["*"]
      }
    },
    onVerbsChange() {
      if (this.ruleForm.verbs.indexOf("*") > -1) {
        this.ruleForm.verbs = ["*"]
      }
    },
    onRuleConfirm() {
      const rule = {
        apiGroup: this.ruleForm.apiGroup,
        resources: this.ruleForm.resources,
        verbs: this.ruleForm.verbs
      }
      switch (this.operation) {

        case "create":
          this.clusterRoleForm.rules.push(rule)
          break
        case "update":
          this.editForm.rules.push(rule)
          break
      }
      this.ruleDialogOpened = false
    },

    onRuleDelete(row) {
      switch (this.operation) {
        case "create":
          this.clusterRoleForm.rules.splice(this.clusterRoleForm.rules.indexOf(row), 1)
          break
        case "update":
          this.editForm.rules.splice(this.clusterRoleForm.rules.indexOf(row), 1)
          break
      }
    },
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
    },

    onConfirm() {
      const req = {
        metadata: {
          name: this.clusterRoleForm.name,
          labels: {
            "kubeoperator.io/role-type": this.clusterRoleForm.type
          }
        },
        rules: []
      }
      for (const rule of this.clusterRoleForm.rules) {
        req.rules.push({
          apiGroups: [rule.apiGroup],
          resources: rule.resources,
          verbs: rule.verbs,
        })
      }
      createClusterRole(this.name, req).then(() => {
        this.list()
        this.createDialogOpened = false
      })
    },
    onEditConfirm() {
      const req = {
        metadata: {
          name: this.editForm.name,
        },
        rules: []
      }
      for (const rule of this.editForm.rules) {
        req.rules.push({
          apiGroups: [rule.apiGroup],
          resources: rule.resources,
          verbs: rule.verbs,
        })
      }
      updateClusterRole(this.name, req.metadata.name, req).then(() => {
        this.list()
        this.editDialogOpened = false
      })
    }
  },
  created() {
    this.list()
  }
}
</script>

<style scoped>

</style>