<template>
  <div>
    <el-form label-position="top" :disabled="isReadOnly">
      <table style="width: 100%" class="tab-table">
        <tr v-for="(row, index) in matchTables" v-bind:key="index">
          <td width="40%">
            <ko-form-item itemType="input" v-model="row.key" />
          </td>
          <td width="16%">
            <ko-form-item itemType="select" v-model="row.operator" :selections="operator_list" />
          </td>
          <td width="40%">
            <ko-form-item v-if="row.operator === 'Exists' || row.operator === 'DoesNotExist'" disabled itemType="input" value="N/A" />
            <ko-form-item v-else itemType="input" v-model="row.value" />
          </td>
          <td>
            <el-button type="text" style="font-size: 10px" @click="handleMatchDelete(index)">
              {{ $t("commons.button.delete") }}
            </el-button>
          </td>
        </tr>
        <tr>
          <td align="left">
            <el-button @click="handleMatchAdd()">{{ $t("commons.button.add") }}</el-button>
          </td>
        </tr>
      </table>
    </el-form>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoMatchTable",
  components: { KoFormItem },
  props: {
    matchTables: Array,
    isReadOnly: Boolean,
  },
  watch: {
    matchTables: {
      handler(newName) {
        this.matchTables = newName
      },
      immediate: true,
    },
  },
  data() {
    return {
      operator_list: [
        { label: "Lt", value: "Lt" },
        { label: "Gt", value: "Gt" },
        { label: "Exists", value: "Exists" },
        { label: "DoesNotExist", value: "DoesNotExist" },
        { label: "In", value: "In" },
        { label: "NotIn", value: "NotIn" },
      ],
      priority: "Preferred",
    }
  },
  methods: {
    handleMatchAdd() {
      var item = {
        key: "",
        operator: "",
        value: "",
      }
      this.matchTables.push(item)
    },
    handleMatchDelete(index) {
      this.matchTables.splice(index, 1)
    },
  },
}
</script>