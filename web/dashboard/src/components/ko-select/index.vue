<template>
  <el-select v-model="namespaceStr" :disabled="disabled" @change="changeNamespace">
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
      disabled: false,
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
    },
    changeNamespace () {
      this.$emit("changeFunc")
    }
  },
  created () {
    if (!sessionStorage.getItem("namespace")) {
      const cluster = this.$route.query.cluster
      getNamespaces(cluster).then(res => {
        this.initData(res.data)
        this.disabled = false
      })
    } else {
      this.initData([sessionStorage.getItem("namespace")])
      this.disabled = true
    }
  }
}
</script>

<style scoped>

</style>
