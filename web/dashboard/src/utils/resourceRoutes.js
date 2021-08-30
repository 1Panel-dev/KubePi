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
        case "Job":
          this.$router.push({
            name: "JobDetail",
            params: { namespace: namespace, name: name },
            query: { yamlShow: false }
          })
          break
        case "HorizontalPodAutoscaler":
          this.$router.push({
            name: "HPADetail",
            params: { namespace: namespace, name: name },
            query: { yamlShow: false }
          })
          break
        case "PersistentVolumeClaim":
          this.$router.push({
            name: "PersistentVolumeClaimDetail",
            params: { namespace: namespace, name: name },
            query: { yamlShow: false }
          })
          break
        case "Endpoints":
          this.$router.push({
            name: "EndpointDetail",
            params: { namespace: namespace, name: name },
            query: { yamlShow: false }
          })
          break
        case "PodDisruptionBudget":
          this.$router.push({
            name: "PDBDetail",
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

