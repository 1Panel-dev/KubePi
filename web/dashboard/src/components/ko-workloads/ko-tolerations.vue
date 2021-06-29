<template>
  <div style="margin-top: 20px">
    <ko-card title="Tolerations">
      <table style="width: 98%" class="tab-table">
        <tr>
          <th scope="col" width="25%" align="left">
            <label>label Key</label>
          </th>
          <th scope="col" width="10%" align="left">
            <label>operator</label>
          </th>
          <th scope="col" width="25%" align="left">
            <label>value</label>
          </th>
          <th scope="col" width="16%" align="left">
            <label>effect</label>
          </th>
          <th scope="col" width="15%" align="left">
            <label>Toleration Seconds(s)</label>
          </th>
          <th align="left"></th>
        </tr>
        <tr v-for="(row, index) in tolerations" v-bind:key="index">
          <td>
            <ko-form-item :withoutLabel="true" itemType="input" v-model="row.key" />
          </td>
          <td>
            <ko-form-item :withoutLabel="true" itemType="select" v-model="row.operator" :selections="operator_list" />
          </td>
          <td>
            <ko-form-item :withoutLabel="true" itemType="input" v-model="row.value" />
          </td>
          <td>
            <ko-form-item :withoutLabel="true" itemType="select" v-model="row.effect" :selections="effect_list" />
          </td>
          <td>
            <ko-form-item :withoutLabel="true" :disabled="row.effect !== 'NoExecute'" itemType="number" v-model.number="row.tolerationSeconds" />
          </td>
          <td>
            <el-button type="text" style="font-size: 10px" @click="handleTolerationsDelete(index)">
              {{ $t("commons.button.delete") }}
            </el-button>
          </td>
        </tr>
        <tr>
          <td align="left">
            <el-button @click="handleTolerationsAdd">Add Toleration</el-button>
          </td>
        </tr>
      </table>
    </ko-card>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"
import KoCard from "@/components/ko-card/index"

export default {
  name: "KoPodScheduling",
  components: { KoFormItem, KoCard },
  props: {
    tolerationsParentObj: Object,
  },
  data() {
    return {
      operator_list: [
        { label: "Exists", value: "Exists" },
        { label: "=", value: "Equal" },
      ],
      effect_list: [
        { label: "All", value: "All" },
        { label: "NoSchedule", value: "NoSchedule" },
        { label: "PreferNoSchedule", value: "PreferNoSchedule" },
        { label: "NoExecute", value: "NoExecute" },
      ],
      namespace_list: [
        { label: "kube-system", value: "kube-system" },
        { label: "kube-public", value: "kube-public" },
        { label: "kube-operator", value: "kube-operator" },
        { label: "default", value: "default" },
      ],
      tolerations: [],
    }
  },
  methods: {
    handleTolerationsAdd() {
      var item = {
        key: "",
        operator: "",
        value: "",
        effect: "",
        tolerationSeconds: "",
      }
      this.tolerations.push(item)
    },
    handleTolerationsDelete(index) {
      this.tolerations.splice(index, 1)
    },
    transformation(parentFrom) {
      if (this.tolerations.length !== 0) {
        parentFrom.tolerations = this.tolerations
      }
    },
  },
  mounted() {
    if (this.tolerationsParentObj) {
      if (this.tolerationsParentObj.tolerations) {
        this.tolerations = this.tolerationsParentObj.tolerations
      }
    }
  },
}
</script>