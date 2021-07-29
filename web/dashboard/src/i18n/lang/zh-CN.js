import el from "element-ui/lib/locale/lang/zh-CN" // 加载element的内容
import fu from "fit2cloud-ui/src/locale/lang/zh-CN" // 加载fit2cloud的内容


const message = {
  commons: {
    message_box: {
      alert: "警告",
      confirm: "确认",
      prompt: "提示",
    },
    personal:{
      exit:"退出"
    },
    button: {
      delete: "删除",
      import: "导入",
      create: "创建",
      cancel: "取消",
      login: "登录",
      confirm: "确定",
      yaml: "YAML",
      back_form: "返回表单",
      add: "增加",
      edit: "编辑",
      edit_yaml: "编辑 YAML",
      search: "搜索",
      upload: "从文件读取",
      view_form: "查看表单",
      view_yaml: "查看 YAML",
      download_yaml: "下载 YAML",
      open_shell:"打开 SHELL",
      back_detail: "返回详情",
      submit: "提交",
    },
    table: {
      name: "名称",
      created_time: "创建时间",
      status: "状态",
      action: "操作",
      creat_by: "创建者",
      time: "时间",
      message: "消息",
      lastUpdateTime: "更新时间",
      type: "类型"
    },
    search: {
      quickSearch: "搜索"
    },
    form: {
      name: "名称",
      detail: "详情"
    },
    validate: {
      limit: "长度在 {0} 到 {1} 个字符",
      input: "请输入{0}",
      select: "请选择{0}",
      required_msg: "输入项不能为空",
      number_limit: "请输入正确的数字",
    },
    msg: {
      delete_success: "删除成功",
      create_success: "创建成功",
      update_success: "修改成功",
    },
    confirm_message: {
      delete: "此操作不可撤销, 是否继续?",
      add_init_container: "当前 initContainers 为空，是否现在添加?",
    },
    login: {
      username: "用户名",
      password: "密码",
      title: "登录 EKKO",
      welcome: "欢迎回来，请输入用户名和密码登录",
      expires: "认证信息已过期，请重新登录",
      table: {
        name: "名称",
        created_time: "创建时间"
      },
      search: {
        quickSearch: "搜索"
      },
      form: {
        name: "名称"
      }
    },
  },
  business: {
    common: {
      label: "标签",
      annotation:"注释",
      basic: "基本信息",
      expand:"展开",
      pack_up: "收起",
      system: "系统信息",
      config: "配置信息",
      resource: "资源信息"
    },
    dashboard: {
      dashboard: "概览"
    },
    cluster: {
      cluster: "集群",
      version: "版本",
      list: "集群列表",
      import: "导入集群",
      nodes: "节点",
      api_server_help: "例如: https://172.16.10.100:8443",
      router_help: "装有 kube-proxy 的任意节点的且可以被访问到的 IP 地址",
      label: "标签",
      description: "描述"
    },
    event: {
      reason: "原因",
      type: "类型",
      time: "时间",
      resource: "对象",
      message: "消息"
    },
    node: {
      ready:"准备就绪",
      role: "角色",
      os: "操作系统",
      arch: "系统架构",
      osImage: "操作系统镜像",
      kernel: "操作系统内核",
      container: "容器 runtime 版本",
      kubelet_version: "kubelet 版本",
      kubeProxy_version: "kubeProxy 版本",
      allocation: "分配"
    },
    pod: {
      image: "镜像",
      ready: "准备就绪",
      restart: "重启",
      type: "类别",
      reason: "原因",
      message: "信息",
      lastHeartbeatTime: "最后检测时间",
      lastTransitionTime:  "最后迁移时间",
      size: "大小",
      resource: "资源信息",
      address: "地址"
    },
    namespace: {
      namespace: "命名空间",
      description: "描述"
    },
    workload: {
      workload: "工作量",
      replicas: "副本数",
      container: "容器",
      initContainer: "初始化容器",
      schedule: "调度",
      lastScheduleTime: "最后的调度",
      suspend: "暂停",
      duration: "耗时",
      lastTransitionTime: "最后更新时间",
      restarts: "重启",
      current: "当前调度",
      desired: "期望期望",
    },
    network: {
      network: "网络",
      port: "端口",
      protocol: "协议",
      rule: "规则",
      host: "主机",
      path: "路径",
      path_type: "路径类型",
      service_name: "Service 名称",
      service_port: "Service 端口",
      target: "目标",
      pod_selector: "Pod 选择器"
    },
    storage: {
      storage: "存储",
      accessModes: "访问模式",
      capacity: "容量 Gib",
      reclaimPolicy: "回收策略",
      storageClass: "存储类",
      claim: "要求",
      assignSc: "绑定 StorageClass",
      DirectoryOrCreateLabel: "DirectoryOrCreate: 如果路径不存在则创建空目录，默认权限为0755 ",
      DirectoryLabel: "Directory: 路径必须存在",
      FileOrCreateLabel: "FileOrCreate: 如果路径上什么都不存在，则创建空文件，默认权限为：0644",
      FileLabel: "File: 路径上必须存在的文件 ",
      SocketLabel: "Socket: 路径上必须存在的 UNIX 套接字 ",
      CharDeviceLabel: "CharDevice: 路径上必须存在的字符设备 ",
      BlockDeviceLabel: "BlockDevice: 路径上必须存在的块设备，例: /dev/sda1 ",
    },
    configuration: {
      configuration: "配置",
      data: "数据",
      no_data: "没有数据可显示",
      type: "类型",
      basic_auth: "基本身份认证",
      ssh_auth: "SSH 身份认证",
      tls_auth: "TLS",
      token_auth: "令牌",
      cluster_ip: "集群 IP",
      selector: "选择"
    },
    access_control: {
      access_control: "访问控制",
      resource_helper: "多个resource请用,分隔"
    },
    custom_resource: {
      custom_resource: "自定义资源",
      full_name: "全名",
      namespaced: "有命名空间的",
      version: "版本",
      scope: "范围",
      names: "允许的名称",
      singular: "单数",
      plural: "复数",
      served: "服务",
      storage: "存储",
      status: "状态"
    },
    user: {
      user_management: "用户与权限",
      user_group: "用户组",
      nickname: "昵称",
      email: "邮箱",
      user_list: "用户列表",
      user_group_list: "用户组列表",
      role_list: "角色列表",
      user: "用户",
      role: "角色",
    }
  },
}


export default {
  ...el,
  ...fu,
  ...message
}
