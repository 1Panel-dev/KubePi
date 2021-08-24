<template>
  <layout-content>

    <el-row>
      <el-col :span="2">
        <el-menu class="menuClass" router :default-active="$route.path"  background-color="#222629">
          <el-menu-item :index="'/clusters/detail/'+name+'/roles/cluster'">{{
              $t('business.cluster.cluster')
            }}
          </el-menu-item>
          <el-menu-item :index="'/clusters/detail/'+name+'/roles/namespace'">{{
              $t('business.cluster.namespace')
            }}
          </el-menu-item>
        </el-menu>
      </el-col>
      <el-col :span="22">
        <router-view></router-view>
      </el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"


export default {
  name: "ClusterRoles",
  props: ["name"],
  components: {LayoutContent},
  data() {
    return {}
  },
  methods: {},
}
</script>

<style scoped>

</style>