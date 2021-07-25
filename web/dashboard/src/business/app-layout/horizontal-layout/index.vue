<template>
  <layout>
    <template v-slot:header>
      <horizontal-header />
    </template>
    <template v-slot:footer>
      <horizontal-footer />
    </template>
  </layout>
</template>

<script>
import HorizontalHeader from "./HorizontalHeader"
import HorizontalFooter from "./HorizontalFooter"
import Layout from "@/components/layout"

export default {
  name: "HorizontalLayout",
  components: { Layout, HorizontalHeader, HorizontalFooter },
}
</script>

<style scoped>
.el-footer {
  background-color: #b3c0d1;
  color: #333;
  text-align: center;
  line-height: 60px;
}
</style>