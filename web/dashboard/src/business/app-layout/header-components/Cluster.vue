<template>
    <span class="cluster-name">{{cluster}}</span>
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
.cluster-name {
  color: var(--kp-text-secondary);
  margin-left: 20px;
}
</style>
