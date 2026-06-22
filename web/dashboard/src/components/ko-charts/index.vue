<template>
  <div class="Echarts">
    <div :id="chartData.name" style="width: 100px; height: 90px"></div>
  </div>
</template>

<script>
export default {
  name: "KoCharts",
  props: {
    chartData: Object,
    index: Number
  },
  data() {
    return {
    }
  },

  methods: {
    loadColor(index) {
      const colors = ['#2563eb','#16a34a','#d97706','#dc2626','#0891b2','#7c3aed','#db2777','#0f766e','#64748b','#f59e0b','#84cc16','#be123c']
      return colors[index]
    },
    initCharts(index) {
      const echarts = require('echarts');
      let myChart = echarts.init(document.getElementById(this.chartData.name));
      let option = {
        tooltip: {
          trigger: 'item'
        },
        series: [
          {
            type: 'pie',
            radius: ['40%', '70%'],
            avoidLabelOverlap: false,
            label: {
              show: false,
              position: 'center'
            },
            emphasis: {
              label: {
                show: false,
              }
            },
            labelLine: {
              show: false
            },
            tooltip:{
              show: false
            },
            data: this.chartData.data,
            itemStyle: {
                color: this.loadColor(index)
            },

          }
        ]
      };
      myChart.setOption(option);
    }
  },
  watch: {
    chartData: {
      handler() {
        this.initCharts(this.index);
      },
      deep:true
    }
  },
  mounted () {
    this.initCharts(this.index);
  }
}
</script>

<style scoped>

</style>
