<template>
  <el-dropdown trigger="click" @command="handleCommand">
    <span class="el-dropdown-link">
        <i class="el-icon-user-solid" style="margin-right: 3px"></i>
        <span>{{ nickName }}</span>
        <i class="el-icon-arrow-down el-icon--right"></i>
    </span>
    <el-dropdown-menu slot="dropdown">
      <el-dropdown-item style="text-align: center" command="profile">
        {{ $t("commons.personal.profile") }}
      </el-dropdown-item>


      <el-dropdown-item style="text-align: center" divided command="exit">
        {{ $t("commons.personal.exit") }}
      </el-dropdown-item>
    </el-dropdown-menu>
  </el-dropdown>
</template>


<script>
import {mapGetters} from "vuex"
import {logout} from "@/api/auth"

export default {
  name: "UserSetting",
  computed: {
    ...mapGetters([
      "nickName"
    ])
  },
  methods: {
    handleCommand(command) {
      switch (command) {
        case "exit":
          this.exit()
          break
        case "profile":
          this.$router.push("/profile")
          break
        default:
          this.aboutDialogVisible = true
          break
      }
    },
    exit() {
      logout().then(() => {
        location.reload()
      })
    },
  }
}
</script>

<style scoped>
.el-dropdown-link {
  cursor: pointer;
}
</style>