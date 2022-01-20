<template>
  <layout-content :header="header" :back-to="{ name: 'ImageRepos' }">
    <el-row v-loading="loading">
      <el-col :span="4"><br/></el-col>
      <el-col :span="12">
        <div class="grid-content bg-purple-light">
          <el-form ref="form" :model="form" :rules="rules" label-width="250px" label-position="left">
            <el-form-item :label="$t('business.image_repos.name')" prop="name">
              <el-input v-model="form.name" :disabled="mode==='edit'"></el-input>
            </el-form-item>
            <el-form-item :label="$t('business.image_repos.type')" prop="type">
              <el-select v-model="form.type" style="width:100%" :disabled="mode==='edit'">
                <el-option :value="'Nexus'" :label="'Nexus'"></el-option>
                <el-option :value="'Harbor'" :label="'Harbor'"></el-option>
                <el-option :value="'DockerRegistry'" :label="'Docker Registry'"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('business.image_repos.endpoint')" prop="endPoint">
              <el-input v-model="form.endPoint" :placeholder="'http://172.16.10.10:8080'"></el-input>
            </el-form-item>
            <el-form-item :label="$t('business.image_repos.version')" prop="version" v-if="form.type === 'Harbor'">
              <el-select v-model="form.version" style="width:100%" >
                <el-option :value="'v1'" :label="'v1'"></el-option>
                <el-option :value="'v2'" :label="'v2'"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('business.image_repos.auth')" prop="auth">
              <el-radio-group v-model="form.auth">
                <el-radio :label="true">{{ $t("commons.bool.true") }}</el-radio>
                <el-radio :label="false">{{ $t("commons.bool.false") }}</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item v-if="form.auth" :label="$t('business.image_repos.username')" prop="credential.username">
              <el-input v-model="form.credential.username"></el-input>
            </el-form-item>
            <el-form-item v-if="form.auth" :label="$t('business.image_repos.password')" prop="credential.password">
              <el-input type="password" v-model="form.credential.password"></el-input>
            </el-form-item>
            <el-form-item v-if="form.type !== 'DockerRegistry'" :label="$t('business.image_repos.repo')"
                          prop="repoName">
              <el-select v-model="form.repoName" style="width:100%" filterable>
                <el-option v-for="(repo,index) in repos" :key="index" :value="repo" :label="repo">
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('business.image_repos.downloadUrl')" prop="downloadUrl">
              <el-input v-model="form.downloadUrl" :placeholder="'172.16.10.10:8081'"></el-input>
            </el-form-item>
            <el-form-item :label="$t('business.image_repos.allow_anonymous')" prop="allowAnonymous">
              <el-radio-group v-model="form.allowAnonymous">
                <el-radio :label="true">{{ $t("commons.bool.true") }}</el-radio>
                <el-radio :label="false">{{ $t("commons.bool.false") }}</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item>
              <div style="float: right">
                <el-button @click="list()" v-if="form.type !== 'DockerRegistry'">{{
                    $t("business.image_repos.load_repo")
                  }}
                </el-button>
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
import {createRepo, getRepo, listInternalRepos, updateRepo} from "@/api/imagerepos"
import Rules from "@/utils/rules"

export default {
  name: "ImageRepoOp",
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
        auth: true,
        credential: {},
        allowAnonymous: false,
        version: "v1"
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
        version:[
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
        allowAnonymous: [
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
      this.checkAuth()
      if (this.mode === "edit") {
        updateRepo(this.name, this.form).then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.update_success")
          })
          this.$router.push({ name: "ImageRepos" })
        }).finally(() => {
            this.isSubmitGoing = false
            this.loading = false
          }
        )
      } else {
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
      }
    },
    list () {
      this.checkAuth()
      this.loading = true
      listInternalRepos(this.form).then(res => {
        this.repos = res.data
        if (res.data === null || res.data.length === 0) {
          this.$message({
            type: "warning",
            message: this.$t("business.image_repos.repo_null")
          })
        }
      }).finally(() => {
        this.loading = false
      })
    },
    getDetail () {
      this.loading = true
      getRepo(this.name).then(res => {
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
      this.$router.push({ name: "ImageRepos" })
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

<style scoped>

</style>
