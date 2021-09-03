<template>
  <el-form>
    <el-row style="margin-top: 20px;">
      <el-col :span="6">
        <el-radio-group style="width: 100%" @change="selectContainer(true)" v-model="selectContainerType">
          <el-radio-button :disabled="isReadOnly && initContainers === null" style="width: 50%" v-for="(item, index) in type_list" :key="index" :label="item.value">{{item.label}}</el-radio-button>
        </el-radio-group>
      </el-col>
      <el-col :span="14" :key="isRefresh">
        <el-form-item>
          <el-select v-if="selectContainerType === 'standardContainers'" @change="selectContainer(false)" style="width:100%" v-model="selectContainerIndex">
            <el-option v-for="(item, index) in containers" :key="index" :label="item.name" :value="index" />
          </el-select>
          <el-select v-if="selectContainerType === 'initContainers'" @change="selectContainer(false)" style="width:100%" v-model="selectContainerIndex">
            <el-option v-for="(item, index) in initContainers" :key="index" :label="item.name" :value="index" />
          </el-select>
        </el-form-item>
      </el-col>
      <el-col :span="2">
        <div>
          <el-button :disabled="isReadOnly" style="width:100%" @click="handleAddContainer">+</el-button>
        </div>
      </el-col>
      <el-col :span="2">
        <div>
          <el-button :disabled="isReadOnly" style="width:100%" @click="handleDeleteContainer">-</el-button>
        </div>
      </el-col>
    </el-row>
  </el-form>
</template>
 
<script>
import { deepClone } from "@/utils/deepClone"
export default {
  name: "KoCommand",
  props: {
    baseParentObj: Object, // form.spec.template.spec
    isReadOnly: Boolean,
  },
  watch: {
    baseParentObj: {
      handler(newObj) {
        if (newObj) {
          if (newObj.containers) {
            this.containers = deepClone(newObj.containers)
          }
          if (newObj.initContainers) {
            this.initContainers = deepClone(newObj.initContainers)
          }
        }
      },
      immediate: true,
      deep: true,
    },
  },
  data() {
    return {
      isRefresh: false,
      initContainers: null,
      containers: [],
      type_list: [
        { label: this.$t("business.workload.initContainer"), value: "initContainers" },
        { label: this.$t("business.workload.standardContainer"), value: "standardContainers" },
      ],
      selectContainerIndex: 0,
      currentContainerIndex: 0,
      selectContainerType: "standardContainers",
      currentContainerType: "standardContainers",
      currentContainer: {},
    }
  },
  methods: {
    selectContainer(isChangeType) {
      this.$emit("gatherFormData")
      this.selectContainerIndex = isChangeType ? 0 : this.selectContainerIndex
      if (this.selectContainerType === "initContainers") {
        if (!this.initContainers || this.initContainers.length === 0) {
          const itemContainer = { name: "init_container-0", image: "", imagePullPolicy: "IfNotPresent" }
          this.initContainers = [itemContainer]
          this.currentContainerType = this.selectContainerType
          this.currentContainerIndex = this.selectContainerIndex
          this.currentContainer = this.initContainers[this.currentContainerIndex]
          this.$emit("addContainer", this.currentContainerType, itemContainer)
          this.$emit("refreshContainer", this.currentContainerType, this.currentContainerIndex, this.currentContainer)
        } else {
          this.currentContainerType = this.selectContainerType
          this.currentContainerIndex = this.selectContainerIndex
          this.currentContainer = this.initContainers[this.currentContainerIndex]
          this.$emit("refreshContainer", this.currentContainerType, this.currentContainerIndex, this.currentContainer)
          this.isRefresh = !this.isRefresh
        }
      } else {
        this.currentContainerType = this.selectContainerType
        this.currentContainerIndex = this.selectContainerIndex
        this.currentContainer = this.containers[this.currentContainerIndex]
        this.$emit("refreshContainer", this.currentContainerType, this.currentContainerIndex, this.currentContainer)
        this.isRefresh = !this.isRefresh
      }
    },
    handleAddContainer() {
      let itemContainer = {}
      itemContainer.image = ""
      itemContainer.imagePullPolicy = "IfNotPresent"
      if (this.selectContainerType === "initContainers") {
        if (this.initContainers) {
          itemContainer.name = "init_container-" + this.initContainers.length.toString()
        } else {
          this.initContainers = []
          itemContainer.name = "init_container-0"
        }
        this.initContainers.push(itemContainer)
      } else {
        itemContainer.name = "container-" + this.containers.length.toString()
        this.containers.push(itemContainer)
      }
      this.$emit("addContainer", this.currentContainerType, itemContainer)
    },
    handleDeleteContainer() {
      this.selectContainerIndex = 0
      if (this.selectContainerType === "initContainers") {
        this.$emit("deleteContainer", this.currentContainerType, this.currentContainerIndex)
        if (this.initContainers.length <= this.currentContainerIndex) {
          return
        }
        if (this.initContainers.length === 1) {
          this.initContainers = null
          this.selectContainerType = "standardContainers"
          this.currentContainerType = "standardContainers"
          this.currentContainer = this.containers[0]
        } else {
          this.initContainers.splice(this.currentContainerIndex, 1)
          this.currentContainer = this.initContainers[0]
        }
      } else {
        if (this.containers.length <= this.currentContainerIndex || this.containers.length === 1) {
          return
        } else {
          this.$emit("deleteContainer", this.currentContainerType, this.currentContainerIndex)
          this.containers.splice(this.currentContainerIndex, 1)
          this.currentContainer = this.containers[0]
        }
      }
      this.currentContainerIndex = 0
      this.isRefresh = !this.isRefresh
      this.$emit("refreshContainer", this.currentContainerType, this.currentContainerIndex, this.currentContainer)
    },
  },
}
</script>