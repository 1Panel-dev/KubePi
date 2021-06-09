<template>
  <layout>
    <template v-slot:header>
      <horizontal-header/>
    </template>
  </layout>
</template>

<script>
import HorizontalHeader from "./HorizontalHeader"
import Layout from "@/components/layout"

export default {
  name: "HorizontalLayout",
  components: { Layout, HorizontalHeader }
}
</script>

<style scoped>

</style>