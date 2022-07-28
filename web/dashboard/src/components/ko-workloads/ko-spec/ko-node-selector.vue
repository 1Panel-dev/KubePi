<template>
  <div>
    <el-form label-position="top" ref="form" :disabled="isReadOnly">
      <el-row>
        <el-col :span="24">
          <el-form-item :label="$t('business.workload.schedule_type')">
            <ko-form-item itemType="radio" v-model="scheduling_type" :radios="scheduling_type_list" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row v-if="scheduling_type === 'specificNode'">
        <el-col :span="24">
          <el-form-item :label="$t('business.workload.node_name')">
            <ko-form-item itemType="select2" v-model="nodeName" :selections="node_name_list" />
          </el-form-item>
        </el-col>
      </el-row>

      <div v-if="scheduling_type === 'nodeSelector'">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.node_name')">
              <ko-form-item itemType="input" disabled v-model="matchNode" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.value')">
              <ko-form-item itemType="select2" @change="getMatchNode" v-model="nodeSelector" :selections="labels" />
            </el-form-item>
          </el-col>
        </el-row>
      </div>
    </el-form>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoNodeScheduling",
  components: { KoFormItem },
  props: {
    nodeSchedulingParentObj: Object,
    isReadOnly: Boolean,
    nodeList: Array,
    nodeSchedulingType: String,
  },
  watch: {
    nodeList: {
      handler(newName) {
        this.node_name_list = []
        this.node_list = []
        if (newName) {
          this.loadNodes(newName)
        }
      },
      immediate: true,
    },
  },
  data() {
    return {
      scheduling_type_list: [
        { label: this.$t("business.workload.any_node"), value: "anyNode" },
        { label: this.$t("business.workload.specific_node"), value: "specificNode" },
        { label: this.$t("business.workload.match_rules"), value: "nodeSelector" },
      ],
      scheduling_type: "anyNode",
      nodeName: "",
      node_name_list: [],
      node_list: [],
      labels: [],
      nodeSelector: "",
      matchNode: "",
    }
  },
  methods: {
    loadNodes(datas) {
      this.node_name_list = []
      this.node_list = datas
      let labels = []
      for (const node of datas) {
        this.node_name_list.push(node.metadata.name)
        for (const key in node.metadata.labels) {
          if (Object.prototype.hasOwnProperty.call(node.metadata.labels, key)) {
            labels.push(key + " : " + (node.metadata.labels[key] || ""))
          }
        }
      }
      this.labels = labels.filter((item, index) => {
        return labels.indexOf(item) === index
      })
      this.getMatchNode()
    },

    getMatchNode() {
      this.matchNode = ""
      if (this.nodeSelector === "") {
        this.matchNode = this.node_name_list.join(",")
        return
      }
      let itemList = []
      for (const node of this.node_list) {
        let isMatch = false
        for (const key in node.metadata.labels) {
          if (Object.prototype.hasOwnProperty.call(node.metadata.labels, key)) {
            if (this.nodeSelector === key + " : " + (node.metadata.labels[key] || "")) {
              isMatch = true
              break
            }
          }
        }
        if (isMatch) {
          itemList.push(node.metadata.name)
        }
      }
      this.matchNode = itemList.join(",")
    },

    transformation(parentFrom) {
      if (this.scheduling_type === "specificNode") {
        parentFrom.nodeName = this.nodeName || undefined
      }
      if (this.scheduling_type === "nodeSelector") {
        parentFrom.nodeSelector = {}
        if (this.nodeSelector.indexOf(" : ") !== -1) {
          parentFrom.nodeSelector[this.nodeSelector.split(" : ")[0]] = this.nodeSelector.split(" : ")[1]
        }
      }
    },
  },
  mounted() {
    this.clusterName = this.$route.query.cluster
    this.nodeSelector = ""
    if (this.nodeSchedulingParentObj) {
      if (this.nodeSchedulingParentObj.nodeSelector) {
        this.scheduling_type = "nodeSelector"
        for (const key in this.nodeSchedulingParentObj.nodeSelector) {
          if (Object.prototype.hasOwnProperty.call(this.nodeSchedulingParentObj.nodeSelector, key)) {
            this.nodeSelector = key + " : " + this.nodeSchedulingParentObj.nodeSelector[key]
          }
        }
      } else if (this.nodeSchedulingParentObj.nodeName) {
        this.scheduling_type = "specificNode"
        this.nodeName = this.nodeSchedulingParentObj.nodeName
      } else {
        this.scheduling_type = "anyNode"
      }
    }
  },
}
</script>
