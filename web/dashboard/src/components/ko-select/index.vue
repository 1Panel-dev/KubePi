<template>
  <el-select v-model="namespaceStr" v-bind="$attrs" v-on="$listeners">
    <el-option v-for="namespace in namespaces"
               :key="namespace"
               :label="namespace"
               :value="namespace">
    </el-option>
  </el-select>
</template>

<script>
import {getNamespaces} from "@/api/auth"

export default {
  name: "KoSelect",
  components: {},
  props: {
    namespace: String
  },
  data () {
    return {
      namespaceStr: "",
      namespaces: []
    }
  },
  methods: {
    initData (namespaces) {
      this.namespaces = namespaces
      this.namespaceStr = namespaces.find(item => item==='default')
      if (this.namespaceStr === undefined || this.namespaceStr === "" ) {
        this.namespaceStr = namespaces[0]
      }
      this.$emit("update:namespace", this.namespaceStr)
    }
  },
  created () {
    if (!sessionStorage.getItem("namespace")) {
      const cluster = this.$route.query.cluster
      getNamespaces(cluster).then(res => {
        this.initData(res.data)
      })
    } else {
      this.initData([sessionStorage.getItem("namespace")])
    }
  }
}
</script>

<style scoped>

</style>
