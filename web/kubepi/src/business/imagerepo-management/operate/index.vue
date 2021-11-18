<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'ImageRepos' }">
    <el-row v-loading="loading">
      <el-col :span="4"><br/></el-col>
      <el-col :span="10">
        <div class="grid-content bg-purple-light">
          <el-form ref="form" :model="form" :rules="rules" label-width="150px" label-position="left">
            <el-form-item :label="$t('business.image_repos.name')" prop="name">
              <el-input v-model="form.name"></el-input>
            </el-form-item>
            <el-form-item :label="$t('business.image_repos.type')" prop="type">
              <el-select v-model="form.type">
                <el-option :value="'Nexus'" :label="'Nexus'"></el-option>
                <el-option :value="'Harbor'" :label="'Harbor'"></el-option>
                <el-option :value="'DockerRegistry'" :label="'Docker Registry'"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('business.image_repos.endpoint')" prop="endpoint">
              <el-input v-model="form.endpoint"></el-input>
            </el-form-item>
            <el-form-item :label="$t('business.image_repos.username')" prop="credential.username">
              <el-input v-model="form.credential.username"></el-input>
            </el-form-item>
            <el-form-item :label="$t('business.image_repos.password')" prop="credential.password">
              <el-input type="password" v-model="form.credential.password"></el-input>
            </el-form-item>
            <el-form-item :label="$t('business.image_repos.repo')" prop="repoName">
              <el-select v-model="form.repoName">
                <el-option v-for="(repo,index) in repos" :key="index" :value="repo" :label="repo">
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('business.image_repos.downloadUrl')" prop="downloadUrl">
              <el-input v-model="form.downloadUrl"></el-input>
            </el-form-item>
            <el-form-item>
              <div style="float: right">
                <el-button @click="list()">{{ $t("business.image_repos.load_repo") }}</el-button>
                <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
                <el-button type="primary" :disabled="isSubmitGoing"  @click="onConfirm()">{{ $t("commons.button.confirm") }}
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
import {createRepo, listInternalRepos} from "@/api/imagerepos"
import Rules from "@/utils/rules"

export default {
  name: "ImageRepoOp",
  components: { LayoutContent },
  props: {},
  data () {
    return {
      loading: false,
      mode: "",
      form: {
        credential:{}
      },
      rules: {
        name: [
          Rules.RequiredRule
        ],
        type: [
          Rules.RequiredRule
        ],
        endpoint: [
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
        ]
      },
      repos: [],
      isSubmitGoing: false
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
      createRepo(this.form).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success")
        })
        this.$router.push({ name: "ImageRepos" })
      }).finally(() => {
          this.isSubmitGoing = false
          this.loading = false
        }
      )
    },
    list () {
      this.loading = true
      listInternalRepos(this.form).then(res => {
        this.repos = res.data
      }).finally(() => {
        this.loading = false
      })
    },
    onCancel () {
      this.$router.push({ name: "ImageRepos" })
    },
  },
  created () {

  }
}
</script>

<style scoped>

</style>
