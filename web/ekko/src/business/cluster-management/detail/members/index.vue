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
          {{ row.name }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ row.createAt | datetimeFormat }}
        </template>
      </el-table-column>
      <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')" fix/>
    </complex-table>


    <el-dialog
        :title="$t('commons.button.'+operation)+$t('business.cluster.member')"
        :visible.sync="formDialogOpened"
        z-index="10"
        width="60%"
        center>
      <el-form :model="memberForm" label-position="left" label-width="144px">
        <el-form-item :label="$t('business.user.user')+$t('commons.table.name')">
          <el-select v-model="memberForm.userName" style="width: 80%" :disabled="operation==='edit'">
            <el-option v-for="(item, index) in getUserOptions" :key="index" :value="item.name">
              {{ item.name }}
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('business.cluster.cluster')+$t('business.cluster.role')">
          <el-radio-group v-model="memberForm.roleType" @change="onRoleTypeChange">
            <el-radio label="admin">{{ $t('business.cluster.administrator') }}</el-radio>
            <el-radio label="viewer">{{ $t('business.cluster.viewer') }}</el-radio>
            <el-radio label="custom">{{ $t('business.cluster.custom') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <div v-if="memberForm.roleType==='custom'">
          <el-form-item>
            <el-select v-model="memberForm.customClusterRoles" multiple style="width: 80%">
              <el-option
                  v-for="(item,index) in getClusterRolesOptions"
                  :key="index"
                  :value="item.metadata.name">
                {{ item.metadata.name }}
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item :label="$t('business.cluster.namespace')+$t('business.cluster.role')">
            <el-button @click="onNamespaceRoleCreate"><i class="el-icon-plus "></i></el-button>
            <table border="1" cellspacing="0" style="width: 80%">
              <thead style="background-color: #1d3e4d">
              <tr>
                <th style="width: 45%">{{ $t('business.cluster.namespace') }}</th>
                <th style="width: 45%">{{ $t('business.cluster.role') }}</th>
                <th>{{ $t('commons.table.action') }}</th>
              </tr>
              </thead>
              <tbody>
              <tr v-if="memberForm.namespaceRoles.length===0">
                <td style="text-align: center" colspan="3">{{ $t('commons.msg.no_data') }}</td>
              </tr>


              <tr v-for="(item, index) in memberForm.namespaceRoles" :key="index">
                <td style="text-align: center">
                  <el-select v-model="item.namespace" style="width:100%">
                    <el-option v-for="(item,index) in getNamespaceOptions"
                               :key="index"
                               :value="item.metadata.name">
                      {{ item.metadata.name }}
                    </el-option>
                  </el-select>
                </td>
                <td>
                  <el-select multiple v-model="item.roles" style="width:100%">
                    <el-option v-for="(item,index) in namespaceRoleOptions"
                               :key="index"
                               :value="item.metadata.name">
                      {{ item.metadata.name }}
                    </el-option>
                  </el-select>
                </td>
                <td style="text-align: center">
                  <el-button @click="onNamespaceRoleDelete(index)" icon="el-icon-delete" size="mini" circle></el-button>
                </td>
              </tr>
              </tbody>
            </table>

          </el-form-item>
        </div>
      </el-form>
      <span slot="footer" class="dialog-footer">
                <el-button @click="formDialogOpened = false">{{ $t('commons.button.cancel') }}</el-button>
                <el-button type="primary" @click="onConfirm">{{ $t('commons.button.confirm') }}</el-button>
            </span>
    </el-dialog>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {createClusterMember, listClusterMembers, listNamespaces} from "@/api/clusters"
import {listUsers} from "@/api/users"
import {listClusterRoles, deleteClusterMember, getClusterMember, updateClusterMember} from "@/api/clusters";


export default {
  name: "ClusterMembers",
  props: ["name"],
  components: {LayoutContent, ComplexTable},
  data() {
    return {

      formDialogOpened: false,
      editDialogOpened: false,
      operation: "create",
      namespaceRoleDialogOpened: false,
      userOptions: [],
      clusterRolesOptions: [],
      namespaceRoleOptions: [],
      namespaceOptions: [],
      memberForm: {
        userName: "",
        customClusterRoles: [],
        namespaceRoles: [],
        roleType: "admin"
      },
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.onEdit(row)
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          }
        },
      ],
      data: [],
    }
  },
  computed: {
    getUserOptions() {
      return this.userOptions.filter((n) => {
        let exists = false
        for (const user of this.data) {
          if (user.name === n.name) {
            exists = true
          }
        }
        return !exists
      })
    },
    getNamespaceOptions() {
      return this.namespaceOptions.filter((n) => {
        let exists = false
        for (const nrs of this.memberForm.namespaceRoles) {
          if (nrs.namespace === n.metadata.name) {
            exists = true
          }
        }
        return !exists
      })
    },
    getClusterRolesOptions() {
      return this.clusterRolesOptions.filter((cr) => {
        if (this.memberForm.roleType === 'custom') {
          return !["admin-cluster", "view-cluster"].includes(cr.metadata.name)
        }
        return true
      })
    }
  },
  methods: {
    list() {
      this.loading = false
      listClusterMembers(this.name).then(data => {
        this.loading = false
        this.data = data.data
      })
    },
    listClusterRoles() {
      listClusterRoles(this.name).then(data => {
        this.clusterRolesOptions = data.data.filter((r) => {
          return r.metadata["labels"]["kubeoperator.io/role-type"] === "cluster"
        })
        this.namespaceRoleOptions = data.data.filter((r) => {
          return r.metadata["labels"]["kubeoperator.io/role-type"] === "namespace"
        })
      })
    },
    onCreate() {
      this.formDialogOpened = true
      this.operation = "create"
      this.memberForm.userName = ""
      this.memberForm.customClusterRoles = []
      this.memberForm.roleType = "admin"
      this.memberForm.namespaceRoles = []
      this.listClusterRoles()
      listUsers().then(data => {
        this.userOptions = data.data;
      })
      listNamespaces(this.name).then(data => {
        this.namespaceOptions = data.data;
      })
    },
    onEdit(row) {
      this.operation = "edit"
      this.memberForm.userName = ""
      this.memberForm.customClusterRoles = []
      this.memberForm.roleType = "admin"
      this.memberForm.namespaceRoles = []
      this.listClusterRoles()
      listNamespaces(this.name).then(data => {
        this.namespaceOptions = data.data;
      })
      getClusterMember(this.name, row.name).then(data => {
        this.memberForm.userName = data.data.name
        console.log(data.data)
        if (data.data.clusterRoles && data.data.clusterRoles.length === 1 && data.data.namespaceRoles.length === 0) {
          if (data.data.clusterRoles[0] === 'admin-cluster') {
            this.memberForm.roleType = 'admin'
          } else if (data.data.clusterRoles[0] === 'view-cluster') {
            this.memberForm.roleType = 'viewer'
          }
        } else {
          this.memberForm.roleType = 'custom'
          this.memberForm.customClusterRoles = data.data.clusterRoles
          this.memberForm.namespaceRoles = data.data.namespaceRoles
        }
      })
      this.formDialogOpened = true
    },
    onRoleTypeChange() {
      if (this.memberForm.roleType === 'admin' || this.memberForm.roleType === 'viewer') {
        this.memberForm.customClusterRoles = []
        this.memberForm.namespaceRoles = []
      }
    },
    onDelete(raw) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: 'warning'
      }).then(() => {
        deleteClusterMember(this.name, raw.name, raw.kind).then(() => {
          this.$message.success(this.$t('commons.msg.delete_success'))
          this.list()
        })
      });
    },
    onNamespaceRoleCreate() {
      for (const nr of this.memberForm.namespaceRoles) {
        if (!nr.namespace || nr.roles.length === 0) {
          this.$message.error(this.$t('business.cluster.namespace_role_form_check_msg'))
          return
        }
      }
      this.memberForm.namespaceRoles.push({
        namespace: "",
        roles: [],
      })
    },
    onNamespaceRoleDelete(index) {
      this.memberForm.namespaceRoles.splice(index, 1)
    },
    onConfirm() {
      let req = {
        name: this.memberForm.userName,
        namespaceRoles: this.memberForm.namespaceRoles,
      }
      switch (this.memberForm.roleType) {
        case "admin":
          req.clusterRoles = ["admin-cluster"]
          break
        case "viewer":
          req.clusterRoles = ["view-cluster"]
          break
        case "custom":
          req.clusterRoles = this.memberForm.customClusterRoles
          break
      }
      switch (this.operation) {
        case "create":
          createClusterMember(this.name, req).then(() => {
            this.formDialogOpened = false
            this.list()
          })
          break;
        case "edit":
          console.log(req)
          updateClusterMember(this.name, this.memberForm.userName, req).then(() => {
            this.formDialogOpened = false
            this.list()
          })
      }
    },
  }
  ,
  created() {
    this.list()
  }
}
</script>

<style scoped>

</style>