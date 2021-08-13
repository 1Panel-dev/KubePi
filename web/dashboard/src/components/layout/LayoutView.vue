<template>
  <main class="view-container" style="overflow:auto">
    <transition name="el-fade-in" mode="out-in">
      <keep-alive :include="caches">
        <router-view></router-view>
      </keep-alive>
    </transition>
  </main>
</template>

<script>
export default {
  name: "LayoutView",
  data() {
    return {
      caches: []
    }
  },
  computed: {
    key() {
      return this.$route.path
    },
  },
  beforeUpdate() {
    let cache = this.$route.meta?.cache
    let name = this.$route.name
    if (cache && name && !this.caches.includes(name)) {
      this.caches.push(name)
    }
  }
}
</script>

<style scoped>

</style>