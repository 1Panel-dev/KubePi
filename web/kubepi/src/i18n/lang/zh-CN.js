import el from "element-ui/lib/locale/lang/zh-CN" // 加载element的内容
import fu from "fit2cloud-ui/src/locale/lang/zh-CN" // 加载fit2cloud的内容


const message = {
    commons: {
        message_box: {
            alert: "警告",
            confirm: "确认",
            prompt: "提示",
        },
        personal: {
            profile: "个人信息",
            exit: "退出",
            project_url: "项目地址",
            issue: "问题反馈",
            talk: "参与讨论",
            star: "点亮 Star",
            version: "版本",
            introduction: "KubePi 是一款简单易用的开源 Kubernetes 可视化管理面板。",
            introduction2: "KubePi 允许管理员导入多个 Kubernetes 集群，并且通过权限控制，将不同 cluster、namespace 的权限分配给指定用户。它允许开发人员管理 Kubernetes 集群中运行的应用程序并对其进行故障排查，供开发人员更好地处理 Kubernetes 集群中的复杂性。",
        },
        help: {
            help: "帮助",
            about: "关于",
            business_support: "商业支持",
            ko_docs: "官方文档",
        },
        button: {
            detail: "详情",
            delete: "删除",
            skip: "跳过",
            import: "导入",
            create: "创建",
            cancel: "取消",
            login: "登录",
            confirm: "确定",
            edit_yaml: "YAML",
            back_form: "返回表单",
            add: "添加",
            edit: "编辑",
            all_select: "全选",
            upload: "上传文件",
            search: "搜索",
            rbac_manage: "授权",
            sync: "同步"
        },
        table: {
            label: "标签",
            name: "名称",
            kind: "类型",
            created_time: "创建时间",
            status: "状态",
            action: "操作",
            creat_by: "创建者",
            built_in: "内置",
            description: "描述",
        },
        header: {
            help_doc: "帮助文档",
            support: "获取支持",
            guide: "引导",
            guide_text: "感谢选择本产品, 是否立即开始注册您的第一个集群?",
        },
        bool: {
            true: "是",
            false: "否"
        },
        search: {
            quickSearch: "搜索"
        },
        form: {
            select_placeholder: "请选择"
        },
        validate: {
            limit: "长度在 {0} 到 {1} 个字符",
            input: "请输入{0}",
            select: "请选择{0}",
            required: "必填项",
            email: "请输入有效的电子邮箱",
            number_limit: "请输入正确的数字",
            password_help: "有效密码:8-30位,英文字母+数字+特殊字符(~!@#$%^&*,可选)",
            name_not_compliant: "该名称不符合命名规范",
        },
        msg: {
            create_success: "创建成功",
            delete_success: "删除成功",
            update_success: "修改成功",
            no_data: "暂无数据",
        },
        confirm_message: {
            delete: "此操作不可撤销, 是否继续?",
            create_success: "创建成功",
            save_success: "保存成功",
        },
        login: {
            username_or_email: "用户名或邮箱",
            password: "密码",
            title: "登录 KubePi",
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
            node: "节点",
            cluster: "集群",
            namespace: "命名空间",
            scope: "作用域",
            version: "版本",
            list: "集群列表",
            import: "导入集群",
            edit: "编辑",
            nodes: "节点",
            label: "标签",
            description: "描述",
            cluster_detail: "集群详情",
            connect_setting: "连接设置",
            connect_direction: "连接方向",
            connect_forward: "正向连接",
            connect_backward: "反向连接",
            authenticate_setting: "认证设置",
            certificate: "证书",
            authenticate_mode: "认证模式",
            expect: "敬请期待",
            management: "管理",
            open_dashboard: "控制台",
            cluster_version: "版本",
            member: "成员",
            role: "角色",
            administrator: "管理员",
            viewer: "只读者",
            custom: "自定义",
            rule: "规则",
            config_file: "kubeconfig 文件",
            config_content: "config 内容",
            hidden_cluster: "隐藏不可访问的集群",
            namespace_role_form_check_msg: "命名空间或者角色列表不能为空",
            api_group: "API 组",
            resource: "资源",
            verb: "动作",
            cluster_role_form_check_msg: "API组,资源,动作不为空",
            user_not_in_cluster: "您不是此集群的成员，无法访问此集群的控制台，请联系管理员添加您为集群成员",
            ready: "正常",
            not_ready: "异常",
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
        group: {
            binding_management: "绑定管理",
            add_user: "添加用户"
        },
        configuration: {
            configuration: "配置"
        },
        access_control: {
            access_control: "访问控制"
        },
        cluster_role: {
            none: "无",
            cluster_administrator: "集群管理员",
            cluster_viewer: "只读用户",
        },
        user: {
            user_management: "用户管理",
            user_group: "用户组",
            username: "用户名",
            nickname: "昵称",
            email: "邮箱",
            user_list: "用户列表",
            user_group_list: "用户组列表",
            role_list: "角色列表",
            user: "用户",
            role: "角色",
            template: "模版",
            base_on_exists_role: "基于已有角色",
            permission: "权限",
            permission_setting: "权限设置",
            password: "密码",
            confirm_password: "确认密码",
            old_password: "原密码",
            new_password: "新密码",
            change_password: "修改密码",
            resource_name: "资源名称",
            please_input_password: "请输入密码",
            please_input_password_agin: "请再次输入密码",
            password_not_equal: "两次输入密码不一致",
            ldap: "LDAP",
            ldap_address: "地址",
            ldap_port: "端口",
            ldap_username: "用户名",
            ldap_password: "密码",
            ldap_filter_dn: "用户过滤 DN",
            ldap_filter_rule: "用户过滤规则",
            ldap_helper: "注意：没有邮箱的用户不会被同步，与本地用户登录名重复的也不会被同步!",
            ldap_sync: "开始同步，稍后请查看用户列表",
            ldap_sync_error: "请先保存",
            type: "类型"
        }
    },
}

const apiObjects = {
    users: "用户",
    roles: "角色",
    clusters: "集群",
}

const apiVerbs = {
    update: "编辑",
    delete: "删除",
    get: "详情",
    list: "列表",
    create: "创建",
    privilege: "特权"
}


const description = {
    i18n_user_administrator: "超级管理员，拥有所有对象的权限",
    i18n_user_manage_cluster: "集群管理员，拥有集群对象的所有权限",
    i18n_user_manage_rbac: "角色与用户管理员，拥有角色和用户对象的所有权限",
    i18n_user_manage_readonly: "只读用户，只拥有所有对象的访问权限",
    i18n_user_common_user: "普通用户，只拥有集群对象访问权限",
    i18n_user_manage_chart: "Chart仓库管理员， 拥有对Chart仓库的所有权限",

    i18n_cluster_owner: "集群拥有者,拥有所有对象权限",
    i18n_cluster_viewer: "集群只读用户,拥有所有对象的只读权限",
    i18n_manage_cluster_rbac: "集群访问控制管理员, 拥有 ClusterRole、ClusterRoleBinding 对象的所有权限",
    i18n_view_cluster_rbac: "集群访问控制只读用户, 拥有 ClusterRole、ClusterRoleBinding 对象的只读权限",
    i18n_manage_cluster_storage: "集群存储管理员,拥有 StorageClass、PersistentVolume 对象的所有权限",
    i18n_view_cluster_storage: "集群存储只读用户,拥有 StorageClass、PersistentVolume 对象的只读权限",
    i18n_manage_namespaces: "命名空间管理员,拥有对 Namespace 对象的所有权限",
    i18n_view_namespaces: "命名空间只读用户,拥有对 Namespace 对象的所有权限",
    i18n_view_events: "集群事件只读用户, 拥有 Events 对象的只读权限",
    i18n_view_nodes: "节点只读用户,拥有 Node 对象的只读权限",

    i18n_manage_nodes: "节点管理员,拥有 Node 对象的只读权限",
    i18n_manage_crd: "自定义资源管理员,拥有 CustomResourceDefinition 对象的所有权限",
    i18n_view_crd: "自定义资源管理员,拥有 CustomResourceDefinition 对象的只读权限",


    i18n_manage_config: "配置管理员, 拥有当前命名空间 ConfigMap、Secret、ResourceQuotas、LimitRanges、HorizontalPodAutoscalers 和 PodDisruptionBudget 对象的所有权限",
    i18n_view_config: "配置只读用户, 拥有当前命名空间 ConfigMap、Secret、ResourceQuotas、LimitRanges、HorizontalPodAutoscalers 和 PodDisruptionBudget 对象的只读权限",
    i18n_namespace_owner: "命名空间拥有者，拥有当前命名空间内的所有对象的所有权限",
    i18n_namespace_viewer: "命名空间只读用户，拥有当前命名空间内的所有对象的只读权限",
    i18n_view_workload: "工作负载只读用户，拥有当前命名空间内 DaemonSet、StatefulSet、Deployment、Job、CronJob和Pod的只读权限",
    i18n_manage_workload: "工作负载管理员,工作负载只读用户，拥有当前命名空间内 DaemonSet、StatefulSet、Deployment、Job、CronJob和Pod的所有权限",
    i18n_manage_storage: "存储管理员，拥有当前命名空间内 PersistentVolumeClaim 对象的所有权限",
    i18n_view_storage: "存储只读用户，拥有当前命名空间内 PersistentVolumeClaim 对象的只读权限",
    i18n_view_service_discovery: "服务发现只读用户,拥有当前命名空间内 Service、Endpoint、Ingress和NetworkPolicy 对象的只读权限",
    i18n_manage_service_discovery: "服务发现管理员,拥有当前命名空间内 Service、Endpoint、Ingress和NetworkPolicy 对象的所有权限",
    i18n_manage_rbac: "命名空间访问控制,拥有当前命名空间内 Role、RoleBinding 和 ServiceAccount 对象的所有权限",
    i18n_view_rbac: "命名空间访问控制 只读用户,拥有当前命名空间内 Role、RoleBinding 和 ServiceAccount 对象的只读权限",

    i18n_manage_appmarket: "应用市场管理员, 拥有对应用市场的所有权限"

}

export default {
    ...el,
    ...fu,
    ...message,
    ...apiObjects,
    ...apiVerbs,
    ...description

}
