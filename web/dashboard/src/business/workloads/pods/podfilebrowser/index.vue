<template>
  <layout-content :header="name" :back-to="{name: 'Pods'}">
    <div>
      <i class="el-icon-folder-opened" @click="openRoot"></i>
      <span v-for="(v,i) in folders" :key="i" >
        &nbsp;&nbsp;&nbsp;&nbsp;
        <span>/</span>
         &nbsp;&nbsp;&nbsp;&nbsp;
        <el-link class="primary" :disabled="i === (folders.length -1 )" @click="linkTo(v)"> {{v}}</el-link>
      </span>
    </div>
    <complex-table :selects.="selects" :data="data" v-loading="loading">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="80" show-overflow-tooltip fix>
        <template v-slot:default="{row}">
          <span v-if="row.isDir"  class="span-link" @click="toFolder(row.name)"><i class="el-icon-folder"></i> {{row.name}} </span>
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
      folder: "/",
      folders: []
    }
  },
  methods: {
    openRoot() {
      this.listFiles("/",[])
    },
    linkTo(folder) {
      let folderDir = "/"
      let folders = []
      for(const v of this.folders) {
        folders.push(v)
        folderDir = folderDir + v + "/"
        if (folder === v) {
          break
        }
      }
      this.listFiles(folderDir,folders)
    },
    toFolder(folder) {
      const folderDir = this.folder + folder + "/"
      const newFolders = this.folders.concat([folder])
      this.listFiles(folderDir,newFolders)
    },
    listFiles(folder,folders) {
      this.loading = true
      this.fileRequest.path = folder
      listPodFiles(this.fileRequest).then(res => {
        this.data = res.data
        this.folder = folder
        this.folders = folders
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
      containerName: this.$route.query.container
    }
    this.listFiles("/",[])
  }
}
</script>

<style scoped>

</style>
