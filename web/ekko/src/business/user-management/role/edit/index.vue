<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{ name: 'Roles' }">

    <el-row>
      <el-col :span="4"><br/></el-col>
      <el-col :span="10">
        <div class="grid-content bg-purple-light">
          <el-form ref="form" :model="role" :rules="rules" label-width="150px" label-position="left">

            <el-form-item :label="$t('commons.table.name')" prop="name" required>
              {{ role.name }}
            </el-form-item>

            <el-form-item :label="$t('commons.table.description')" prop="description">
              <el-input v-model="role.description"></el-input>
            </el-form-item>


            <el-form-item :label="$t('business.user.permission_setting')">
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
                  <template #header>
                    <div>
                      <el-checkbox :indeterminate="indeterminate" v-model="allSelect"
                                   @change="onAllSelectChange">
                        {{ $t('commons.button.all_select') }}
                      </el-checkbox>
                    </div>
                  </template>
                  <template slot-scope="scope">
                    <div v-for="(item,index) in scope.row.verbs" :key="index">
                      <el-checkbox
                          v-model="item.enable">
                        {{ $t(item.name) }}
                      </el-checkbox>
                    </div>
                  </template>
                </el-table-column>
              </el-table>
            </el-form-item>

            <el-form-item style="float: right">
              <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
              <el-button type="primary" @click="onConfirm">{{ $t("commons.button.confirm") }}</el-button>
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
import {getRole, updateRole} from "@/api/roles"
import Rule from "@/utils/rules"


export default {
  name: "RoleEdit",
  props: ["name"],
  components: {LayoutContent},
  data() {
    return {
      loading: false,
      isSubmitGoing: false,
      allSelect: false,
      indeterminate: true,
      role: {},
      resources: [],
      roles: [],
      rules: {
        name: [Rule.RequiredRule]
      },
    }
  },
  methods: {
    onConfirm() {
      if (this.isSubmitGoing) {
        return
      }
      let isFormReady = false
      this.$refs["form"].validate((valid) => {
        if (valid) {
          isFormReady = true
        }
      })
      if (!isFormReady) {
        return
      }
      this.isSubmitGoing = true
      this.loading = true
      const policyRules = []
      this.resources.forEach((r) => {
        const rule = {resource: [r.name], verbs: []}
        r.verbs.forEach((v) => {
          if (v.enable) {
            rule.verbs.push(v.name)
          }
        })
        policyRules.push(rule)
        this.role.rules = policyRules
      })
      updateRole(this.name, this.role).then(data => {
        this.$message.success(this.$t("commons.msg.update_success"))
        this.$router.push({name: "RoleEdit", params: {name: data.data.name}})
        this.loading = false
        this.isSubmitGoing = false
        this.$router.push({name: "Roles"})
      })
    },
    onCancel() {
      this.$router.push({name: "Roles"})
    },
    onAllSelectChange() {
      this.indeterminate = false
      this.resources.forEach((r) => {
        r.verbs.forEach(v => {
          v.enable = this.allSelect
        })
      })
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
