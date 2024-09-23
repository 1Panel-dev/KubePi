<template>
    <span style="color: white;margin-left: 20px">
        <el-select v-model="autorefresh"  @change="changeAutorefresh">
                  <el-option v-for="(vs,index) in autorefreshs" :key="vs.value"
                             :label="vs.label"
                             :value="vs.value"
                  ></el-option>
        </el-select>
    </span>
</template>

<script>
/*定时刷新设置*/
export default {
  name: "autorefresh",
  components: {},
  props: {},
  data () {
    return {
      autorefresh: '-1',
      autorefreshs:[
        {value:'-1',label:"手动刷新"},
        {value:'5',label:"5秒自动刷新"},
        {value:'10',label:"10秒自动刷新"},
        {value:'30',label:"30秒自动刷新"},
      ]
    }
  },
  methods: {
    async changeAutorefresh(){
        await this.$store.dispatch("app/setAutorefresh", this.autorefresh)
    }
  },
  async mounted () {
    this.autorefresh =localStorage.getItem("autorefresh")
    if(!this.autorefresh || this.autorefresh === 'undefined' || this.autorefresh === ''){
      this.autorefresh = '-1'
    }
    await this.$store.dispatch("app/setAutorefresh", this.autorefresh)

  },
  created () {
    
  }
}
</script>

<style scoped>

</style>

