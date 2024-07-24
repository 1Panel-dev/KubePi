<template>
  <el-dropdown trigger="click" @command="handleCommand">
    <span class="el-dropdown-link">
      <i class="iconfont iconbangzhu1" style="margin-right: 3px"></i>
      <span>{{ $t('commons.help.help') }}</span>
      <i class="el-icon-arrow-down el-icon--right"></i>
    </span>
    <el-dropdown-menu slot="dropdown">
      <el-dropdown-item command="docs">{{ $t('commons.help.ko_docs') }}</el-dropdown-item>
      <el-dropdown-item command="dev">{{ $t('commons.help.dev_doc') }}</el-dropdown-item>
      <el-dropdown-item class="iconfont iconguanyu" command="about">{{$t("commons.help.about")}}</el-dropdown-item>
    </el-dropdown-menu>
    <el-dialog class="ko-dialog" :show-close="false" :visible.sync="aboutDialogVisible" width="50%">
      <div class="aboutBackground" style="padding: 20px 25px;">
        <img style="margin-left: 0;" :src="require('@/assets/KubePi-red.png')" class="sidebar-logo" alt="Sidebar Logo">
        <p>{{ $t("commons.personal.introduction") }}</p>
        <p>{{ $t("commons.personal.introduction2") }}</p>
        <strong>{{ $t("commons.personal.version") }}: v1.7.0</strong>
      </div>
      <div style="padding:15px 20px;box-shadow:rgb(69 70 70) 0 -14px 24px -12px;">
        <el-row style="color: #ffffff;text-align: center">
          <el-col :span="8">
            <el-link @click="toGithub" class="iconfont iconhuaban88"><span>{{
                $t("commons.personal.project_url")
              }}</span></el-link>
          </el-col>
          <el-col :span="8">
            <el-link @click="toIssue" class="iconfont iconbug">{{ $t("commons.personal.issue") }}</el-link>
          </el-col>
          <el-col :span="8">
            <el-link @click="toGithubStar" class="iconfont iconStarStar">{{ $t("commons.personal.star") }}</el-link>
          </el-col>
        </el-row>
      </div>
    </el-dialog>
  </el-dropdown>
</template>

<script>
export default {
  name: "Help",
  data() {
    return {
      aboutDialogVisible: false,
    }
  },
  methods: {
    handleCommand(command) {
      switch (command) {
        case "docs":
          window.open("https://github.com/1Panel-dev/KubePi/wiki", "_blank")
          break
        case "dev":
          window.open("/kubepi/swagger/index.html", "_blank");
          break
        default:
          this.aboutDialogVisible = true
          break
      }
    },
    toGithub() {
      window.open("https://github.com/1Panel-dev/KubePi", "_blank")
    },
    toIssue() {
      window.open("https://github.com/1Panel-dev/KubePi/issues", "_blank")
    },
    toGithubStar() {
      window.open("https://github.com/1Panel-dev/KubePi", "_blank")
    },
  },
}
</script>

<style lang="scss" scoped>
  @import "~@/styles/business/header-menu.scss";

.sidebar-logo {
  height: $logo-height;
  vertical-align: middle;
}

.el-dropdown-link {
  cursor: pointer;
}
</style>
