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
        <el-row>
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.node_name')">
              <ko-form-item itemType="input" disabled v-model="matchNode" />
            </el-form-item>
          </el-col>
        </el-row>
        <table style="width: 100%" class="tab-table">
          <tr>
            <th scope="col" width="50%" align="left">
              <label>{{ $t('business.workload.value') }}</label>
            </th>
            <th align="left"></th>
          </tr>
          <tr v-for="(row, index) in nodeSelector" v-bind:key="index">
            <td>
              <ko-form-item itemType="select2" @change="getMatchNode" v-model="row.value" :selections="labels" />
            </td>
            <td>
              <el-button type="text" style="font-size: 10px" @click="handleLabelsDelete(index)">
                {{ $t("commons.button.delete") }}
              </el-button>
            </td>
          </tr>
          <tr>
            <td align="left">
              <el-button @click="handleLabelsAdd()">
                {{ $t('business.workload.add') }}{{ $t('business.workload.rule') }}
              </el-button>
            </td>
          </tr>
        </table>
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
      nodeSelector: [],
      matchNode: "",
    }
  },
  methods: {
    handleLabelsAdd() {
      var item = {
        value: "",
      }
      this.nodeSelector.push(item)
    },
    handleLabelsDelete(index) {
      this.nodeSelector.splice(index, 1)
    },
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
      for (const node of this.node_list) {
        let isAllMatch = true
        for (const selector of this.nodeSelector) {
          let isOneMatch = false
          for (const key in node.metadata.labels) {
            if (Object.prototype.hasOwnProperty.call(node.metadata.labels, key)) {
              if (selector.value === key + " : " + (node.metadata.labels[key] || "")) {
                isOneMatch = true
                break
              }
            }
          }
          if (!isOneMatch) {
            isAllMatch = false
            break
          }
        }
        if (isAllMatch) {
          this.matchNode += node.metadata.name + ","
        }
      }
    },

    transformation(parentFrom) {
      if (this.scheduling_type === "specificNode") {
        parentFrom.nodeName = this.nodeName || undefined
      }
      if (this.scheduling_type === "nodeSelector") {
        parentFrom.nodeSelector = {}
        for (const item of this.nodeSelector) {
          if (item.value.indexOf(" : ") !== -1) {
            parentFrom.nodeSelector[item.value.split(" : ")[0]] = item.value.split(" : ")[1]
          }
        }
      }
    },
  },
  mounted() {
    this.clusterName = this.$route.query.cluster
    this.nodeSelector = []
    if (this.nodeSchedulingParentObj) {
      if (this.nodeSchedulingParentObj.nodeSelector) {
        this.scheduling_type = "nodeSelector"
        for (const key in this.nodeSchedulingParentObj.nodeSelector) {
          if (Object.prototype.hasOwnProperty.call(this.nodeSchedulingParentObj.nodeSelector, key)) {
            this.nodeSelector.push({
              value: key + " : " + this.nodeSchedulingParentObj.nodeSelector[key],
            })
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
