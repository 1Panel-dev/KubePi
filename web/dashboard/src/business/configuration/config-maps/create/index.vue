<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{name: 'ConfigMaps'}" v-loading="loading">
    <div class="grid-content bg-purple-light">
      <el-row :gutter="20">
        <el-form label-position="top" :model="form">
          <el-col :span="12">
            <el-form-item :label="$t('commons.table.name')">
              <el-input clearable v-model="form.metadata.name"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.namespace.namespace')">
              <el-select v-model="form.metadata.namespace">
                <el-option label="All Namespaces" value=""></el-option>
                <el-option v-for="namespace in namespaces"
                           :key="namespace.metadata.name"
                           :label="namespace.metadata.name"
                           :value="namespace.metadata.name">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">

          </el-col>
        </el-form>
      </el-row>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {listNamespace} from "@/api/namespaces"

export default {
  name: "ConfigMapCreate",
  components: { LayoutContent },
  props: {
    name: String,
  },
  data () {
    return {
      loading: false,
      form: {
        metadata: {}
      },
      namespaces: []
    }
  },
  methods: {},
  created () {
    listNamespace(this.$store.state.user.cluster).then(res => {
      this.namespaces = res.items
    })
  }
}
</script>

<style scoped>

</style>