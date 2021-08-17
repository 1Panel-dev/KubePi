<template>
  <div class="Echarts">
    <div :id="chartData.name" style="width: 100px; height: 90px"></div>
  </div>
</template>

<script>
export default {
  name: "KoCharts",
  props: {
    chartData: Object
  },
  methods: {
    initCharts() {
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
              normal: {
                color: function (p) {
                  if (p.data.value > 0) {
                    return ['#158fd2']
                  }else {
                    return ['#2e3c3e']
                  }
                }
              }
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
        this.initCharts();
      },
      deep:true
    }
  },
  mounted () {
    this.initCharts();
  }
}
</script>

<style scoped>

</style>
