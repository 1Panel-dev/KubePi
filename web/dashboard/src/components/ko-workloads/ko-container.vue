<template>
  <div>
    <el-form label-position="top" ref="form" :rules="rules" :model="form" :disabled="isReadOnly">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.container_name')" prop="name">
            <ko-form-item itemType="input" @input="changeName" v-model="form.name" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="4">
          <el-form-item :label="$t('business.workload.list_image')">
            <el-select style="width: 100%" v-model="repo.name" @change="changeRepo(repo.name)">
              <el-option :value="''" :label="$t('business.workload.repo_disabled')"></el-option>
              <el-option v-for="(item,index) in repos" :key="index" :value="item.name" :label="item.name">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item :label="$t('business.workload.container_image')" prop="image" :rules="inputRule" v-if="repo.name===''">
            <ko-form-item placeholder="e.g. nginx:latest" itemType="input" v-model="form.image" />
          </el-form-item>
          <el-form-item :label="$t('business.workload.container_image')" prop="image" :rules="selectRule" v-else>
            <el-select v-model="form.image" @change="changeImage(form.image)" style="width: 100%" filterable>
              <el-option v-for="(item,index) in repo.images" :key="index" :value="item" :label="item">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="4">
          <el-form-item :label="$t('business.workload.pull_policy')" prop="imagePullPolicy">
            <ko-form-item itemType="select" :noClear="true" v-model="form.imagePullPolicy" :selections="image_pull_policy_list" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row v-if="repo.name !== ''" :gutter="20" style="margin-left: 2px">
        <el-checkbox v-model="checked">{{$t('business.workload.repo_help')}}</el-checkbox>
      </el-row>
    </el-form>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"
import Rule from "@/utils/rules"
import { getRepo, listImages } from "../../../../kubepi/src/api/imagerepos"

export default {
  name: "KoContainer",
  components: { KoFormItem },
  props: {
    containerParentObj: Object,
    containerType: String,
    secretList: Array,
    isReadOnly: Boolean,
    repoList: Array,
    metadata: Object,
  },
  watch: {
    repoList: {
      handler(newObj) {
        this.repos = []
        if (newObj) {
          this.repos = newObj
        }
      },
      immediate: true,
      deep: true,
    },
    metadata: {
      handler(newObj) {
        if (newObj?.labels) {
          this.cluster = this.$route.query.cluster
          let itemName = ""
          if (this.containerType === "initContainers") {
            itemName = "kubepi-repo-init-" + this.containerParentObj.name
          } else {
            itemName = "kubepi-repo-" + this.containerParentObj.name
          }
          for (const key in newObj.labels) {
            if (key === itemName) {
              this.repo.name = newObj.labels[itemName]
              this.changeRepo(this.repo.name)
              break
            }
          }
        }
      },
      immediate: true,
      deep: true,
    },
  },
  data() {
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
        repo: {},
      },
      checked: false,
      rules: {
        name: [Rule.CommonNameRule],
      },
      selectRule: Rule.SelectRule,
      inputRule: Rule.RequiredRule,
      image_pull_policy_list: [
        { label: "Always", value: "Always" },
        { label: "IfNotPresent", value: "IfNotPresent" },
        { label: "Never", value: "Never" },
      ],
      cluster: "",
      repos: [],
    }
  },
  methods: {
    changeName(val) {
      this.$emit("updateContanerList", val)
    },
    checkIsValid() {
      let isValid = true
      this.$refs["form"].validate((valid) => {
        isValid = valid
      })
      return isValid
    },
    transformation(parentFrom, metadata) {
      parentFrom.name = this.form.name || undefined
      parentFrom.image = this.form.image || undefined
      parentFrom.imagePullPolicy = this.form.imagePullPolicy || undefined
      if (!metadata.labels) {
        metadata.labels = {}
      } else {
        for (let key in metadata.labels) {
          if (key.indexOf("kubepi-repo-") !== -1) {
            delete metadata.labels[key]
          }
        }
      }
      if (this.repo.name === "") {
        delete metadata.labels["operation"]
        return
      }
      let secrets = ""
      if (this.containerType === "standardContainers") {
        secrets = "kubepi-repo-" + parentFrom.name
      } else {
        secrets = "kubepi-repo-init-" + parentFrom.name 
      }
      metadata.labels[secrets] = this.repo.name
      metadata.labels["operation"] = this.checked ? "update" : "check"
    },
    changeRepo(repo) {
      this.repo.images = []
      if (repo === "") {
        this.repo.repo = {}
        this.repo.name = ""
        return
      }
      listImages(this.cluster, repo).then((res) => {
        this.repo.images = res.data
      })
      getRepo(repo).then((res) => {
        this.repo.repo = res.data
      })
    },
    changeImage(image) {
      this.form.image = image
    },
  },
  mounted() {
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
  created() {
    this.cluster = this.$route.query.cluster
  },
}
</script>
