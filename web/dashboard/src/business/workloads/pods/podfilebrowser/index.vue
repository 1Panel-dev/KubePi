<template>
  <layout-content :header="name">
    <complex-table :selects.="selects" :data="data" v-loading="loading">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="80" show-overflow-tooltip fix>
        <template v-slot:default="{row}">
          <span v-if="row.isDir"  class="span-link" @click="changeFolder(row.name)"><i class="el-icon-folder"></i> {{row.name}} </span>
          <span v-else>{{row.name}}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.pod.size')" prop="size">
      </el-table-column>
      <el-table-column :label="$t('business.pod.permission')" prop="mode">
      </el-table-column>
      <el-table-column :label="$t('business.pod.last_update')" prop="modTime">
        <template v-slot:default="{row}">
          {{ row.modTime | dateFormat }}
        </template>
      </el-table-column>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {listPodFiles} from "@/api/pods"
import ComplexTable from "@/components/complex-table"
export default {
  name: "PodFileBrowser",
  components: { ComplexTable, LayoutContent },
  props: {
    name: String,
    namespace: String
  },
  data () {
    return {
      fileRequest : {},
      selects: [],
      data:[],
      loading: false,
      folder: "/"
    }
  },
  methods: {
    changeFolder(folder) {
      this.folder += folder
      this.fileRequest.path = this.folder
      this.listFiles()
    },

    listFiles() {
      this.loading = true
      listPodFiles(this.fileRequest).then(res => {
        this.data = res.data
      }).finally(() => {
        this.loading = false
      })
    }
  },
  created () {
    this.fileRequest = {
      cluster: this.$route.query.cluster,
      namespace: this.namespace,
      podName: this.name,
      containerName: this.$route.query.container,
      path: this.folder
    }
    this.listFiles()
  }
}
</script>

<style scoped>

</style>
