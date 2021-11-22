<template>
  <layout-content :header="$t('business.cluster.cluster_detail')" :back-to="{ name: 'Clusters' }">
<!--    <el-menu class="menuClass" router :default-active="activeRoute()" mode="horizontal" background-color="#222629">-->
<!--    </el-menu>-->

    <el-menu class="menuClass" router :default-active="activeRoute()" mode="horizontal"  background-color="#222629">
      <el-menu-item :index="membersPath">{{ $t('business.cluster.member') }}</el-menu-item>
      <el-menu-item :index="rolesPath">{{ $t('business.cluster.role') }}</el-menu-item>
      <el-menu-item :index="reposPath">{{ $t('business.cluster.repo') }}</el-menu-item>
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
    return {
      membersPath: `/clusters/detail/${this.name}/members`,
      rolesPath: `/clusters/detail/${this.name}/roles`,
      reposPath: `/clusters/detail/${this.name}/repos`,
    }
  },
  methods: {

    activeRoute() {
      if (this.$route.path.indexOf(this.rolesPath) !== -1) {
        return this.rolesPath
      }
      return this.$route.path
    }

  },
}
</script>

<style scoped>

</style>
