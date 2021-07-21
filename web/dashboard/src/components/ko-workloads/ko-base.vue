<template>
  <el-row :gutter="20">
    <el-col :span="12">
      <el-row>
        <el-col :span="8">
          <el-select style="margin-top: 33px;width: 100%" @change="selectContainer(true)" v-model="selectContainerType">
            <el-option v-for="(item, index) in type_list" :key="index" :label="item.label" :value="item.value" />
          </el-select>
        </el-col>
        <el-col :span="12" :key="isRefresh">
          <el-form-item :label="$t('business.workload.container')">
            <el-select v-if="selectContainerType === 'standardContainers'" @change="selectContainer(false)" style="width:100%" v-model="selectContainerIndex">
              <el-option v-for="(item, index) in containers" :key="index" :label="item.name" :value="index" />
            </el-select>
            <el-select v-if="selectContainerType === 'initContainers'" @change="selectContainer(false)" style="width:100%" v-model="selectContainerIndex">
              <el-option v-for="(item, index) in initContainers" :key="index" :label="item.name" :value="index" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="2">
          <div style="margin-top: 33px">
            <el-button style="width:100%" @click="handleAddContainer">+</el-button>
          </div>
        </el-col>
        <el-col :span="2">
          <div style="margin-top: 33px">
            <el-button style="width:100%" @click="handleDeleteContainer">-</el-button>
          </div>
        </el-col>
      </el-row>
    </el-col>
  </el-row>
</template>
 
<script>
import { deepClone } from "@/utils/deepClone"
export default {
  name: "KoCommand",
  props: {
    baseParentObj: Object, // form.spec.template.spec
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
        { label: "Init Container", value: "initContainers" },
        { label: "Standard Container", value: "standardContainers" },
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
      // if (!this.isValid) {
      //   this.$notify({ title: this.$t("commons.message_box.prompt"), message: this.unValidInfo })
      //   this.selectContainerIndex = this.currentContainerIndex
      //   this.selectContainerType = this.currentContainerType
      //   return
      // }
      this.$emit("gatherFormData")
      this.selectContainerIndex = isChangeType ? 0 : this.selectContainerIndex
      if (this.selectContainerType === "initContainers") {
        if (!this.initContainers || this.initContainers.length === 0) {
          this.$confirm(this.$t("commons.confirm_message.add_init_container"), this.$t("commons.message_box.prompt"), {
            confirmButtonText: this.$t("commons.button.confirm"),
            cancelButtonText: this.$t("commons.button.cancel"),
            type: "warning",
          })
            .then(() => {
              const itemContainer = { name: "initContainer-0", image: "", imagePullPolicy: "ifNotPresent" }
              this.initContainers = [itemContainer]
              this.currentContainerType = this.selectContainerType
              this.currentContainerIndex = this.selectContainerIndex
              this.currentContainer = this.initContainers[this.currentContainerIndex]
              this.$emit("addContainer", this.currentContainerType, itemContainer)
              this.$emit("refreshContainer", this.currentContainerType, this.currentContainerIndex, this.currentContainer)
            })
            .catch(() => {
              this.selectContainerType = this.currentContainerType
              this.selectContainerIndex = this.currentContainerIndex
            })
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
      itemContainer.imagePullPolicy = "ifNotPresent"
      if (this.selectContainerType === "initContainers") {
        if (this.initContainers) {
          itemContainer.name = "initContainer-" + this.initContainers.length.toString()
        } else {
          this.initContainers = []
          itemContainer.name = "initContainer-0"
        }
        this.initContainers.push(itemContainer)
      } else {
        itemContainer.name = "Container-" + this.containers.length.toString()
        this.containers.push(itemContainer)
      }
      this.$emit("addContainer", this.currentContainerType, itemContainer)
    },
    handleDeleteContainer() {
      this.selectContainerIndex = 0
      this.$emit("deleteContainer", this.currentContainerType, this.currentContainerIndex)
      if (this.selectContainerType === "initContainers") {
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