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
              normal: {
                color: function () {
                  const colors = ['#fee82c','#abcd2c','#db2d45','#cd4d89','#1491d1','#f17f43','#108774','#7c4899','#282a88','#008B8B','#90EE90','#8B0000']
                  return [colors[index]]
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
