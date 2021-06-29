<template>
  <div>
    <el-button icon="el-icon-s-home" @click="toHome()"></el-button>
  </div>
</template>

<script>

export default {
  name: "Home",
  data () {
    return {}
  },
  methods: {
    toHome () {
      window.open("/ekko", '_blank');
    }
  }
}
</script>

<style scoped>

</style>