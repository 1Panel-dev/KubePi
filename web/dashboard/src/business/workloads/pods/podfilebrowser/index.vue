<template>
  <layout-content :header="name" :back-to="{name: 'Pods'}">
    <el-dropdown @command="openPage">
      <el-button size="mini" icon="el-icon-plus" type="primary">{{ $t("commons.button.create") }}<i
              class="el-icon-arrow-down el-icon--right"></i></el-button>
      <el-dropdown-menu slot="dropdown">
        <el-dropdown-item command="file_create">
          <span> {{ $t("business.pod.create_file") }}</span>
        </el-dropdown-item>
        <el-dropdown-item command="folder_create">
          <span> {{ $t("business.pod.create_folder") }}</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </el-dropdown>
    &nbsp;&nbsp;&nbsp;
    <el-dropdown @command="openPage">
      <el-button size="mini" icon="el-icon-upload2" type="primary">{{ $t("business.pod.upload") }}<i
              class="el-icon-arrow-down el-icon--right"></i></el-button>
      <el-dropdown-menu slot="dropdown">
        <el-dropdown-item command="file_upload">
          <span> {{ $t("business.pod.upload_file") }}</span>
        </el-dropdown-item>
        <el-dropdown-item command="folder_upload">
          <span> {{ $t("business.pod.upload_folder") }}</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </el-dropdown>
    <div style="margin-top: 20px;height:20px">
      <el-link><i class="el-icon-folder-opened" @click="openRoot"></i></el-link>
      <span v-for="(v,i) in folders" :key="i">
        &nbsp;&nbsp;&nbsp;&nbsp;
        <span>/</span>
         &nbsp;&nbsp;&nbsp;&nbsp;
        <el-link class="primary" :disabled="i === (folders.length -1 )" @click="linkTo(v)"><span
                style="font-size: 17px"> {{ v }}</span> </el-link>
      </span>
    </div>
    <complex-table :selects.="selects" :data="data" v-loading="loading">
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="200" show-overflow-tooltip fix>
        <template v-slot:default="{row}">
          <span v-if="row.isDir" class="span-link" @click="toFolder(row)"><i
                  class="el-icon-folder"></i> {{ row.name }}{{ row.link ? " -> " + row.link : "" }} </span>
          <el-link v-else @click="catFile(row)"><i
                  class="el-icon-tickets"></i>{{ row.name }}{{ row.link ? " -> " + row.link : "" }}
          </el-link>
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
          <el-tooltip placement="bottom" :content="$t('business.pod.download')">
            <el-button circle size="mini" icon="el-icon-download" @click="download(row)">
            </el-button>
          </el-tooltip>
          <el-dropdown style="margin-left: 10px" @command="handleClick($event,row)" :hide-on-click="false">
            <el-button circle icon="el-icon-more" size="mini"/>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item icon="el-icon-edit" command="edit" :disabled="row.isDir">{{ $t("commons.button.edit") }}
              </el-dropdown-item>
              <el-dropdown-item icon="el-icon-edit-outline" command="rename">{{ $t("business.pod.rename") }}
              </el-dropdown-item>
              <el-dropdown-item icon="el-icon-delete" command="delete">{{ $t("commons.button.delete") }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </complex-table>
    <el-dialog
            :title="$t('business.pod.create_folder')"
            :visible.sync="openAddFolder"
            :close-on-click-modal="false"
            width="30%">
      <el-form label-position="top" :model="folderForm" ref="folderForm" :rules="rules">
        <el-form-item :label="$t('commons.table.name')" prop="name">
          <el-input clearable v-model="folderForm.name" :placeholder="$t('business.pod.name_helper')"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
      <el-button @click="openAddFolder = false">{{ $t("commons.button.cancel") }}</el-button>
      <el-button type="primary" @click="folderCreate">{{ $t("commons.button.confirm") }}</el-button>
      </span>
    </el-dialog>
    <el-dialog
            :title="$t('business.pod.rename')+ '   ' +this.renameForm.oldName"
            :visible.sync="openRenamePage"
            :close-on-click-modal="false"
            width="30%">
      <el-form label-position="top" :model="renameForm" ref="renameForm" :rules="rules">
        <el-form-item :label="$t('commons.table.name')" prop="name">
          <el-input clearable v-model="renameForm.name"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
      <el-button @click="openRenamePage = false">{{ $t("commons.button.cancel") }}</el-button>
      <el-button type="primary" @click="rename()">{{ $t("commons.button.confirm") }}</el-button>
      </span>
    </el-dialog>
    <el-dialog
            :title="editFile?$t('business.pod.edit_file'):$t('business.pod.create_file')"
            :visible.sync="openAddFile"
            :before-close="handleFileClose"
            :close-on-click-modal="false"
            width="60%">
      <el-form label-position="top" :model="fileForm" ref="fileForm" :rules="rules">
        <el-form-item :label="$t('commons.table.name')" prop="name">
          <el-input clearable v-model="fileForm.name" :disabled="editFile"
                    :placeholder="$t('business.pod.name_helper')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('business.pod.file_content')" prop="content">
          <el-input type="textarea" :autosize="{ minRows: 15, maxRows: 20}" v-model="fileForm.content"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
      <el-button @click="handleFileClose()">{{ $t("commons.button.cancel") }}</el-button>
      <el-button type="primary" @click="fileCreate">{{ $t("commons.button.confirm") }}</el-button>
      </span>
    </el-dialog>
    <el-dialog
            :title="$t('business.pod.upload')"
            :visible.sync="openUpload"
            :close-on-click-modal="false"
            :before-close="handleUploadClose"
            width="30%">
      <el-upload :on-change="onUploadChange" ref="upload" action="" :auto-upload="false" class="upload-demo"
                 :multiple="true">
        <el-button>{{ $t("business.pod.choose_file") }}</el-button>
        <div slot="tip" class="el-upload__tip">{{ $t("business.pod.upload_tip") }}</div>
      </el-upload>
      <span slot="footer" class="dialog-footer">
        <el-button @click="handleUploadClose" :disabled="uploadLoading">{{ $t("commons.button.cancel") }}</el-button>
        <el-button type="primary" @click="upload" :loading="uploadLoading">{{ $t("commons.button.confirm") }}</el-button>
      </span>
    </el-dialog>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {
  createFile,
  createFolder,
  delFolder,
  listPodFiles,
  openFile,
  renameFile, updateFile, uploadFile
} from "@/api/pods"
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
      editFile: false,
      openRenamePage: false,
      openUpload: false,
      folderForm: {},
      renameForm: {},
      fileForm: {
        name: "",
        content: ""
      },
      rules: {
        name: [Rule.RequiredRule],
      },
      uploadAction: "",
      files: [],
      uploadLoading: false
    }
  },
  methods: {
    openRoot () {
      this.listFiles("/", [])
    },
    handleClick (btn, row) {
      switch (btn) {
        case "edit":
          this.catFile(row)
          break
        case "delete":
          this.folderDelete(row.name)
          break
        case "rename":
          this.openRename(row.name)
          break
      }
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
    toFolder (row) {
      if (!this.checkLink(row)) {
        return
      }
      const folder = row.name
      const folderDir = this.getPath(folder)
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
    openPage (type) {
      switch (type) {
        case "file_create":
          this.openFileCreate()
          break
        case "folder_create":
          this.openFolderCreate()
          break
        case "file_upload":
          this.openUploadPage(false)
          break
        case "folder_upload":
          this.openUploadPage(true)
          break
        default:
          this.openFileCreate()
          break
      }
    },
    openFolderCreate () {
      this.openAddFolder = true
      this.folderForm = {}
      this.$refs["folderForm"].resetFields()
    },
    openFileCreate () {
      this.openAddFile = true
      this.fileForm = {}
    },
    openRename (name) {
      this.openRenamePage = true
      this.renameForm = {
        oldName: name
      }
      this.$refs["renameForm"].resetFields()
    },
    openUploadPage (dir) {
      this.files = []
      this.openUpload = true
      if (dir) {
        this.$nextTick(() => {
          this.$refs.upload.$children[0].$refs.input.webkitdirectory = true
        })
      } else {
        this.$nextTick(() => {
          this.$refs.upload.$children[0].$refs.input.webkitdirectory = false
        })
      }
    },
    handleUploadClose () {
      this.openUpload = false
      this.$refs.upload.clearFiles()
    },
    handleFileClose () {
      this.openAddFile = false
      this.editFile = false
      this.$refs["fileForm"].resetFields()
    },
    getPath (name) {
      if (this.folder === "/") {
        return this.folder + name
      } else {
        return this.folder + "/" + name
      }
    },
    getLinkPath (row) {
      if (row.link !== "") {
        return row.path
      } else {
        return this.getPath(row.name)
      }
    },
    folderCreate () {
      this.$refs["folderForm"].validate((valid) => {
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
      this.$refs["fileForm"].validate((valid) => {
        if (valid) {
          this.fileRequest.path = this.getPath(this.fileForm.name)
          this.fileRequest.content = this.fileForm.content

          if (this.editFile) {
            updateFile(this.fileRequest).then(() => {
              this.openAddFile = false
              this.$message({
                type: "success",
                message: this.$t("commons.msg.update_success"),
              })
              this.listFiles(this.folder, this.folders)
            })
          } else {
            createFile(this.fileRequest).then(() => {
              this.openAddFile = false
              this.$message({
                type: "success",
                message: this.$t("commons.msg.create_success"),
              })
              this.listFiles(this.folder, this.folders)
            })
          }
        }
      })
    },
    catFile (row) {
      if (!this.checkLink(row)) {
        return
      }
      if (row.size < 1024 * 1024 * 10) {
        this.loading = true
        this.fileRequest.path = this.getPath(row.name)
        openFile(this.fileRequest).then((res) => {
          this.openAddFile = true
          this.editFile = true
          this.fileForm.name = row.name
          this.fileForm.content = res.data
        }).finally(() => {
          this.loading = false
        })
      } else {
        this.$confirm(
          this.$t("commons.confirm_message.change_to_download"),
          this.$t("commons.message_box.prompt"), {
            confirmButtonText: this.$t("commons.button.confirm"),
            cancelButtonText: this.$t("commons.button.cancel"),
            type: "info",
          }).then(() => {
          this.download(row)
        })
      }
    },
    rename () {
      this.$refs["renameForm"].validate((valid) => {
        if (valid) {
          this.fileRequest.path = this.getPath(this.renameForm.name)
          this.fileRequest.oldPath = this.getPath(this.renameForm.oldName)
          renameFile(this.fileRequest).then(() => {
            this.openRenamePage = false
            this.listFiles(this.folder, this.folders)
            this.$message({
              type: "success",
              message: this.$t("commons.msg.update_success"),
            })
          })
        }
      })
    },
    download (row) {
      if (!this.checkLink(row)) {
        return
      }
      const url = this.getUrl(row.name)
      //区分文件和目录,目录tar打包下载,文件直接下载(应对部分镜像tar命令异常导致下载不了)
      if(row["mode"].startsWith("d"))
      window.open("/kubepi/api/v1/pod/files/download/folder" + url, "_blank")
        else
      window.open("/kubepi/api/v1/pod/files/download/file" + url, "_blank")
    },
    getUrl (name) {
      this.fileRequest.path = this.getPath(name)
      let url = ""
      const keys = Object.keys(this.fileRequest)
      for (let i = 0; i < keys.length; i++) {
        if (!this.fileRequest[keys[i]]) {
          continue
        }
        if (url) {
          url += `&${keys[i]}=${this.fileRequest[keys[i]]}`
        } else {
          url += `?${keys[i]}=${this.fileRequest[keys[i]]}`
        }
      }
      return url
    },
    upload () {
      this.uploadLoading = true
      const formData = new FormData()
      const files = this.files
      for (let i = 0; i < files.length; i++) {
        formData.append("files", files[i].raw)
      }
      uploadFile(formData, this.fileRequest).then(() => {
        this.listFiles(this.folder, this.folders)
        this.$message({
          type: "success",
          message: this.$t("commons.msg.upload_success"),
        })
        this.handleUploadClose()
      }).finally(() => {
        this.uploadLoading = false
      })
    },
    onUploadChange (file) {
      this.files.push(file)
    },
    checkLink (row) {
      if (row.link !== "") {
        this.$message({
          showClose: true,
          type: "info",
          message: this.$t("business.pod.link_tip"),
        })
        return false
      }
      return true
    },
  },
  created () {
    this.fileRequest = {
      cluster: this.$route.query.cluster,
      namespace: this.namespace,
      podName: this.name,
      containerName: this.$route.query.container
    }
    this.listFiles("/", [])
  },
}
</script>

<style scoped>

</style>
