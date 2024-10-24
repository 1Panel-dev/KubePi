<template>
  <layout-content :header="header" :back-to="{ name: 'Metrics' }">
    <el-row v-loading="loading">
      <el-col :span="4"><br/></el-col>
      <el-col :span="12">
        <div class="grid-content bg-purple-light">
          <el-form ref="form" :model="form" :rules="rules" label-width="250px" label-position="left">
            <el-form-item :label="$t('commons.table.name')" prop="name">
              <el-input v-model="form.name" :disabled="mode==='edit'"></el-input>
            </el-form-item>
            <el-form-item :label="$t('business.monitor.metrics.type')" prop="type">
              <el-select v-model="form.type" style="width:100%" :disabled="mode==='edit'">
                <el-option :value="'Prometheus'" :label="'Prometheus'"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('business.monitor.metrics.endpoint')" prop="endPoint">
              <el-input v-model="form.endPoint" :placeholder="'http://192.168.56.101:9090'"></el-input>
            </el-form-item>
            <el-form-item :label="$t('business.monitor.metrics.auth')" prop="auth">
              <el-radio-group v-model="form.auth">
                <el-radio :label="true">{{ $t("commons.bool.true") }}</el-radio>
                <el-radio :label="false">{{ $t("commons.bool.false") }}</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item v-if="form.auth" :label="$t('business.monitor.metrics.username')" prop="credential.username">
              <el-input v-model="form.credential.username"></el-input>
            </el-form-item>
            <el-form-item v-if="form.auth" :label="$t('business.monitor.metrics.password')" prop="credential.password">
              <el-input type="password" v-model="form.credential.password"></el-input>
            </el-form-item>
            <el-form-item>
              <div style="float: right">
                <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
                <el-button type="primary" :disabled="isSubmitGoing" @click="onConfirm()">{{
                    $t("commons.button.confirm")
                  }}
                </el-button>
              </div>
            </el-form-item>
          </el-form>
        </div>
      </el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {createMetrics, getMetrics, updateMetrics} from "@/api/monitor"
import Rules from "@/utils/rules"

export default {
  components: { LayoutContent },
  props: {
    name: String
  },
  data () {
    return {
      loading: false,
      mode: "",
      header: this.$t('commons.button.create'),
      form: {
        auth: false,
        credential: {},
      },
      rules: {
        name: [
          Rules.RequiredRule,
          Rules.CommonNameRule
        ],
        type: [
          Rules.RequiredRule
        ],
        endPoint: [
          Rules.RequiredRule
        ],
        auth: [
          Rules.RequiredRule
        ],
        credential: {
          username: [
            Rules.RequiredRule
          ],
          password: [
            Rules.RequiredRule
          ],
        },
        repoName: [
          Rules.RequiredRule
        ],
        downloadUrl: [
          Rules.RequiredRule
        ],
      },
      repos: [],
      isSubmitGoing: false,
      searchRequest: {
        page:1,
        limit:20,
        search:"",
      }
    }
  },
  methods: {
    onConfirm () {
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
      this.checkAuth()
      if (this.mode === "edit") {
        updateMetrics(this.name, this.form).then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.update_success")
          })
          this.$router.push({ name: "Metrics" })
        }).finally(() => {
            this.isSubmitGoing = false
            this.loading = false
          }
        )
      } else {
        createMetrics(this.form).then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.create_success")
          })
          this.$router.push({ name: "Metrics" })
        }).finally(() => {
            this.isSubmitGoing = false
            this.loading = false
          }
        )
      }
    },
    list () {
      this.checkAuth()
      const request = {...this.form}
      request.page = 1
      request.limit = 20
      this.search(request)
    },
    searchByName(name){
      const request = {...this.form}
      request.page = this.searchRequest.page
      request.limit = this.searchRequest.limit
      request.search = name
      this.search(request)
    },
    getDetail () {
      this.loading = true
      getMetrics(this.name).then(res => {
        this.form = res.data
        this.loading = false
      })
    },
    checkAuth () {
      if (this.form.auth === false) {
        this.form.credential = {
          username: "",
          password: ""
        }
      }
    },
    onCancel () {
      this.$router.push({ name: "Metrics" })
    },
  },
  created () {
    this.mode = this.$route.query.mode
    if (this.mode === "edit") {
      this.getDetail()
      this.header = this.$t('commons.button.edit')
    }
  }
}
</script>
