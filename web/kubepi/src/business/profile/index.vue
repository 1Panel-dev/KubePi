<template>
  <layout-content :header="$t('commons.personal.profile')">
    <el-row v-loading="loading">
      <el-col :span="4"><br/></el-col>
      <el-col :span="10">
        <div class="grid-content bg-purple-light">
          <el-form ref="profileForm" :rules="profileRules" :model="profileForm" label-width="150px"
                   label-position="left">

            <el-form-item :label="$t('business.user.username')" prop="name">
              {{ profileForm.name }}
            </el-form-item>

            <el-form-item :label="$t('business.user.nickname')" prop="nickname">
              <span v-if="!editModeOpened">{{ profileForm.nickname }}</span>
              <span v-if="editModeOpened"><el-input v-model="profileForm.nickname"></el-input></span>
            </el-form-item>

            <el-form-item :label="$t('business.user.email')" prop="email">
              <span v-if="!editModeOpened">{{ profileForm.email }}</span>
              <span v-if="editModeOpened"><el-input v-model="profileForm.email"></el-input></span>
            </el-form-item>

            <el-form-item>
              <div style="float: right">
                <el-button v-if="!editModeOpened" @click="onEditProfile">{{ $t("commons.button.edit") }}</el-button>
                <el-button v-if="editModeOpened" @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
                <el-button v-if="editModeOpened" type="primary" @click="onConfirm()">{{
                    $t("commons.button.confirm")
                  }}
                </el-button>
              </div>
            </el-form-item>
          </el-form>
        </div>
      </el-col>
      <el-col :span="4"><br/></el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getCurrentUser, updateProfile} from "@/api/auth"
import Rules from "@/utils/rules"

export default {
  name: "UserProfile",
  components: {LayoutContent},
  data() {
    return {
      loading: false,
      isSubmitGoing: false,
      roleOptions: [],
      editModeOpened: false,
      profileForm: {
        name: "",
        nickname: "",
        email: "",
      },
      profileRules: {
        nickname: [
          Rules.RequiredRule
        ],
        email: [
          Rules.RequiredRule,
          Rules.EmailRule
        ],
      },
    }
  },
  methods: {
    onConfirm() {
      if (this.isSubmitGoing) {
        return
      }
      let isFormReady = false
      this.$refs["profileForm"].validate((valid) => {
        if (valid) {
          isFormReady = true
        }
      })
      if (!isFormReady) {
        return
      }
      this.isSubmitGoing = true
      this.loading = true
      updateProfile({"nickName": this.profileForm.nickname, "email": this.profileForm.email}).then(() => {
        this.isSubmitGoing = false
        this.loading = false
        this.editModeOpened = false
        this.$message.success(this.$t('commons.msg.update_success'))
        window.location.reload()
      })

    },
    onEditProfile() {
      this.editModeOpened = true
    },
    onCancel() {
      this.editModeOpened = false
      this.onVueCreated()
    },
    onVueCreated() {
      getCurrentUser().then(data => {
        this.profileForm.name = data.data.name
        this.profileForm.nickname = data.data.nickName
        this.profileForm.email = data.data.email
      })
    }
  },

  created() {
    this.onVueCreated()
  }
}
</script>

<style scoped>
</style>
