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
          {{ row.createAt }}
        </template>
      </el-table-column>
      <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')" fix/>
    </complex-table>


    <el-dialog
        title="创建命名空间角色"
        :visible.sync="namespaceRoleDialogOpened"
        width="40%"
        center z-index="20">
      <el-form label-position="left" label-width="144px" :model="namespaceRoleForm">
        <el-form-item label="命名空间">
          <el-select v-model="namespaceRoleForm.namespace" style="width:100%">
            <el-option v-for="(item,index) in getNamespaceOptions"
                       :key="index"
                       :value="item.metadata.name">
              {{ item.metadata.name }}
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="角色">
          <el-select multiple v-model="namespaceRoleForm.roles" style="width:100%">
            <el-option v-for="(item,index) in namespaceRoleOptions"
                       :key="index"
                       :value="item.metadata.name">
              {{ item.metadata.name }}
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
                <el-button @click="namespaceRoleDialogOpened=false">取 消</el-button>
                <el-button type="primary" @click="onNamespaceRoleConfirm">确 定</el-button>
        </span>
    </el-dialog>

    <el-dialog
        title="添加成员"
        :visible.sync="createDialogOpened"
        z-index="10"
        width="60%"
        center>
      <el-form :model="memberForm" label-position="left" label-width="144px">
        <el-form-item label="用户名称">
          <el-select v-model="memberForm.userName" style="width: 80%">
            <el-option v-for="(item, index) in userOptions" :key="index" :value="item.name">
              {{ item.name }}
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="集群角色">
          <el-radio-group v-model="memberForm.roleType">
            <el-radio label="admin">管理者</el-radio>
            <el-radio label="viewer">只读者</el-radio>
            <el-radio label="custom">自定义</el-radio>
          </el-radio-group>
        </el-form-item>
        <div v-if="memberForm.roleType==='custom'">
          <el-form-item>
            <el-select v-model="memberForm.customClusterRoles" multiple style="width: 80%"
                       placeholder="请选择">
              <el-option
                  v-for="(item,index) in clusterRolesOptions"
                  :key="index"
                  :value="item.metadata.name">
                {{ item.metadata.name }}
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="命名空间角色">
            <el-button @click="onNamespaceRoleCreate"><i class="el-icon-plus "></i></el-button>
            <el-table
                :data="memberForm.namespaceRoles"
                border
                style="width: 80%">
              <el-table-column
                  prop="namespace"
                  label="命名空间"
              >
              </el-table-column>
              <el-table-column
                  label="角色"
              >
                <template v-slot:default="{row}">
                  <span v-for="v in row.roles" :key="v">{{ v }}<br/></span>
                </template>
              </el-table-column>
              <el-table-column width="256x">
                <template v-slot:default="{row}">
                  <el-button
                      size="mini" circle @click="onNamespaceRoleDelete(row)">
                    <i class="el-icon-delete"></i>
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-form-item>
        </div>
      </el-form>
      <span slot="footer" class="dialog-footer">
                <el-button @click="createDialogOpened = false">取 消</el-button>
                <el-button type="primary" @click="onConfirm">确 定</el-button>
            </span>
    </el-dialog>

    <el-dialog
        title="编辑成员"
        :visible.sync="editDialogOpened"
        width="60%"
        center>
      <el-form :model="editForm" label-position="left" label-width="144px">
        <el-form-item label="用户名称">
          <el-select v-model="editForm.userName" style="width: 80%" disabled="disabled">
            <el-option v-for="(item, index) in userOptions" :key="index" :value="item.name">
              {{ item.name }}
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="集群角色">
          <el-radio-group v-model="editForm.roleType">
            <el-radio label="admin">管理者</el-radio>
            <el-radio label="viewer">只读者</el-radio>
            <el-radio label="custom">自定义</el-radio>
          </el-radio-group>
        </el-form-item>
        <div v-if="editForm.roleType==='custom'">
          <el-form-item>
            <el-select v-model="editForm.customClusterRoles" multiple style="width: 80%"
                       placeholder="请选择">
              <el-option
                  v-for="(item,index) in clusterRolesOptions"
                  :key="index"
                  :value="item.metadata.name">
                {{ item.metadata.name }}
              </el-option>
            </el-select>
          </el-form-item>
        </div>
      </el-form>
      <span slot="footer" class="dialog-footer">
                <el-button @click="editDialogOpened = false">取 消</el-button>
                <el-button type="primary" @click="onEditConfirm">确 定</el-button>
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
      createDialogOpened: false,
      editDialogOpened: false,
      namespaceRoleDialogOpened: false,
      userOptions: [],
      clusterRolesOptions: [],
      namespaceRoleOptions: [],
      namespaceOptions: [],
      namespaceRoleForm: {
        namespace: "",
        roles: [],
      },
      memberForm: {
        userName: "",
        customClusterRoles: [],
        namespaceRoles: [],
        roleType: "admin"
      },
      editForm: {
        userName: "",
        customClusterRoles: [],
        roleType: ""
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
    onCreate() {
      this.createDialogOpened = true
      listClusterRoles(this.name).then(data => {
        this.clusterRolesOptions = data.data.filter((r) => {
          return r.metadata["labels"]["kubeoperator.io/role-type"] === "cluster"
        })
        this.namespaceRoleOptions = data.data.filter((r) => {
          return r.metadata["labels"]["kubeoperator.io/role-type"] === "namespace"
        })
      })

      listUsers().then(data => {
        this.userOptions = data.data;
      })
      listNamespaces(this.name).then(data => {
        this.namespaceOptions = data.data;
      })
    },

    onEdit(row) {
      listClusterRoles(this.name).then(data => {
        this.clusterRolesOptions = data.data
        getClusterMember(this.name, row.name).then(data => {
          this.editForm.userName = data.data.name
          if (data.data.clusterRoles.length === 1) {
            if (data.data.clusterRoles[0] === 'Admin Cluster') {
              this.editForm.roleType = 'admin'
            } else if (data.data.clusterRoles[0] === 'View Cluster') {
              this.editForm.roleType = 'viewer'
            } else {
              this.editForm.roleType = 'custom'
              this.editForm.customClusterRoles = data.data.clusterRoles
            }
          } else {
            this.editForm.roleType = 'custom'
            this.editForm.customClusterRoles = data.data.clusterRoles
          }
        })
      })
      this.editDialogOpened = true
    },
    onDelete(raw) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: 'warning'
      }).then(() => {
        deleteClusterMember(this.name, raw.name, raw.kind).then(() => {
          this.$message.success("删除成功")
          this.list()
        })
      });
    },
    onNamespaceRoleCreate() {
      this.namespaceRoleForm.namespaceRoles = []
      this.namespaceRoleDialogOpened = true
      this.namespaceRoleForm.namespace = ""
    },
    onEditConfirm() {
      let req = {
        name: this.editForm.userName,
      }
      switch (this.editForm.roleType) {
        case "admin":
          req.clusterRoles = ["Admin Cluster"]
          break
        case "viewer":
          req.clusterRoles = ["View Cluster"]
          break
        case "custom":
          req.clusterRoles = this.editForm.customClusterRoles
          break
      }
      updateClusterMember(this.name, this.editForm.userName, req).then(() => {
        this.editDialogOpened = false
        this.list()
      })
    },
    onNamespaceRoleConfirm() {
      this.memberForm.namespaceRoles.push({
        namespace: this.namespaceRoleForm.namespace,
        roles: this.namespaceRoleForm.roles,
      })
      this.namespaceRoleDialogOpened = false
    },
    onNamespaceRoleDelete(row) {
      this.memberForm.namespaceRoles.splice(this.memberForm.namespaceRoles.indexOf(row), 1)

    },
    onConfirm() {
      let req = {
        name: this.memberForm.userName,
        namespaceRoles: this.memberForm.namespaceRoles,
      }
      switch (this.memberForm.roleType) {
        case "admin":
          req.clusterRoles = ["Admin Cluster"]
          break
        case "viewer":
          req.clusterRoles = ["View Cluster"]
          break
        case "custom":
          req.clusterRoles = this.memberForm.customClusterRoles
          break
      }
      createClusterMember(this.name, req).then(() => {
        this.createDialogOpened = false
        this.list()
      })
    },
  },
  created() {
    this.list()
  }
}
</script>

<style scoped>

</style>