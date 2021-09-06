<template>
  <el-button circle  type="text" :class="['sidebar-toggle-button', icon]" @click="toggle"></el-button>
</template>

<script>
import {mapGetters} from "vuex";

export default {
  name: "SidebarToggleButton",
  methods: {
    toggle() {
      this.$store.dispatch('app/toggleSideBar');
    }
  },
  computed: {
    ...mapGetters([
      'sidebar'
    ]),
    icon() {
      return this.sidebar.opened ? "el-icon-arrow-left" : "el-icon-arrow-right"
    },
  }
}
</script>

<style lang="scss" scoped>
  .sidebar-toggle-button.el-button {
    font-size: 18px;
    background-color: #1c1e20;
    color: #3884c5;
  }

</style>
