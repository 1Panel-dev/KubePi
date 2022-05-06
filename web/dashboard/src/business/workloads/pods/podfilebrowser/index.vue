<template>
  <layout-content :header="fileRequest.pod">
    <div>

    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {listPodFiles} from "@/api/pods"
export default {
  name: "PodFileBrowser",
  components: { LayoutContent },
  props: {
    name: String,
    namespace: String
  },
  data () {
    return {
      fileRequest : {}
    }
  },
  methods: {
    listFiles() {
      listPodFiles(this.fileRequest).then(res => {
        console.log(res)
      })
    }
  },
  created () {
    this.fileRequest = {
      cluster: this.$route.query.cluster,
      namespace: this.namespace,
      podName: this.name,
      containerName: this.$route.query.container,
    }
    this.listFiles()
  }
}
</script>

<style scoped>

</style>
