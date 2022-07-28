<template>
    <span style="color: white;margin-left: 20px">{{cluster}}</span>
</template>

<script>
export default {
  name: "ClusterName",
  components: {},
  props: {},
  data () {
    return {
      cluster: ""
    }
  },
  methods: {},
  created () {
    this.cluster = this.$store.getters.cluster
  }
}
</script>

<style scoped>

</style>
