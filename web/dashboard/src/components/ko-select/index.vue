<template>
  <el-select v-model="namespaceStr" v-bind="$attrs" v-on="$listeners" :disabled="disabled">
    <el-option v-for="namespace in namespaces"
               :key="namespace"
               :label="namespace"
               :value="namespace"
               >
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
      namespaces: [],
      disabled: false
    }
  },
  methods: {
    initData (namespaces) {
      this.namespaces = namespaces
      if (this.namespace !== "") {
        this.namespaceStr = this.namespace
        return
      }
      this.namespaceStr = namespaces.find(item => item === "default")
      if (this.namespaceStr === undefined || this.namespaceStr === "") {
        this.namespaceStr = namespaces[0]
      }
      this.$emit("update:namespace", this.namespaceStr)
      this.$emit("change", this.namespaceStr)
    },
  },
  watch: {
    namespaceStr: function (newValue) {
      this.$emit("update:namespace", newValue)
    }
  },
  created () {
    if (!sessionStorage.getItem("namespace")) {
      const cluster = this.$route.query.cluster
      getNamespaces(cluster).then(res => {
        this.initData(res.data)
      })
      this.disabled = false
    } else {
      this.disabled = true
      this.initData([sessionStorage.getItem("namespace")])
    }
  }
}
</script>

<style scoped>

</style>
