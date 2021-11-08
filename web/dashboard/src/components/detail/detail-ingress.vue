<template>
  <div>
    <h3>{{ $t("business.network.ingress_router") }}</h3>
    <el-alert v-if="ingress.length === 0" type="info">
      <div slot="title">
        <i class="el-icon-info"></i>
        <span> {{ $t('business.network.ingress_detail_help', [resourceType, name]) }}</span>
      </div>
    </el-alert>
    <div v-for="(item, index) in ingress" :key="index">
      <span>{{ $t('business.network.host') }}: {{ item.host }}</span>
      <complex-table v-if="ingress.length !== 0" :data="item.details">
        <el-table-column :label="$t('business.network.path_type')" prop="pathType" />
        <el-table-column :label="$t('business.network.path')" prop="path" />
        <el-table-column :label="$t('business.workload.name')" prop="backend.service.name" />
        <el-table-column :label="$t('business.workload.port')" prop="backend.service.port.number" />
      </complex-table>
      <br><br>
    </div>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"
import { listIngressWithNs } from "@/api/ingress"
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
      ingress: [],
    }
  },
  methods: {
    search() {
      this.loading = true
      if (!checkPermissions({ scope: "namespace", apiGroup: "", resource: "ingresses", verb: "list" })) {
        return
      }
      this.ingress = []
      listIngressWithNs(this.cluster, this.namespace, this.name).then((res) => {
        for (const ing of res.items) {
          if ((ing.metadata.name === this.name)) {
            if (ing.spec.rules) {
              for (const i of ing.spec.rules) {
                var item = {
                  host: i.host,
                }
                if (i.http) {
                  item.details = i.http.paths
                }
                this.ingress.push(item)
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
