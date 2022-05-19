<template>
  <el-dropdown trigger="click" @command="handleCommand">
    <span class="el-dropdown-link">
        <font-awesome-icon class="language-icon" style="margin-right: 3px"
                           :icon="['fas', 'globe']"/>
        <span>{{ languageMap[language] }}</span>
        <i class="el-icon-arrow-down el-icon--right"></i>
    </span>
    <el-dropdown-menu slot="dropdown">
      <el-dropdown-item v-for="(value, key) in languageMap" :key="key" :index="key" :command="key">
        <span>{{ value }}</span>
        <i class="el-icon-check" v-if="key === language"/>
      </el-dropdown-item>
    </el-dropdown-menu>
  </el-dropdown>
</template>

<script>

export default {
  name: "LanguageSwitch",
  data() {
    return {
      languageMap: {
        "zh-CN": "中文(简体)",
        "en-US": "English",
      }
    }
  },
  computed: {
    language() {
      return this.$store.getters.language
    }
  },
  methods: {
    setLanguage(lang) {
      this.$store.dispatch("user/setLanguage", lang).then(() => {
        location.reload();
      })
    },
    handleCommand(command) {
      switch (command) {
        case "zh-CN":
          this.setLanguage("zh-CN")
          break
        case "en-US":
          this.setLanguage("en-US")
          break
        default:
          this.setLanguage("zh-CN")
          break
      }
    }
  }
}
</script>

<style lang="scss">
@import "~@/styles/business/header-menu.scss";

.header-menu {
  .language-icon {
    width: 24px;
  }
}

.el-dropdown-link {
  cursor: pointer;
}

.header-menu-popper {
  .el-icon-check {
    margin-left: 10px;
    color: $--color-primary;
  }
}
</style>
