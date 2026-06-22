<template>
    <span class="auto-refresh">
        <el-select v-model="autorefresh"  @change="changeAutorefresh">
                  <el-option v-for="vs in autorefreshs" :key="vs.value"
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
    }
  },
  computed: {
    autorefreshs () {
      return [
        {value: '-1', label: this.$t("commons.auto_refresh.manual")},
        {value: '5', label: this.$t("commons.auto_refresh.seconds", [5])},
        {value: '10', label: this.$t("commons.auto_refresh.seconds", [10])},
        {value: '30', label: this.$t("commons.auto_refresh.seconds", [30])},
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
.auto-refresh {
  color: var(--kp-text-secondary);
  margin-left: 20px;
}
</style>
