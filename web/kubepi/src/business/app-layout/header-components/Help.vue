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
    <el-dialog class="ko-dialog" :show-close="false" :visible.sync="aboutDialogVisible" width="560px">
      <div class="aboutBackground about-panel">
        <div class="about-heading">
          <img :src="require('@/assets/KubePi-red.png')" class="about-logo" alt="KubePi Logo">
          <div class="about-version">{{ $t("commons.personal.version") }}: v2.0.0</div>
        </div>
        <p class="about-text">{{ $t("commons.personal.introduction") }}</p>
        <p class="about-text">{{ $t("commons.personal.introduction2") }}</p>
      </div>
      <div class="about-footer">
        <el-row class="about-links">
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

.about-logo {
  width: 132px;
  height: auto;
  max-height: 44px;
  object-fit: contain;
  vertical-align: middle;
}

.el-dropdown-link {
  cursor: pointer;
}

.about-panel {
  padding: 28px 28px 24px;
}

.about-heading {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 22px;
  padding: 16px 18px;
  background: linear-gradient(180deg, #f8fafc 0%, #ffffff 100%);
  border: 1px solid var(--kp-border);
  border-radius: 8px;
  flex-wrap: wrap;
}

.about-version {
  display: inline-block;
  padding: 3px 9px;
  color: var(--kp-primary);
  background-color: var(--kp-primary-soft);
  border: 1px solid #bfdbfe;
  border-radius: 4px;
  font-size: 12px;
  line-height: 20px;
}

.about-text {
  margin: 0 0 12px;
  color: var(--kp-text-secondary);
  line-height: 1.7;
}

.about-text:last-child {
  margin-bottom: 0;
}

.about-footer {
  padding: 14px 24px 16px;
  border-top: 1px solid var(--kp-border);
  background-color: var(--kp-surface-soft);
}

.about-links {
  color: var(--kp-text-secondary);
  text-align: center;
}

.about-links .el-link {
  color: var(--kp-primary);
  font-weight: 500;
}
</style>
