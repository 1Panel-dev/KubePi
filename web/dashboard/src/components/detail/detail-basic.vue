<template>
  <div>
    <h3>{{ $t("business.common.basic") }}</h3>
    <table style="width: 100%" class="myTable">
      <tr>
        <td>{{ $t("commons.table.name") }}</td>
        <td colspan="4">{{ item.metadata.name }}</td>
      </tr>
      <tr v-if="item.metadata.namespace">
        <td>{{ $t("business.namespace.namespace") }}</td>
        <td colspan="4">{{ item.metadata.namespace }}</td>
      </tr>
      <tr v-if="item.metadata.finalizers">
        <td>Finalizers</td>
        <td colspan="4">
          <div v-for="value in item.metadata.finalizers" v-bind:key="value" class="myTag">
            <el-tag type="info" size="small">
              {{ value.length > 100 ? value.substring(0, 100) + "..." : value }}
            </el-tag>
          </div>
        </td>
      </tr>
      <tr v-if="hasPodContainers()">
        <td>{{ $t("business.pod.image") }}</td>
        <td colspan="4">
          <div v-for="(item,index) in containers" v-bind:key="index" class="myTag">
            <el-tag type="info" size="small">
              {{ item.image }}
            </el-tag>
          </div>
        </td>
      </tr>
      <tr v-if="item.kind === 'Namespace'">
        <td>Resource Quotas</td>
        <td colspan="4">
          <div v-for="(item,index) in resourceQuotas" v-bind:key="index" class="my-tag">
            <el-link @click="toResourceQuota(item)">{{item.metadata.name}}</el-link>
          </div>
        </td>
      </tr>
      <tr v-if="item.kind === 'Namespace'">
        <td>Limit Ranges</td>
        <td colspan="4">
          <div v-for="(item,index) in limitRanges" v-bind:key="index" class="my-tag">
            <el-link @click="toLimitRange(item)">
              {{item.metadata.name}}
            </el-link>
          </div>
        </td>
      </tr>
      <tr>
        <td>{{ $t("commons.table.created_time") }}</td>
        <td colspan="4">{{ item.metadata.creationTimestamp | age }}</td>
      </tr>
      <tr>
        <td>{{ $t("business.common.label") }}</td>
        <td colspan="4">
          <ko-detail-key-value v-if="Object.keys(item.metadata).length > 0" :valueObj="item.metadata.labels"></ko-detail-key-value>
        </td>
      </tr>
      <tr>
        <td>{{ $t("business.common.annotation") }}</td>
        <td colspan="4">
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
import {listResourceQuotaByNamespace} from "@/api/resourcequota"
import {listLimitRangeByNamespace} from "@/api/limitranges"

export default {
  name: "KoDetailBasic",
  components: { KoDetailKeyValue },
  props: {
    item: Object,
    yamlShow: Boolean,
  },
  data() {
    return {
      show: false,
      containers: [],
      limitRanges: [],
      resourceQuotas: [],
      cluster: ""
    }
  },
  methods: {
    showYaml() {
      this.show = !this.show
      this.$emit("update:yamlShow", this.show)
    },
    hasPodContainers() {
      if (this.item.spec?.template?.spec?.containers) {
        this.containers = this.item.spec.template.spec.containers
        return true
      } else if (this.item.spec?.jobTemplate?.spec?.template?.spec?.containers) {
        this.containers = this.item.spec.jobTemplate.spec.template.spec.containers
        return true
      } else {
        return false
      }
    },
    listResources() {
      listResourceQuotaByNamespace(this.cluster,this.item.metadata.name).then(res => {
        this.resourceQuotas = res.items
      })
      listLimitRangeByNamespace(this.cluster,this.item.metadata.name).then(res => {
        this.limitRanges = res.items
      })
    },
    toResourceQuota(row){
      this.$router.push({
        name: "ResourceQuotaDetail",
        params: { namespace: row.metadata.namespace, name: row.metadata.name }
      })
    },
    toLimitRange(row) {
      this.$router.push({
        name: "LimitRangeDetail",
        params: {namespace: row.metadata.namespace, name: row.metadata.name}
      })
    }
  },
  created() {
    this.show = this.yamlShow
    this.cluster = this.$route.query.cluster
    if (this.item.kind === 'Namespace') {
      this.listResources()
    }
  },
}
</script>

<style scoped>
</style>
