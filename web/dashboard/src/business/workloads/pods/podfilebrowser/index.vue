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
          <span @click="openFileCreate"> {{ $t("business.pod.create_file") }}</span>
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
          <el-link v-else type="info" @click="openFile(row.name)">{{ row.name }}</el-link>
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
    <el-dialog
            :title="$t('business.pod.create_file')"
            :visible.sync="openAddFile"
            width="60%">
      <el-form label-position="top" :model="fileForm" ref="form" :rules="rules">
        <el-form-item :label="$t('commons.table.name')" required prop="name">
          <el-input clearable v-model="fileForm.name"></el-input>
        </el-form-item>
        <el-form-item :label="$t('business.pod.file_content')" prop="content">
          <el-input type="textarea" :autosize="{ minRows: 15, maxRows: 20}"  v-model="fileForm.content"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
      <el-button @click="openAddFile = false">{{ $t("commons.button.cancel") }}</el-button>
      <el-button type="primary" @click="fileCreate">{{ $t("commons.button.confirm") }}</el-button>
      </span>
    </el-dialog>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {createFile, createFolder, delFolder, listPodFiles, openFile} from "@/api/pods"
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
      fileForm: {
        name: "",
        content: ""
      },
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
    openFileCreate() {
      this.openAddFile = true
      this.fileForm = {}
    },
    getPath(name) {
      if (this.folder === "/") {
        return this.folder +  name
      }else {
        return this.folder + "/" + name
      }
    },
    folderCreate () {
      this.$refs["form"].validate((valid) => {
        if (valid) {
          this.fileRequest.path = this.getPath(this.folderForm.name)
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
        this.fileRequest.path = this.getPath(name)
        delFolder(this.fileRequest).then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.delete_success"),
          })
          this.listFiles(this.folder, this.folders)
        })
      })
    },
    fileCreate () {
      this.$refs["form"].validate((valid) => {
        if (valid) {
          this.fileRequest.path = this.getPath(this.fileForm.name)
          this.fileRequest.content = this.fileForm.content
          createFile(this.fileRequest).then(() => {
            this.openAddFile = false
            this.$message({
              type: "success",
              message: this.$t("commons.msg.create_success"),
            })
            this.listFiles(this.folder, this.folders)
          })
        }
      })
    },
    openFile(name) {
      if (this.folder === "/") {
        this.fileRequest.path = this.folder +  name
      }else {
        this.fileRequest.path = this.folder + "/" + name
      }
      this.fileRequest.path = this.getPath(name)
      openFile(this.fileRequest).then((res) => {
        this.openAddFile = true
        this.fileForm.name = name
        this.fileForm.content = res.data
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
