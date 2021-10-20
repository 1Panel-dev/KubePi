<template>
  <layout-content :header="$t('business.chart.create')" :back-to="{ name: 'Repos' }">
    <el-row v-loading="loading">
      <el-col :span="4"><br/></el-col>
      <el-col :span="10">
        <div class="grid-content bg-purple-light">
          <el-form ref="form" :model="form" :rules="rules" label-width="150px" label-position="left">
            <el-form-item :label="$t('business.chart.name')" prop="name">
              <el-input v-model="form.name"></el-input>
            </el-form-item>
            <!--            <el-form-item :label="$t('business.chart.type')" prop="type">-->
            <!--              <el-select v-model="form.type" style="width:100%">-->
            <!--                <el-option :value="'http'" :label="$t('business.chart.http')"></el-option>-->
            <!--                <el-option :value="'git'" :label="$t('business.chart.git')"></el-option>-->
            <!--              </el-select>-->
            <!--            </el-form-item>-->
            <!--            <div v-if="form.type==='git'">-->
            <!--              <el-form-item :label="'URL'" prop="spec.gitRepo">-->
            <!--                <el-input v-model="form.spec.gitRepo"></el-input>-->
            <!--              </el-form-item>-->
            <!--              <el-form-item :label="$t('business.chart.branch')" prop="spec.gitBranch">-->
            <!--                <el-input v-model="form.spec.gitBranch"></el-input>-->
            <!--              </el-form-item>-->
            <!--            </div>-->
            <el-form-item :label="'URL'" prop="url">
              <el-input v-model="form.url"></el-input>
            </el-form-item>
            <el-form-item :label="$t('business.chart.auth')" prop="type">
              <el-select v-model="form.type">
                <el-option :value="'None'" :label="$t('business.chart.none')"></el-option>
                <el-option :value="'Basic'" :label="$t('business.chart.basic')"></el-option>
                <!--                <el-option :value="'SSH'" :label="'SSH'"></el-option>-->
              </el-select>
            </el-form-item>
            <div v-if="form.type==='Basic'">
              <el-form-item :label="$t('business.chart.username')" prop="username">
                <el-input v-model="form.username"></el-input>
              </el-form-item>
              <el-form-item :label="$t('business.chart.password')" prop="password">
                <el-input type="password" v-model="form.password"></el-input>
              </el-form-item>
            </div>
            <div v-if="form.type==='SSH'">
              <el-form-item :label="$t('business.chart.publicKey')" prop="publicKey">
                <el-input v-model="form.publicKey"></el-input>
              </el-form-item>
              <el-form-item :label="$t('business.chart.privateKey')" prop="privateKey">
                <el-input v-model="form.privateKey"></el-input>
              </el-form-item>
            </div>
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
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import Rules from "@/utils/rules"
import {createRepo} from "@/api/charts"

export default {
  name: "RepoCreate",
  components: { LayoutContent },
  props: {},
  data () {
    return {
      loading: false,
      cluster: "",
      form: {
        type: "None"
      },
      rules: {
        name: [
          Rules.RequiredRule
        ],
        username: [
          Rules.RequiredRule
        ],
        type: [
          Rules.RequiredRule
        ],
        password: [
          Rules.RequiredRule
        ],
        url: [
          Rules.UrlRule
        ],
      },
      isSubmitGoing: false
    }
  },
  methods: {
    onCancel () {
      this.$router.push({ name: "Repos" })
    },
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
      this.form.name = this.form.name.replace(" ", "")
      createRepo(this.cluster, this.form).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success")
        })
        this.$router.push({ name: "Repos" })
      }).finally(
        () => {
          this.isSubmitGoing = false
          this.loading = false
        }
      )
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
  }
}
</script>

<style scoped>

</style>
