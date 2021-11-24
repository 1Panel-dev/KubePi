<template>
  <div>
    <el-form label-position="top" ref="form" :rules="rules" :model="form" :disabled="isReadOnly">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.container_name')" prop="name">
            <ko-form-item itemType="input" @input="changeName" v-model="form.name"/>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="8">
          <el-form-item :label="$t('business.workload.container_image')" prop="image">
            <ko-form-item placeholder="e.g. nginx:latest" itemType="input" v-model="form.image"/>
          </el-form-item>
        </el-col>
        <el-col :span="4">
          <el-form-item :label="$t('business.workload.list_image')">
            <el-select v-model="repo.name" @change="changeRepo(repo.name)">
              <el-option :value="''" :label="$t('business.workload.repo_disabled')"></el-option>
              <el-option v-for="(item,index) in repos" :key="index" :value="item.repo" :label="item.repo">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item :label="$t('business.workload.repo_image')">
            <el-select v-model="repo.image" @change="changeImage(repo.image)" style="width: 100%">
              <el-option v-for="(item,index) in repo.images" :key="index" :value="item" :label="item">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="4">
          <el-form-item :label="$t('business.workload.pull_policy')" prop="imagePullPolicy">
            <ko-form-item itemType="select" :noClear="true" v-model="form.imagePullPolicy"
                          :selections="image_pull_policy_list"/>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"
import Rule from "@/utils/rules"
import {listClusterRepos} from "../../../../kubepi/src/api/clusters"
import {getRepo, listImages} from "../../../../kubepi/src/api/imagerepos"

export default {
  name: "KoContainer",
  components: { KoFormItem },
  props: {
    containerParentObj: Object,
    secretList: Array,
    isReadOnly: Boolean,
  },
  watch: {},
  data () {
    return {
      form: {
        name: "",
        image: "",
        imagePullPolicy: "",
      },
      repo: {
        name: "",
        images: [],
        image: "",
        repo: {}
      },
      rules: {
        name: [Rule.CommonNameRule],
        image: [Rule.RequiredRule],
      },
      image_pull_policy_list: [
        { label: "Always", value: "Always" },
        { label: "IfNotPresent", value: "IfNotPresent" },
        { label: "Never", value: "Never" },
      ],
      cluster: "",
      repos: []
    }
  },
  methods: {
    changeName (val) {
      this.$emit("update:ContanerList", val)
    },
    checkIsValid () {
      let isValid = true
      this.$refs["form"].validate((valid) => {
        isValid = valid
      })
      return isValid
    },
    transformation (parentFrom) {
      parentFrom.name = this.form.name || undefined
      parentFrom.image = this.form.image || undefined
      parentFrom.imagePullPolicy = this.form.imagePullPolicy || undefined
    },
    listRepos () {
      listClusterRepos(this.cluster).then(res => {
        this.repos = res.data
      })
    },
    changeRepo (repo) {
      if (repo === "") {
        this.repo.images = []
        this.repo.repo = {}
        this.repo.name = ""
        return
      }
      listImages(this.cluster, repo).then(res => {
        this.repo.images = res.data
      })
      getRepo(repo).then(res => {
        this.repo.repo = res.data
      })
    },
    getRepoForSecret (workload, namespace) {
      if (this.repo.name !== "") {
        const auths = {
          auths: {
            [this.repo.repo.endPoint]: {
              username: this.repo.repo.credential.username,
              password: this.repo.repo.credential.password
            }
          }
        }
        const { Base64 } = require("js-base64")
        const data = {
          [".dockerconfigjson"]: Base64.encode(JSON.stringify(auths))
        }
        return {
          apiVersion: "v1",
          kind: "Secret",
          metadata: {
            name: workload + "-" + this.repo.name,
            namespace: namespace,
          },
          data: data,
          type: "kubernetes.io/dockerconfigjson"
        }
      } else {
        return {}
      }
    },
    changeImage (image) {
      this.form.image = image
    }
  },
  mounted () {
    if (this.containerParentObj) {
      if (this.containerParentObj.name) {
        this.form.name = this.containerParentObj.name
      }
      if (this.containerParentObj.image) {
        this.form.image = this.containerParentObj.image
      }
      if (this.containerParentObj.imagePullPolicy) {
        this.form.imagePullPolicy = this.containerParentObj.imagePullPolicy
      }
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.listRepos()
  }
}
</script>
