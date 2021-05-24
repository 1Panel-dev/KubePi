import el from "element-ui/lib/locale/lang/zh-CN" // 加载element的内容
import fu from "fit2cloud-ui/src/locale/lang/zh-CN" // 加载fit2cloud的内容


const message = {
  business: {
    dashboard: {
      dashboard: "概览"
    },
    cluster: {
      cluster: "集群"
    },
    workload: {
      workload: "工作量"
    },
    network: {
      network: "网络"
    },
    storage: {
      storage: "存储"
    },
    configuration: {
      configuration: "配置"
    },
    access_control: {
      access_control: "访问控制"
    }
  }
}


export default {
  ...el,
  ...fu,
  ...message
}