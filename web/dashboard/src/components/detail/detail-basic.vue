<template>
  <div>
    <table style="width: 100%" class="myTable">
      <tr>
        <th scope="col" width="30%" align="left">
          <h3>{{ $t("business.common.basic") }}</h3></th>
        <th scope="col"></th>
      </tr>
      <tr>
        <td>{{ $t("commons.table.name") }}</td>
        <td>{{ item.metadata.name }}</td>
      </tr>
      <tr v-if="item.metadata.namespace">
        <td>{{ $t("business.namespace.namespace") }}</td>
        <td>{{ item.metadata.namespace }}</td>
      </tr>
      <tr v-if="item.metadata.finalizers">
        <td>Finalizers</td>
        <td>
          <div v-for="value in item.metadata.finalizers" v-bind:key="value" class="myTag">
            <el-tag type="info" size="small">
              {{ value.length > 100 ? value.substring(0, 100) + "..." : value }}
            </el-tag>
          </div>
        </td>
      </tr>
      <tr>
        <td>{{ $t("commons.table.created_time") }}</td>
        <td>{{ item.metadata.creationTimestamp | age }}</td>
      </tr>
      <tr>
        <td>{{ $t("business.common.label") }}</td>
        <td>
          <ko-detail-key-value v-if="Object.keys(item.metadata).length > 0"
                               :valueObj="item.metadata.labels"></ko-detail-key-value>
        </td>
      </tr>
      <tr>
        <td>{{ $t("business.common.annotation") }}</td>
        <td>
          <ko-detail-key-value v-if="Object.keys(item.metadata).length > 0" :valueObj="item.metadata.annotations"></ko-detail-key-value>
        </td>
      </tr>
    </table>
    <div class="bottom-button">
      <el-button @click="showYaml">{{ $t("commons.button.view_yaml") }}</el-button>
    </div>
  </div>
</template>

<script>
import KoDetailKeyValue from "@/components/detail/detail-key-value"

export default {
  name: "KoDetailBasic",
  components: { KoDetailKeyValue },
  props: {
    item: Object,
    yamlShow: Boolean
  },
  data () {
    return {
      show: false
    }
  },
  methods: {
    showYaml () {
      this.show = !this.show
      this.$emit("update:yamlShow", this.show)
    }
  },
  created () {
    this.show = this.yamlShow
  }
}
</script>

<style scoped>

</style>
