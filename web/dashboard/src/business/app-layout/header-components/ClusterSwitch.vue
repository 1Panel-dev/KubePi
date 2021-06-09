<template>
  <el-select v-model="menu" placeholder="请选择" >
<!--    <el-option-group-->
<!--            v-for="group in options"-->
<!--            :key="group.label"-->
<!--            :label="group.label">-->
<!--      <el-option-->
<!--              v-for="item in group.options"-->
<!--              :key="item.value"-->
<!--              :label="item.label"-->
<!--              :value="item.value">-->
<!--      </el-option>-->
<!--    </el-option-group>-->
    <el-option-group>
      <el-option value="global" @click.native="changeMenu('global')" >全局</el-option>
    </el-option-group>
    <el-option-group label="集群">
      <el-option value="cluster1"  @click.native="changeMenu('cluster1')">集群1</el-option>
      <el-option value="cluster2">集群2</el-option>
    </el-option-group>
  </el-select>
<!--  <el-dropdown trigger="click">-->
<!--    <span class="el-dropdown-link">-->
<!--      <i class="iconfont iconnamesapce" style="color: #3884c5;margin-right: 3px" :icon="['fas', 'globe']"/>-->
<!--      <span>切换</span>-->
<!--      <i class="el-icon-arrow-down el-icon&#45;&#45;right"></i>-->
<!--    </span>-->
<!--    <el-dropdown-menu slot="dropdown">-->
<!--      <el-dropdown-item>-->
<!--        <span>全局</span>-->
<!--      </el-dropdown-item>-->
<!--      <el-dropdown-item>-->
<!--        <span>集群</span>-->
<!--      </el-dropdown-item>-->
<!--    </el-dropdown-menu>-->
<!--  </el-dropdown>-->
</template>

<script>

export default {
  name: "ClusterSwitch",
  data() {
    return {
      menu: "global"
    }
  },
  methods: {
    changeMenu(type){
      this.$store.commit("SET_MENU", type)
      location.reload()
      console.log(type)
    }
  }
}
</script>

<style scoped>

</style>