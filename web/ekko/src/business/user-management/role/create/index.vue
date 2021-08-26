<template>
  <layout-content :header="$t('commons.button.create')+$t('business.user.role')" :back-to="{ name: 'Roles' }">
    <el-row>
      <el-col :span="4"><br/></el-col>
      <el-col :span="10">
        <div class="grid-content bg-purple-light">
          <el-form ref="form" :model="form" :rules="rules" label-width="150px" label-position="left">

            <el-form-item :label="$t('commons.table.name')" prop="name" required>
              <el-input v-model="form.name"></el-input>
            </el-form-item>

            <el-form-item :label="$t('commons.table.description')" prop="description">
              <el-input v-model="form.description"></el-input>
            </el-form-item>
            <el-form-item :label="$t('business.user.permission_setting')">
              <el-table
                  :data="resources"
                  v-loading="loading"
                  style="width: 100%">
                <el-table-column
                    :label="$t('business.user.resource_name')"
                    min-width="200">
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
            <el-form-item>
              <div style="float: right">
                <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
                <el-button type="primary" @click="onConfirm()">{{ $t("commons.button.confirm") }}
                </el-button>
              </div>
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
import {createRole} from "@/api/roles"
import Rule from "@/utils/rules"
import {listApiResource} from "@/api/apis"


export default {
  name: "RoleCreate",
  components: {LayoutContent},
  data() {
    return {
      loading: false,
      isSubmitGoing: false,
      allSelect: false,
      indeterminate: true,
      resources: [],
      rules: {
        name: [Rule.RequiredRule],
      },
      form: {
        name: "",
        description: "",
      },
    }
  },
  methods: {
    onAllSelectChange() {
      this.indeterminate = false
      this.resources.forEach((r) => {
        r.verbs.forEach(v => {
          v.enable = this.allSelect
        })
      })
    },
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

      const policyRules = []
      this.resources.forEach((r) => {
        const rule = {resource: [r.name], verbs: []}
        r.verbs.forEach((v) => {
          if (v.enable) {
            rule.verbs.push(v.name)
          }
        })
        policyRules.push(rule)
      })

      const req = {
        "apiVersion": "v1",
        "kind": "Role",
        "name": this.form.name,
        "description": this.form.description,
        "rules": policyRules
      }
      createRole(req).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success")
        })
        this.$router.push({name: "Roles"})
      }).finally(() => {
        this.isSubmitGoing = false
      })
    },
    onCancel() {
      this.$router.push({name: "Roles"})
    },

    initData() {
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
      })
      this.loading = false
    }
  },
  created() {
    this.initData()
  }
}
</script>

<style scoped>
</style>
