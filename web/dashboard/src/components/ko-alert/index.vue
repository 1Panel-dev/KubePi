<template>
  <div v-if="messages.length > 0">
    <div  v-for="(row, index) in messages" v-bind:key="index">
      <el-alert type="error" :title="m"></el-alert>
    </div>
  </div>
</template>

<script>
export default {
  name: "KoAlert",
  components: {},
  props: {
    messages: Array
  },
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
