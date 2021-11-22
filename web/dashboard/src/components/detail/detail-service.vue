<template>
  <div>
    <h3>{{ $t("business.common.service") }}</h3>
    <el-alert v-if="services.length === 0" type="info">
      <div slot="title">
        <i class="el-icon-info"></i>
        <span> {{ $t('business.network.service_detail_help', [resourceType, name]) }}</span>
      </div>
    </el-alert>
    <complex-table v-if="services.length !== 0" :data="services">
      <el-table-column :label="$t('business.workload.name')" prop="name" />
      <el-table-column :label="$t('business.workload.protocol')" prop="protocol" />
      <el-table-column :label="$t('business.workload.port')" prop="port" />
      <el-table-column :label="$t('business.workload.node_port')" prop="nodePort" />
      <el-table-column :label="$t('business.workload.container_port')" prop="targetPort" />
    </complex-table>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"
import { listServicesWithNs } from "@/api/services"
import { checkPermissions } from "@/utils/permission"

export default {
  name: "KoDetailSvcs",
  components: { ComplexTable },
  props: {
    cluster: String,
    namespace: String,
    name: String,
    resourceType: String,
  },
  data() {
    return {
      loading: false,
      services: [],
    }
  },
  methods: {
    search() {
      this.loading = true
      if (!checkPermissions({ scope: "namespace", apiGroup: "", resource: "services", verb: "list" })) {
        return
      }
      this.services = []
      listServicesWithNs(this.cluster, this.namespace, this.name).then((res) => {
        for (const svc of res.items) {
          if (svc.metadata.name === this.name) {
            if (svc.spec.ports) {
              for (const p of svc.spec.ports) {
                this.services.push(p)
              }
            }
          }
        }
        this.loading = false
      })
    },
  },
  created() {
    this.search()
  },
}
</script>

<style scoped>
</style>
