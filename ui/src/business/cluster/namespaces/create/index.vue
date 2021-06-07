<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'Namespaces' }">
    <el-row>
      <div style="float: right">
        <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
        <el-button v-loading="loading" @click="onSubmit" type="primary">
          {{$t("commons.button.create") }}
        </el-button>
      </div>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"

export default {
  name: "NamespaceCreate",
  components: { LayoutContent },
  data () {
    return {
      form:{

      }
    }
  },
  methods: {
    onCancel () {
      this.$router.push({ name: "Namespaces" })
    },
    onSubmit () {

    }
  }
}
</script>

<style scoped>

</style>