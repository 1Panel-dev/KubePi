<template>
  <layout-content :header="$t('business.cluster.cluster_detail')" :back-to="{ name: 'Clusters' }">
    <el-menu class="menuClass" router :default-active="$route.path" mode="horizontal">
      <el-menu-item :index="'/clusters/detail/'+name+'/members'">{{ $t('business.cluster.member') }}</el-menu-item>
      <el-menu-item :index="'/clusters/detail/'+name+'/clusterroles'">{{ $t('business.cluster.role') }}</el-menu-item>
    </el-menu>
    <br/>
    <div class="detailClass">
      <router-view></router-view>
    </div>
  </layout-content>

</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"

export default {
  name: "ClusterEdit",
  props: ["name"],
  components: {LayoutContent},
  data() {
    return {}
  },
  methods: {},
}
</script>

<style scoped>
.menuClass {
  position: fixed;
  z-index: 1;
  width: calc(100vw - 330px);
}

.detailClass {
  margin-top: 60px;
  height: calc(100% - 140px);
  overflow: auto;
}
</style>