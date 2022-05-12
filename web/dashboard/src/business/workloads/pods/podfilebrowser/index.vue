<template>
  <layout-content :header="name" :back-to="{name: 'Pods'}">
    <el-dropdown>
      <el-button size="mini" icon="el-icon-plus">{{ $t("commons.button.create") }}<i
              class="el-icon-arrow-down el-icon--right"></i></el-button>
      <el-dropdown-menu slot="dropdown">
        <el-dropdown-item>
          <span @click="openFolderCreate"> {{ $t("business.pod.create_folder") }}</span>
        </el-dropdown-item>
        <el-dropdown-item>
          <span @click="openAddFile = true"> {{ $t("business.pod.create_file") }}</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </el-dropdown>
    <div style="margin-top: 10px">
      <i class="el-icon-folder-opened" @click="openRoot"></i>
      <span v-for="(v,i) in folders" :key="i">
        &nbsp;&nbsp;&nbsp;&nbsp;
        <span>/</span>
         &nbsp;&nbsp;&nbsp;&nbsp;
        <el-link class="primary" :disabled="i === (folders.length -1 )" @click="linkTo(v)"> {{ v }}</el-link>
      </span>
    </div>
    <complex-table :selects.="selects" :data="data" v-loading="loading">
      <!--      <el-table-column type="selection" fix></el-table-column>-->
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="80" show-overflow-tooltip fix>
        <template v-slot:default="{row}">
          <span v-if="row.isDir" class="span-link" @click="toFolder(row.name)"><i
                  class="el-icon-folder"></i> {{ row.name }} </span>
          <span v-else>{{ row.name }}</span>
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
      <el-table-column width="90px" :label="$t('commons.table.action')">
        <template v-slot:default="{row}">
          <el-button circle size="mini" icon="el-icon-download">
          </el-button>
          <el-button circle size="mini" icon="el-icon-delete" @click="folderDelete(row.name)">
          </el-button>
        </template>
      </el-table-column>
    </complex-table>
    <el-dialog
            :title="$t('business.pod.create_folder')"
            :visible.sync="openAddFolder"
            width="30%">
      <el-form label-position="top" :model="folderForm" ref="form" :rules="rules">
        <el-form-item :label="$t('commons.table.name')" required prop="name">
          <el-input clearable v-model="folderForm.name"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
      <el-button @click="openAddFolder = false">{{ $t("commons.button.cancel") }}</el-button>
      <el-button type="primary" @click="folderCreate">{{ $t("commons.button.confirm") }}</el-button>
      </span>
    </el-dialog>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {createFolder, delFolder, listPodFiles} from "@/api/pods"
import ComplexTable from "@/components/complex-table"
import Rule from "@/utils/rules"

export default {
  name: "PodFileBrowser",
  components: { ComplexTable, LayoutContent },
  props: {
    name: String,
    namespace: String
  },
  data () {
    return {
      fileRequest: {},
      selects: [],
      data: [],
      loading: false,
      folder: "/",
      folders: [],
      openAddFolder: false,
      openAddFile: false,
      folderForm: {},
      rules: {
        name: [Rule.RequiredRule, Rule.CommonNameRule],
      }
    }
  },
  methods: {
    openRoot () {
      this.listFiles("/", [])
    },
    linkTo (folder) {
      let folderDir = "/"
      let folders = []
      for (const v of this.folders) {
        folders.push(v)
        folderDir = folderDir + v + "/"
        if (folder === v) {
          break
        }
      }
      this.listFiles(folderDir, folders)
    },
    toFolder (folder) {
      const folderDir = this.folder + folder + "/"
      const newFolders = this.folders.concat([folder])
      this.listFiles(folderDir, newFolders)
    },
    listFiles (folder, folders) {
      this.loading = true
      this.fileRequest.path = folder
      listPodFiles(this.fileRequest).then(res => {
        this.data = res.data
        this.folder = folder
        this.folders = folders
      }).finally(() => {
        this.loading = false
      })
    },
    openFolderCreate () {
      this.openAddFolder = true
      this.folderForm = {}
    },
    folderCreate () {
      this.$refs["form"].validate((valid) => {
        if (valid) {
          this.fileRequest.path = this.folder + "/" + this.folderForm.name
          createFolder(this.fileRequest).then(() => {
            this.openAddFolder = false
            this.$message({
              type: "success",
              message: this.$t("commons.msg.create_success"),
            })
            this.listFiles(this.folder, this.folders)
          })
        }
      })
    },
    folderDelete (name) {
      this.$confirm(
        this.$t("commons.confirm_message.delete"),
        this.$t("commons.message_box.prompt"), {
          confirmButtonText: this.$t("commons.button.confirm"),
          cancelButtonText: this.$t("commons.button.cancel"),
          type: "warning",
        }).then(() => {
        this.fileRequest.path = this.folder + "/" + name
        delFolder(this.fileRequest).then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.delete_success"),
          })
          this.listFiles(this.folder, this.folders)
        })
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
    this.listFiles("/", [])
  }
}
</script>

<style scoped>

</style>
