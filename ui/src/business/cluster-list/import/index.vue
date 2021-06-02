<template>
  <layout-content :header="$t('business.cluster.import')"  :back-to="{ name: 'ClusterList' }">
    <el-row>
      <el-col :span="4"><br/></el-col>
      <el-col :span="10">
        <div class="grid-content bg-purple-light">
            <el-form :model="form" label-position="left" :rules="rules" label-width="120px">
              <el-form-item  :label="$t('commons.table.name')" prop="name">
                  <el-input style="width: 100%" v-model="form.name" clearable></el-input>
              </el-form-item>
              <el-form-item style="float: right">
                <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
                <el-button v-loading="loading" @click="onSubmit" type="primary">{{ $t("commons.button.create") }}</el-button>
              </el-form-item>
            </el-form>
        </div>
      </el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"

export default {
  name: "ClusterImport",
  components: { LayoutContent },
  data() {
    return {
      form: {
        name: ""
      },
      rules: {},
      loading: false
    }
  },
  methods: {
    onCancel() {
      this.$router.push({ name: "ClusterList" })

    },
    onSubmit() {

    }
  },
}
</script>

<style scoped>

</style>