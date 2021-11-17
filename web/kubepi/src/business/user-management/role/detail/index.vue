<template>
  <layout-content :header="$t('commons.button.detail')" :back-to="{ name: 'Roles' }">

    <el-row>
      <el-col :span="4"><br/></el-col>
      <el-col :span="10">
        <div class="grid-content bg-purple-light">
          <el-form ref="form" :model="role" label-width="150px" label-position="left">

            <el-form-item :label="$t('commons.table.name')" prop="name" required>
              {{ role.name }}
            </el-form-item>

            <el-form-item :label="$t('commons.table.description')" prop="description">
              {{ translate(role.description) }}
            </el-form-item>


            <el-form-item :label="$t('business.user.permission')">
              <el-table
                  :data="resources"
                  v-loading="loading"
                  style="width: 100%">
                <el-table-column
                    :label="$t('business.user.resource_name')"
                    min-width="180">
                  <template v-slot:default="{row}">
                    {{ $t(row.name) }}
                  </template>
                </el-table-column>
                <el-table-column>
                  <template slot-scope="scope">
                    <div v-for="(item,index) in scope.row.verbs" :key="index">
                      <el-checkbox disabled
                                   v-model="item.enable">
                        {{ $t(item.name) }}
                      </el-checkbox>
                    </div>
                  </template>
                </el-table-column>
              </el-table>
            </el-form-item>
            <el-form-item style="float: right">
              <el-button @click="onEdit()"
                         :disabled="role.builtIn || !checkPermissions([{resource: 'roles', verb: 'update'}])">
                {{ $t("commons.button.edit") }}
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-col>
      <el-col :span="4"><br/></el-col>
    </el-row>


  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {listApiResource} from "@/api/apis"
import {getRole} from "@/api/roles"
import {checkPermissions} from "@/utils/permission";


export default {
  name: "RoleDetail",
  props: ["name"],
  components: {LayoutContent},
  data() {
    return {
      loading: false,
      role: {},
      resources: [],
      roles: [],
    }
  },
  methods: {
    translate(a) {
      if (a) {
        return a.startsWith("i18n_") ? this.$t(a) : a
      }
    },
    checkPermissions(r) {
      return checkPermissions(r)
    },
    onEdit() {
      this.$router.push({name: "RoleEdit", params: {name: this.name}})
    },
    onCancel() {
      this.$router.push({name: "Roles"})
    },
    onCreated() {
      const sortFunc = function (a, b) {
        if (a.name < b.name) {
          return 1
        }
        if (a.name > b.name) {
          return -1
        }
        return 0
      }
      this.loading = true
      listApiResource().then(data => {
        for (const key in data.data) {
          if (key) {
            const item = {name: key, verbs: []}
            for (const verb of data.data[key]) {
              item.verbs.push({
                name: verb,
                enable: false,
              })
            }
            item.verbs.sort(sortFunc)
            this.resources.push(item)
          }
        }
        this.resources.sort(sortFunc)
        getRole(this.name).then(data => {
          this.role = data.data
          this.role.rules.forEach(rule => {
            rule.resource.forEach(res => {
              this.resources.forEach(cres => {
                if (res === cres.name || res === "*") {
                  rule.verbs.forEach(verb => {
                    cres.verbs.forEach(cverb => {
                      if (verb === cverb.name || verb === "*") {
                        cverb.enable = true
                      }
                    })
                  })
                }
              })
            })
          })
        })
      })
      this.loading = false
    }
  },
  created() {
    this.onCreated()
  }
}
</script>

<style scoped>
</style>
