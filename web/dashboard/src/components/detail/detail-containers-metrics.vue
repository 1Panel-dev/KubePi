<template>
  <table width="100%" v-show="currentPrometheusServer!=null">
      <tr>
       <td style="width:150px">Prometheus Sources:</td>
       <td style="width:150px" align="left">
      <el-select v-model="currentPrometheusServerLabel" filterable @change="handlePrometheusSwitch" v-if="currentPrometheusServer">
 
          <el-option v-for="value in prometheusServers" :key="value.namespace+'.'+value.name" :label="value.namespace+'.'+value.name" :value="value.namespace+'.'+value.name"></el-option>

      </el-select>
      </td>
      <td style="width:150px">time range select:</td>
       <td style="width:150px" align="left">
      <el-select v-model="selectedTimeLabel" filterable @change="handleSelectTime" v-if="currentPrometheusServer">
 
          <el-option v-for="value in times" :key="value.secs" :label="value.label" :value="value.label"></el-option>
      
      </el-select>
      </td>
      <td style="width:150px">refresh interval:</td>
       <td style="width:150px" align="left">
      <el-select v-model="refresh_interval" filterable @change="handleSelectRefreshInterval" v-if="currentPrometheusServer">
 
          <el-option v-for="value in refresh_intervals" :key="value.refresh_interval" :label="value.label" :value="value.refresh_interval"></el-option>
      
      </el-select>
      </td>
      <td style="width:150px"><el-button type="primary" size="small" @click="loadMetrics">Refresh</el-button></td>
      <td style="width:(100%-1050px)"></td>

      </tr>
      <tr>
       <td colspan="8">
          <div class="Echarts">
            <div :id="chartId" style="width: 100%;height: 300px"></div>
          </div>
       </td>
      </tr>
     </table>
</template>
<script>

