<template>
  <layout-content :header="$t('commons.form.edit')">
    <el-row></el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"

export default {
  name: "ConfigMapEdit",
  components: { LayoutContent },
  props: {},
  data () {
    return {}
  },
  methods: {},
  created () {
  }
}
</script>

<style scoped>

</style>