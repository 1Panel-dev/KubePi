<template>
  <div>
    <el-form :disabled="isReadOnly">
      <table style="width: 98%" class="tab-table">
        <tr>
          <th scope="col" width="25%" align="left">
            <label>{{$t('business.workload.key')}}</label>
          </th>
          <th scope="col" width="10%" align="left">
            <label>{{$t('business.workload.operator')}}</label>
          </th>
          <th scope="col" width="25%" align="left">
            <label>{{$t('business.workload.value')}}</label>
          </th>
          <th scope="col" width="16%" align="left">
            <label>{{$t('business.workload.effect')}}</label>
          </th>
          <th scope="col" width="15%" align="left">
            <label>{{$t('business.workload.toleration_seconds')}}</label>
          </th>
          <th align="left"></th>
        </tr>
        <tr v-for="(row, index) in tolerations" v-bind:key="index">
          <td>
            <ko-form-item itemType="input" v-model="row.key" />
          </td>
          <td>
            <ko-form-item itemType="select" v-model="row.operator" :selections="operator_list" />
          </td>
          <td>
            <ko-form-item v-if="row.operator === 'Exists'" disabled placeholder="N/A" itemType="input" v-model="row.value" />
            <ko-form-item v-if="row.operator !== 'Exists'" itemType="input" v-model="row.value" />
          </td>
          <td>
            <ko-form-item itemType="select" v-model="row.effect" :selections="effect_list" />
          </td>
          <td>
            <ko-form-item :disabled="row.effect !== 'NoExecute'" itemType="number" :deviderName="$t('business.workload.seconds')" v-model.number="row.tolerationSeconds" />
          </td>
          <td>
            <el-button type="text" style="font-size: 10px" @click="handleTolerationsDelete(index)">
              {{ $t("commons.button.delete") }}
            </el-button>
          </td>
        </tr>
        <tr>
          <td align="left">
            <el-button @click="handleTolerationsAdd">{{$t('business.workload.add')}}{{$t('business.workload.toleration')}}</el-button>
          </td>
        </tr>
      </table>
    </el-form>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoPodScheduling",
  components: { KoFormItem },
  props: {
    tolerationsParentObj: Object,
    isReadOnly: Boolean,
  },
  data() {
    return {
      operator_list: [
        { label: "Exists", value: "Exists" },
        { label: "=", value: "Equal" },
      ],
      effect_list: [
        { label: this.$t("business.workload.all"), value: "All" },
        { label: this.$t("business.workload.no_schedule"), value: "NoSchedule" },
        { label: this.$t("business.workload.prefer_no_schedule"), value: "PreferNoSchedule" },
        { label: this.$t("business.workload.no_execute"), value: "NoExecute" },
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
      let tolerations = []
      for (const item of this.tolerations) {
        tolerations.push({
          key: item.key || undefined,
          operator: item.operator || undefined,
          value: item.value || undefined,
          effect: item.effect || undefined,
          tolerationSeconds: item.tolerationSeconds || undefined,
        })
      }
      parentFrom.tolerations = tolerations.length !== 0 ? tolerations : undefined
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