/*检查已安装在k8s集群中prometheus的服务端*/
import {
  listCustomResourceDefinitions,
} from "@/api/customresourcedefinitions"
import {listPrometheuses,getContainerMetricsMemAndFs,getContainerMetricsCpu} from "@/api/crds/monitoring.coreos.com/v1/prometheus"
/*用户检查prometheus使用的service*/
import {listServicesWithNs} from "@/api/services"
import moment from 'moment'
function arrayDuplicateRemovalSort(arr1,arr2){
const arr = arr1.concat(arr2)
return [...new Set( arr.sort( ( a , b )  => a > b ? 1 : -1 ) ) ]
}
export default {
  name: "KoDetailContainerMetrics",
  props: {
    cluster: String,
    containerInfo: Object,
    yamlInfo: Object,
  },
  watch: {
    cluster:{
      async handler(newValue,oldValue) {
        if (newValue!=oldValue) {
           await this.loadPrometheuses()
           this.loadMetrics()
           this.stopTimeTick()
           this.startTimeTick()
        }
      },
      immediate: true,
      deep: true,
    },
    containerInfo: {
      async handler(newValue) {
        if (newValue) {
           await this.loadPrometheuses()
           this.loadMetrics()
           this.stopTimeTick()
           this.startTimeTick()
        }
      },
      immediate: true,
      deep: true,
    },
  },
  data() {
    return {
      loading: false,
      //prometheus服务器列表
      prometheusServers:[],
      currentPrometheusServer:null,
      currentPrometheusServerLabel:null,
      //定时
      intervalId:null,
      chartId: this.cluster+"|"+this.yamlInfo.metadata.namespace +"|"+this.yamlInfo.metadata.name+"|"+ this.containerInfo.name,
      myChart:null,
      //时间跨度选择
      secs: 1800,/*查询时间长度*/
      time_interval: 1,/*坐标间隔*/
      times:[
        {
          label:"Half Hour ago",
          secs:1800,/*半小时*/
          time_interval: 1
        },
        {
          label:"1 Hour ago",
          secs:3600,/*1小时*/
          time_interval: 2
        },
        {
          label:"12 Hour ago",
          secs:3600*12,/*12小时*/
          time_interval: 24
        },
        {
          label:"24 Hour ago",
          secs:3600*24,/*24小时*/
          time_interval: 48
        }
      ],
      selectedTimeLabel:"Half Hour ago",
      //刷新频率
      refresh_interval: 60000,/*一分钟*/
      refresh_intervals:[
        {
          refresh_interval: -1,/*不自动刷新*/
          label:"no auto refresh"
        },
        {
          refresh_interval: 60000,/*1分钟*/
          label:"1 Minutes"
        },
        {
          refresh_interval: 300000,/*5分钟*/
          label:"5 Minutes"
        },
        {
          refresh_interval: 600000,/*10分钟*/
          label:"10 Minutes"
        }
      ]
   }
  },
  async mounted() {
    await this.loadPrometheuses()
  },
  methods: {
    //加载Prometheus服务器列表
    async loadPrometheuses(){
      let prometheusServers=[]
  
      const crds =await listCustomResourceDefinitions(this.cluster)
      //检查是否已安装普罗
      const crd =crds.items.find(item=>item.metadata.name=="prometheuses.monitoring.coreos.com")
      

      let prometheusServers_res={items:[]}
      try{
         if(crd)
          prometheusServers_res=await  listPrometheuses(this.cluster)
      }catch(e){
          console.log(e)
      }

      if(prometheusServers_res.items && prometheusServers_res.items.length>0){
      for(let i=0,s=prometheusServers_res.items.length;i<s;i++){
         const prometheusServerItem=prometheusServers_res.items[i]
         const name=prometheusServerItem.metadata.name
         const namespace=prometheusServerItem.metadata.namespace
         const prefix=prometheusServerItem.spec.routePrefix

         const services_res=await listServicesWithNs(this.cluster,namespace)
         const service=services_res.items.find(item=>item.spec.selector["prometheus"]==name)
         if(service){
           prometheusServers.push({
              name:name,
              namespace:namespace,
              service:service.metadata.name,
              service_port: (service.spec.ports && service.spec.ports.length>0 )?service.spec.ports[0].port : 80,
              prefix:prefix
           })
         }

      }
      this.prometheusServers=prometheusServers
      this.currentPrometheusServer=prometheusServers[0]
      this.currentPrometheusServerLabel=prometheusServers[0].namespace+"."+prometheusServers[0].name
     }
    },
    //切换Prometheus服务器
    handlePrometheusSwitch(value){
      this.currentPrometheusServer=this.prometheusServers.find(item=>value==item.namespace+"."+item.name)
      if(this.currentPrometheusServer){
       this.stopTimeTick()
       this.loadMetrics()
       this.startTimeTick()
      }
    },
    //选择查询时间范围和坐标间隔
    handleSelectTime(value){
      this.selectedTime=this.times.find(item=> value==item.label)
      this.secs=this.selectedTime.secs
      this.time_interval=this.selectedTime.time_interval;
      if(this.currentPrometheusServer && this.selectedTime){
       

       this.stopTimeTick()
       this.loadMetrics()
       this.startTimeTick()
      }
    },
    //选择刷新频率
    handleSelectRefreshInterval(value){
      this.refresh_interval=value
      if(this.currentPrometheusServer && this.selectedTime){
       

       this.stopTimeTick()
       this.loadMetrics()
       this.startTimeTick()
      }
    },
    startTimeTick(){
      if(this.refresh_interval<0){
        return
      }
      if (this.intervalId != null) {
        return;
      }
      // 计时器为空，操作
      this.intervalId = setInterval(() => {
         //console.log("读取指标" + new Date());
         this.loadMetrics()
      }, this.refresh_interval);
    },
    stopTimeTick(){
      if (!this.intervalId) {
        return;
      }
      clearInterval(this.intervalId); //清除计时器
      this.intervalId = null; //设置为null
    },
    loadMetrics(){
      if(!this.currentPrometheusServer){
        return
      } 
      if(this .myChart!=null){
        this .myChart.showLoading({
          text: 'loading',
          color: '#c23531',
          textColor: '#000',
          maskColor: 'black',
          zlevel: 0,
        })
      }
     
    Promise.all([getContainerMetricsMemAndFs(
        this.cluster,
        this.currentPrometheusServer.namespace,
        this.currentPrometheusServer.service,
        this.currentPrometheusServer.prefix,
        this.currentPrometheusServer.service_port,
        this.yamlInfo.metadata.namespace,
        this.yamlInfo.metadata.name,
        this.containerInfo.name,
        this.secs /*半小时*/
      ),getContainerMetricsCpu(
        this.cluster,
        this.currentPrometheusServer.namespace,
        this.currentPrometheusServer.service,
        this.currentPrometheusServer.prefix,
        this.currentPrometheusServer.service_port,
        this.yamlInfo.metadata.namespace,
        this.yamlInfo.metadata.name,
        this.containerInfo.name,
        this.secs /*半小时*/
      )]).then(res=>{
        if(this .myChart!=null){
           this .myChart.hideLoading()
        }
        this.initCharts(res)
      })
    },
    initCharts(res) {
      const echarts = require('echarts');
      if(this .myChart==null)
      this .myChart = echarts.init(document.getElementById(this.chartId));
      //取不到指标
      if(!res|| res.length<2){
        return
      }
      if(! res[0].data.result[0] || !res[0].data.result[0].values){
        return
      }
      if(! res[1].data.result[0] || !res[1].data.result[0].values){
        return
      }
      // 设置图表的配置项
      let option = {
           tooltip: {
              trigger: 'axis', // 显示横坐标值
           },
           legend: {
            data: ['cpu','fs','mem'],
            textStyle: {
              color: '#FFF',
            },
           },
           xAxis: {
             type: 'category',
             boundaryGap: false,
             //min是半小时前，max是当前时间
             data: arrayDuplicateRemovalSort(res[0].data.result[0].values.map(items=>{
               return   moment(items[0]*1000).format('HH:mm')
             }),res[1].data.result.length>0?res[1].data.result[0].values.map(items=>{
               return   moment(items[0]*1000).format('HH:mm')
             }):[]),
             axisLabel: {
               interval:this.time_interval, // 坐标轴刻度间隔
             }
           },
           yAxis: [{
             type: 'value',
             splitLine: {
               show: false, //去除网格线
             },
             axisLine: {
            //y轴线的颜色以及宽度
              show: true,
              lineStyle: {
                color: '#FFFFFF',
                width: 1,
                type: 'solid',
              },
             },
          
             splitLine: {
              show: true,
              lineStyle: {
                 color: 'rgba(255,255,255,0.3)',
              },
             },
             name: 'c', //单位
             nameLocation: 'end', // (单位个也就是在在Y轴的最顶部)
           },
           {
             type: 'value',
             splitLine: {
               show: false, //去除网格线
             },
             axisLine: {
            //y轴线的颜色以及宽度
              show: true,
              lineStyle: {
                color: '#FFFFFF',
                width: 1,
                type: 'solid',
              },
             },
          
             splitLine: {
              show: true,
              lineStyle: {
                 color: 'rgba(255,255,255,0.3)',
              },
             },
             name: 'MiB', //单位
             nameLocation: 'end', // (单位个也就是在在Y轴的最顶部)
          
           },
           {
             type: 'value',
             splitLine: {
               show: false, //去除网格线
             },
             axisLine: {
            //y轴线的颜色以及宽度
              show: true,
              lineStyle: {
                color: '#FFFFFF',
                width: 1,
                type: 'solid',
              },
             },
          
             splitLine: {
              show: true,
              lineStyle: {
                 color: 'rgba(255,255,255,0.3)',
              },
             },
             name: 'MiB', //单位
             nameLocation: 'end', // (单位个也就是在在Y轴的最顶部)
          
           }
          ],
           series: [
           {
                name: "cpu",
                data: res[1].data.result.length>0 ? res[1].data.result[0].values.map(items=>Number(items[1])):[],
                type: 'line',
                tooltip: {
                    valueFormatter: function (value) {
                        return value + ' c';
                    }
                }
                
           },
           {
                name: "fs",
                data: res[0].data.result[0].values.map(items=>Number(items[1])/1024.00/1024.00),
                type: 'line',
                tooltip: {
                    valueFormatter: function (value) {
                        return value + ' MiB';
                    }
                }
           },
           {
                name: "mem",
                data: res[0].data.result[1].values.map(items=>Number(items[1])/1024.00/1024.00),
                type: 'line',
                tooltip: {
                    valueFormatter: function (value) {
                        return value + ' MiB';
                    }
                }
           }
           ]
      };
      this .myChart.setOption(option);
    }
  },
  destroyed(){
    // 在页面销毁后，清除计时器
    this.stopTimeTick();
  }
}
</script>