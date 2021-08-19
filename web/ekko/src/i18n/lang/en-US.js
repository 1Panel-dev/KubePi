import el from "element-ui/lib/locale/lang/en"
import fu from "fit2cloud-ui/src/locale/lang/en_US" // 加载fit2cloud的内容

const message = {
    commons: {
        message_box: {
            alert: "alter",
            confirm: "confirm",
            prompt: "prompt",
        },
        personal: {
            profile: "profile",
            exit: "logout"
        },
        button: {
            detail: "Detail",
            delete: "Delete",
            skip: "Skip",
            import: "Import",
            create: "Create",
            cancel: "Cancel",
            login: "Login",
            confirm: "Confirm",
            add: "Add",
            edit: "Edit",
            all_select: "All",
            upload: "Upload",
            search: "Search"
        },
        table: {
            name: "name",
            kind: "kind",
            created_time: "created at",
            status: "status",
            action: "action",
            creat_by: "created by",
            built_in: "built in",
            description: "description",
        },
        header: {
            help_doc: "document",
            support: "support",
            guide: "guide",
            guide_text: "感谢选择本产品, 是否立即开始注册您的第一个集群?",
        },
        bool: {
            true: "true",
            false: "false"
        },
        form: {
            select_placeholder: "please select"
        },
        validate: {
            limit: "长度在 {0} 到 {1} 个字符",
            input: "please input {0}",
            select: "please select {0}",
            required: "required",
            email: "please input a valid email",
        },
        msg: {
            create_success: "create success",
            delete_success: "delete success",
            update_success: "update success",
            no_data: "no data",
        },
        confirm_message: {
            delete: "此操作不可撤销, 是否继续?",
            create_success: "create success",
            save_success: "save success",
        },
        login: {
            username_or_email: "username or email",
            password: "password",
            title: "login EKKO",
            welcome: "Welcome back, please enter your user name and password to log in",
            expires: "认证信息已过期，请重新登录",
        },
    },
    business: {
        cluster: {
            cluster: "Cluster",
            namespace: "Namespace",
            scope: "scope",
            version: "version",
            list: "Cluster",
            import: "import cluster",
            edit: "edit",
            nodes: "nodes",
            label: "label",
            description: "description",
            cluster_detail: "Detail",
            connect_setting: "connect setting",
            connect_direction: "connect direction",
            connect_forward: "forward",
            connect_backward: "backward",
            authenticate_setting: "authenticate setting",
            certificate: "certificate",
            authenticate_mode: "authenticate mode",
            expect: "expect",
            management: "management",
            open_dashboard: "Console",
            cluster_version: "version",
            member: "Member",
            role: "Role",
            administrator: "administrator",
            viewer: "viewer",
            custom: "custom",
            rule: "rule",
            config_file: "config file",
            config_content: "config content",
            hidden_cluster: "hide inaccessible clusters",
            namespace_role_form_check_msg: "namespace or role list cannot be empty"
        },
        cluster_role: {
            none: "无",
            cluster_administrator: "集群管理员",
            cluster_viewer: "只读用户",
        },
        user: {
            user_management: "RBAC",
            username: "username",
            nickname: "nickname",
            email: "email",
            user_list: "User list",
            role_list: "Role list",
            user: "User",
            role: "Role",
            template: "template",
            base_on_exists_role: "base on exists role",
            permission: "permission",
            permission_setting: "permission setting",
            password: "password",
            confirm_password: "confirm password",
            old_password: "old password",
            new_password: "new password",
            change_password: "change password",
            resource_name: "resource name",
            please_input_password: "please input password",
            please_input_password_agin: "please input password again",
            password_not_equal: "two passwords are inconsistent"

        }
    },
}

export default {
    ...el,
    ...fu,
    ...message
};