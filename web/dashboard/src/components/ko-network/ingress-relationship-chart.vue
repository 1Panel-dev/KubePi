<template>
  <div class="Echarts">
    <div :id="chartId" :style="'width:' + this.chartWidth + 'px ;height: 300px'"></div>
  </div>
</template>

<script>
import { getService } from "@/api/services"
import { listPodsWithNsSelector } from "@/api/pods"
export default {
  name: "KoIngressRelationshipChart",
  props: {
    namespace: String,
    cluster: String,
    ingressObject: Object,
  },
  data() {
    return {
      myChart: null,
      chartWidth: window.innerWidth,
      chartId: this.cluster + "|" + this.namespace + "|" + this.ingressObject.metadata.name + "|ingress|relationship-chart",
    }
  },
  methods: {
    //获取service节点
    async getServicesNodes() {
      let linkDataArray = [];
      let nodes = []
      let servicesNodes = []
      const rules = this.ingressObject.spec.rules;
      const serviceObjectList = {}; //用于service节点去重
      for (let i = 0, s = rules.length; i < s; i++) {
        const rule = rules[i];
        if (rule.http && rule.http.paths) {
          for (let j = 0, l = rule.http.paths.length; j < l; j++) {
            const path = rule.http.paths[j];
            if (path.backend && path.backend["service"]) {
              let serviceObject = null;
              try {
                serviceObject = await getService(this.cluster, this.namespace, path.backend["service"].name)

              } catch (e) {
                console.log(e);
              }
              if (serviceObject && !serviceObjectList[serviceObject.metadata.name]) {
                serviceObjectList[serviceObject.metadata.name] = true

                nodes.push({
                  name: "Service " + serviceObject.metadata.namespace + "." + serviceObject.metadata.name,
                  itemStyle: {
                    color: '#FF9900',			// 关系图节点标记的颜色
                  },
                  draggable: true,
                  symbolSize: [50, 50],
                  nodeType: "Service",
                  nodeName: serviceObject.metadata.name,
                  nodeNamespace: serviceObject.metadata.namespace
                })
                linkDataArray.push({
                  source: "Ingress " + this.namespace + "." + this.ingressObject.metadata.name,
                  target: "Service " + serviceObject.metadata.namespace + "." + serviceObject.metadata.name,
                  'symbol': ['circle', 'arrow'],
                  'symbolSize': [5, 5],
                })
                servicesNodes.push(serviceObject)
              }

            }
          }
        }
      }
      return { linkDataArray: linkDataArray, nodes: nodes, servicesNodes: servicesNodes }
    },
    //获取pod节点
    async getPodsNodes(servicesNodes) {
      let linkDataArray = [];
      let nodes = []
      let podsNodes = []
      let podMap = {}
      for (let i = 0, s = servicesNodes.length; i < s; i++) {
        const serviceObject = servicesNodes[i];
        let selectors = []
        for (let key in serviceObject.spec.selector) {
          selectors.push(key + "=" + serviceObject.spec.selector[key])
        }
        if(!selectors || selectors.length==0){
          continue
        }
        const pods = await listPodsWithNsSelector(this.cluster, this.namespace, selectors.join(","), "")

        for (let j = 0, l = pods.items.length; j < l; j++) {
          const podObject = pods.items[j];
          if (!podMap[podObject.metadata.name]) {
            podMap[podObject.metadata.name] = true
            nodes.push({
              name: "Pod " + podObject.metadata.namespace + "." + podObject.metadata.name,
              itemStyle: {
                color: '#00FF00',			// 关系图节点标记的颜色
              },
              draggable: true,
              symbolSize: [50, 50],
              nodeType: "Pod",
              nodeName: podObject.metadata.name,
              nodeNamespace: podObject.metadata.namespace
            })
            linkDataArray.push({

              source: "Service " + serviceObject.metadata.namespace + "." + serviceObject.metadata.name,
              target: "Pod " + podObject.metadata.namespace + "." + podObject.metadata.name,
              'symbol': ['circle', 'arrow'],
              'symbolSize': [5, 5],
            })
            podsNodes.push(podObject)
          }

        }
      }
      return { linkDataArray: linkDataArray, nodes: nodes, podsNodes: podsNodes }
    },
    //展示关系图
    async initCharts() {
      if (document.getElementById(this.chartId) == null) {
        return
      }
      const echarts = require('echarts');
      if (this.myChart == null) {
        this.myChart = echarts.init(document.getElementById(this.chartId));

        window.addEventListener('resize', this.chartResize);
      }
      var charts = {
        nodes: [ // 节点
          {
            name: "Ingress " + this.namespace + "." + this.ingressObject.metadata.name,
            itemStyle: {
              color: '#387DFF',			// 关系图节点标记的颜色
            },
            draggable: true,
            symbolSize: [50, 50],
          }
        ],
        linesData: [
        ]
      }

      //加入service节点
      const serviceResult = await this.getServicesNodes()

      charts.nodes = charts.nodes.concat(serviceResult.nodes)
      charts.linesData = charts.linesData.concat(serviceResult.linkDataArray)


      //加入pod节点
      const podResult = await this.getPodsNodes(serviceResult.servicesNodes)
      charts.nodes = charts.nodes.concat(podResult.nodes)
      charts.linesData = charts.linesData.concat(podResult.linkDataArray)


      var option = {
        backgroundColor: '#ccc',	// 背景颜色
        tooltip: {                  // 提示框的配置
          formatter: function (param) {
            if (param.dataType === 'edge') {
              //return param.data.category + ': ' + param.data.target;
              return param.data.target;
            }
            //return param.data.category + ': ' + param.data.name;
            return param.data.name;
          }
        },

        series: [{
          type: "graph",          // 系列类型:关系图
          top: '10%',             // 图表距离容器顶部的距离
          roam: true,             // 是否开启鼠标缩放和平移漫游。默认不开启。如果只想要开启缩放或者平移，可以设置成 'scale' 或者 'move'。设置成 true 为都开启
          focusNodeAdjacency: true,   // 是否在鼠标移到节点上的时候突出显示节点以及节点的边和邻接节点。[ default: false ]
          force: {                // 力引导布局相关的配置项，力引导布局是模拟弹簧电荷模型在每两个节点之间添加一个斥力，每条边的两个节点之间添加一个引力，每次迭代节点会在各个斥力和引力的作用下移动位置，多次迭代后节点会静止在一个受力平衡的位置，达到整个模型的能量最小化。
            // 力引导布局的结果有良好的对称性和局部聚合性，也比较美观。
            repulsion: 1000,            // [ default: 50 ]节点之间的斥力因子(关系对象之间的距离)。支持设置成数组表达斥力的范围，此时不同大小的值会线性映射到不同的斥力。值越大则斥力越大
            edgeLength: [150, 100]      // [ default: 30 ]边的两个节点之间的距离(关系对象连接线两端对象的距离,会根据关系对象值得大小来判断距离的大小)，
            // 这个距离也会受 repulsion。支持设置成数组表达边长的范围，此时不同大小的值会线性映射到不同的长度。值越小则长度越长。如下示例:
            // 值最大的边长度会趋向于 10，值最小的边长度会趋向于 50      edgeLength: [10, 50]
          },
          layout: "force",            // 图的布局。[ default: 'none' ]
          // 'none' 不采用任何布局，使用节点中提供的 x， y 作为节点的位置。
          // 'circular' 采用环形布局;'force' 采用力引导布局.
          // 标记的图形
          //symbol: "path://M19.300,3.300 L253.300,3.300 C262.136,3.300 269.300,10.463 269.300,19.300 L269.300,21.300 C269.300,30.137 262.136,37.300 253.300,37.300 L19.300,37.300 C10.463,37.300 3.300,30.137 3.300,21.300 L3.300,19.300 C3.300,10.463 10.463,3.300 19.300,3.300 Z",
          symbol: ['rect', 'arrow'],
          lineStyle: {            // 关系边的公用线条样式。其中 lineStyle.color 支持设置为'source'或者'target'特殊值，此时边会自动取源节点或目标节点的颜色作为自己的颜色。
            normal: {
              color: '#000',          // 线的颜色[ default: '#aaa' ]
              width: 1,               // 线宽[ default: 1 ]
              type: 'solid',          // 线的类型[ default: solid ]   'dashed'    'dotted'
              opacity: 0.5,           // 图形透明度。支持从 0 到 1 的数字，为 0 时不绘制该图形。[ default: 0.5 ]
              curveness: 0            // 边的曲度，支持从 0 到 1 的值，值越大曲度越大。[ default: 0 ]
            }
          },
          label: {                // 关系对象上的标签
            normal: {
              show: true,                 // 是否显示标签
              position: "inside",         // 标签位置:'top''left''right''bottom''inside''insideLeft''insideRight''insideTop''insideBottom''insideTopLeft''insideBottomLeft''insideTopRight''insideBottomRight'
              textStyle: {                // 文本样式
                fontSize: 16
              }
            }
          },
          edgeLabel: {                // 连接两个关系对象的线上的标签
            normal: {
              show: false,
              textStyle: {
                fontSize: 14
              }
            }
          },
          data: charts.nodes,
          links: charts.linesData
        }],

        animationEasingUpdate: "quinticInOut",          // 数据更新动画的缓动效果。[ default: cubicOut ]    "quinticInOut"
        animationDurationUpdate: 100                    // 数据更新动画的时长。[ default: 300 ]
      };
      this.myChart.setOption(option);

      const me = this
      //点击事件
      this.myChart.on('click', function (params) {
        if (params.dataType == 'node') {
          if (params.data.nodeType == 'Pod') {
            me.$router.push({
              name: "PodDetail",
              params: { name: params.data.nodeName, namespace: params.data.nodeNamespace },
              query: { yamlShow: false },
            })
          }
          if (params.data.nodeType == 'Service') {
            me.$router.push({
              name: "ServiceDetail",
              params: { name: params.data.nodeName, namespace: params.data.nodeNamespace },
              query: { yamlShow: false },
            })
          }
        }
      });

    }
  },
  watch: {
    cluster: {
      async handler(newValue, oldValue) {
        if (newValue != oldValue) {
          await this.initCharts()
        }
      },
      immediate: true,
      deep: true,
    },
    namespace: {
      async handler(newValue, oldValue) {
        if (newValue != oldValue) {
          await this.initCharts()
        }
      },
      immediate: true,
      deep: true,
    },
    ingressObject: {
      async handler(newValue, oldValue) {
        if (newValue != oldValue) {
          await this.initCharts()
        }
      },
      immediate: true,
      deep: true,
    },
  },

  computed: {
  },
  async mounted() {
    await this.initCharts()
  },
  created() {
  },
  destroyed() {
  }
}
</script>

<style scoped></style>
