<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.access_control.role')">
      <el-row :gutter="20">
        <el-form label-position="top" ref="form" :model="roleRef">
          <el-col :span="6">
            <el-form-item :label="$t('business.configuration.kind')" required>
              <el-select v-model="roleRef.kind" style="width: 100%" @change="changeKind(roleRef.kind)">
                <el-option label="ClusterRole" :value="'ClusterRole'"></el-option>
                <el-option label="Role" v-if="namespace" :value="'Role'"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item :label="$t('business.configuration.name')" required>
              <el-select v-model="roleRef.name" style="width: 100%">
                <el-option v-for="(role,index) in roles"
                           :key="index" :label="role.metadata.name" :value="role.metadata.name"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-form>
      </el-row>
    </ko-card>
    <div style="margin-top: 10px">
      <ko-card :title="$t('business.access_control.object')">
        <el-form label-position="top" ref="form" :model="form">
          <el-card v-for="(row,ind) in form.subjects" v-bind:key="ind"
                   style="background-color: #292a2e;margin-top: 10px;position: relative">
            <div style="float: right; padding: 3px 0;position: relative;z-index: 1">
              <el-button type="text" @click="removeSubject(ind)">{{ $t("commons.button.delete") }}
              </el-button>
            </div>
            <div>
              <el-row :gutter="20">
                <el-col :span="6">
                  <el-form-item :label="$t('business.configuration.kind')" prop="kind">
                    <el-select v-model="row.kind" style="width: 100%" @change="changeSubjectKind(row)">
                      <el-option label="User" :value="'User'"></el-option>
                      <el-option label="Group" :value="'Group'"></el-option>
                      <el-option label="ServiceAccount" :value="'ServiceAccount'"></el-option>
                    </el-select>
                  </el-form-item>
                </el-col>
                <div v-if="row.kind === 'ServiceAccount'">
                  <el-col :span="6">
                    <el-form-item :label="$t('business.namespace.namespace')" prop="namespace">
                      <el-select v-model="row.namespace" @change="changeNamespace(row,ind)">
                        <el-option v-for="namespace in namespaces"
                                   :key="namespace"
                                   :label="namespace"
                                   :value="namespace">
                        </el-option>
                      </el-select>
                    </el-form-item>
                  </el-col>
                  <el-col :span="6">
                    <el-form-item label="ServiceAccount" prop="serviceAccount">
                      <el-select v-model="row.name" v-loading="loading">
                        <el-option v-for="item in serviceAccountArray[ind]"
                                   :key="item.metadata.name"
                                   :label="item.metadata.name"
                                   :value="item.metadata.name">
                        </el-option>
                      </el-select>
                    </el-form-item>
                  </el-col>
                </div>
                <div v-else>
                  <el-col :span="6">
                    <el-form-item :label="$t('commons.table.name')" prop="name">
                      <el-input v-model="row.name"></el-input>
                    </el-form-item>
                  </el-col>
                </div>
              </el-row>
            </div>
          </el-card>
        </el-form>
        <el-button @click="addSubject" style="margin-top: 10px">
          {{ $t("commons.button.add") }}{{ $t("business.access_control.object") }}
        </el-button>
      </ko-card>
    </div>
  </div>
</template>

<script>
import KoCard from "@/components/ko-card"
import {listClusterRoles} from "@/api/clusterroles"
import {listNamespaceRoles} from "@/api/roles"
import {getNamespaces} from "@/api/auth"
import {listServiceAccountsWithNs} from "@/api/serviceaccounts"
import {checkPermissions} from "@/utils/permission"

export default {
  name: "KoRoleObject",
  components: { KoCard },
  props: {
    clusterName: String,
    namespace: String,
    roleRefObj: Object,
    subjectArray: Array,
  },
  data () {
    return {
      roleRef: {
        kind: "ClusterRole",
        apiGroup: "rbac.authorization.k8s.io",
        name: ""
      },
      form: {
        subjects: [],
      },
      roles: [],
      namespaces: [],
      serviceAccountArray: [[{
        metadata: { name: "" }
      }]],
      loading: false
    }
  },
  watch: {
    namespace (newValue, oldValue) {
      if (oldValue !== newValue) {
        this.changeKind(this.roleRef.kind)
      }
    },
  },
  methods: {
    changeKind (kind) {
      this.roles = []
      this.roleRef.name = ""
      if (kind === "Role") {
        if (checkPermissions({ scope: "namespace", apiGroup: "rbac.authorization.k8s.io", resource: "roles", verb: "list" })) {
          listNamespaceRoles(this.clusterName, this.namespace).then(res => {
            this.roles = res.items
          })
        }
      } else {
        if (checkPermissions({ scope: "namespace", apiGroup: "rbac.authorization.k8s.io", resource: "clusterroles", verb: "list" })) {
          listClusterRoles(this.clusterName).then(res => {
            this.roles = res.items
          })
        }
      }
    },
    removeSubject (index) {
      this.serviceAccountArray.splice(index, 1)
      this.form.subjects.splice(index, 1)
    },
    addSubject () {
      this.serviceAccountArray.push([])
      this.form.subjects.push({
        kind: "",
        name: ""
      })
    },
    changeSubjectKind (row) {
      if (row.kind === "ServiceAccount") {
        row.namespace = ""
        delete row.apiGroup
      }
      if (row.kind === "User" || row.kind === "Group") {
        row.apiGroup = "rbac.authorization.k8s.io"
        delete row.namespace
      }
    },
    changeNamespace (row, index) {
      this.loading = true
      if (checkPermissions({ scope: "namespace", apiGroup: "", resource: "serviceaccounts", verb: "list" })) {
        listServiceAccountsWithNs(this.clusterName, row.namespace).then(res => {
          this.serviceAccountArray[index] = res.items
          this.loading = false
        })
      }
    }
  },
  created () {
    getNamespaces(this.clusterName).then(res => {
      this.namespaces = res.data
    })
    this.changeKind(this.roleRef.kind)
    if (!this.roleRefObj) {
      this.$emit("update:roleRefObj", this.roleRef)
    } else {
      this.roleRef = this.roleRefObj
    }
    if (!this.subjectArray) {
      this.$emit("update:subjectArray", this.form.subjects)
    } else {
      this.form.subjects = this.subjectArray
    }
  }
}
</script>

<style scoped>

</style>
