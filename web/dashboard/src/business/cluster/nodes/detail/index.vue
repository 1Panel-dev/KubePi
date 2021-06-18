<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{ name: 'Nodes' }" v-loading="loading">

  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getNode} from "@/api/nodes"

export default {
  name: "NodeDetail",
  components: { LayoutContent },
  props: {
    name:String,
    cluster: String,
  },
  data () {
    return {
      loading: false,
      item: {}
    }
  },
  methods: {
    getNodeByName() {
      getNode(this.cluster,this.name).then(res =>{
        this.item = res
      })
    }
  },
  created () {
    this.getNodeByName()
  }
}
</script>

<style scoped>

</style>