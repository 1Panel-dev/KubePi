<template>
  <layout-content :header="$t('business.cluster.import')" :back-to="{ name: 'ClusterList' }">
    <el-row>
      <el-col :span="4"><br/></el-col>
      <el-col :span="10">
        <div class="grid-content bg-purple-light">
          <el-form :model="form" label-position="left" :rules="rules" label-width="120px">
            <el-form-item :label="$t('commons.table.name')" prop="name">
              <el-input v-model="form.name" clearable></el-input>
            </el-form-item>
            <el-form-item label="apiServer" prop="apiServer">
              <el-input v-model="form.apiServer" :placeholder="$t('business.cluster.api_server_help')"
                        clearable></el-input>
            </el-form-item>
            <el-form-item label="router" prop="router">
              <el-input v-model="form.router" :placeholder="$t('business.cluster.router_help')" clearable></el-input>
            </el-form-item>
            <el-form-item label="token" prop="token">
              <el-input type="textarea" :autosize="{ minRows: 2}" style="width: 100%" v-model="form.token"
                        clearable></el-input>
            </el-form-item>
            <el-form-item style="float: right">
              <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
              <el-button v-loading="loading" @click="onSubmit" type="primary">{{
                  $t("commons.button.create")
                }}
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {create} from "@/api/clusters"

export default {
  name: "ClusterImport",
  components: { LayoutContent },
  data () {
    return {
      form: {
        name: "",
        apiServer: "",
        router: "",
        token: ""
      },
      rules: {},
      loading: false
    }
  },
  methods: {
    onCancel () {
      this.$router.push({ name: "ClusterList" })

    },
    onSubmit () {
      create(this.form).then(() => {
        this.$message({
          type: "success",
          message: this.$t("创建成功")
        })
      })
    }
  },
}
</script>

<style scoped>

</style>