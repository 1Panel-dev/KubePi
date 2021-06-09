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
      edit_yaml: "YAML",
      back_form: "返回表单",
      add: "增加"
    },
    table: {
      name: "名称",
      created_time: "创建时间",
      status: "状态",
      action: "操作",
      creat_by: "创建者",
    },
    search: {
      quickSearch: "搜索"
    },
    form: {
      name: "名称"
    },
    validate: {
      limit: "长度在 {0} 到 {1} 个字符",
      input: "请输入{0}",
      select: "请选择{0}",
    },
    msg: {
      delete_success: "删除成功"
    },
    confirm_message: {
      delete: "此操作不可撤销, 是否继续?",
      create_success: "创建成功",
      save_success: "保存成功",
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
    namespace: {
      description: "描述"
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