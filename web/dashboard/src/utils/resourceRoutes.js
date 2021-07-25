export const mixin = {
  methods: {
    toResource (type, namespace, name) {
      switch (type) {
        case "Pod":
          this.$router.push({
            name: "PodDetail",
            params: { namespace: namespace, name: name },
            query: { yamlShow: false }
          })
          break
        case "Service":
          this.$router.push({
            name: "ServiceDetail",
            params: { namespace: namespace, name: name },
            query: { yamlShow: false }
          })
          break
        default:
          break
      }
    }
  }
}

