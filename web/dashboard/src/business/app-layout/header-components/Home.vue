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
      window.open("/kubepi", '_self');
    }
  }
}
</script>

<style scoped>

</style>