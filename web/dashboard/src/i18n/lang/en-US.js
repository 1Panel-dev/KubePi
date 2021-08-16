import el from "element-ui/lib/locale/lang/en"
import fu from "fit2cloud-ui/src/locale/lang/en_US" // 加载fit2cloud的内容

const message = {
  commons: {
    message_box: {
      alert: "Warning",
      confirm: "Confirm",
      prompt: "Prompt",
    },
    personal: {
      exit: "Exit"
    },
    button: {
      delete: "Delete",
      import: "Import",
      create: "Create",
      cancel: "Cancel",
      login: "Login",
      confirm: "OK",
      yaml: "YAML",
      back_form: "Back to Form",
      add: "Add",
      edit: "Edit",
      edit_yaml: "Edit YAML",
      search: "Search",
      upload: "Read from file",
      view_form: "View Form",
      view_yaml: "View YAML",
      download_yaml: "Download YAML",
      open_shell: "Open SHELL",
      back_detail: "Back to Details",
      submit: "Submit",
      show_diff: "Diff Contrast",
      continue_edit: "Continue editing",
      terminal: "Terminal",
      logs: "Log",
    },
    table: {
      name: "Name",
      created_time: "Created Time",
      status: "Status",
      action: "Operation",
      creat_by: "Creator",
      time: "Time",
      message: "Message",
      lastUpdateTime: "Update Time",
      type: "Type"
    },
    search: {
      quickSearch: "Search"
    },
    form: {
      name: "Name",
      detail: "Details",
      setting: "Setting",
      parameters: "Parameters"
    },
    validate: {
      limit: "The length is between {0} and {1} characters",
      input: "Please enter {0}",
      select: "Please select {0}",
      required_input_msg: "The input item cannot be empty",
      required_select_msg: "Select item cannot be empty",
      number_limit: "Please enter the correct number",
      name_rules: "Only lowercase English, numbers and-are supported",
      params_not_complete: "Parameter is not complete",
    },
    msg: {
      delete_success: "Delete successfully",
      create_success: "Created successfully",
      update_success: "Modified successfully",
    },
    confirm_message: {
      delete: "This operation cannot be undone, do you want to continue?",
      add_init_container: "The current initContainers is empty, do you want to add it now?",
    },
    login: {
      username: "Username",
      password: "Password",
      title: "Log in to EKKO",
      welcome: "Welcome back, please enter your username and password to log in",
      expires: "The authentication information has expired, please log in again",
      table: {
        name: "Name",
        created_time: "Created time"
      },
      search: {
        quickSearch: "Search"
      },
      form: {
        name: "Name"
      }
    },
  },
  business: {
    common: {
      label: "Label",
      annotation: "Annotation",
      basic: "Basic Information",
      expand: "Expand",
      pack_up: "Pack Up",
      system: "System Information",
      config: "Configuration Information",
      resource: "Resource Information"
    },
    dashboard: {
      dashboard: "Dashboard"
    },
    cluster: {
      cluster: "Cluster",
      version: "Version",
      list: "Cluster List",
      import: "Import Cluster",
      nodes: "Node",
      api_server_help: "eg: https://172.16.10.100:8443",
      router_help: "The accessible IP address of any node with kube-proxy installed",
      label: "Label",
      description: "Description"
    },
    event: {
      reason: "Reason",
      type: "Type",
      time: "Time",
      resource: "Resource",
      message: "Message"
    },
    node: {
      ready: "Ready",
      role: "Role",
      os: "OS",
      arch: "Arch",
      osImage: "OS Image",
      kernel: "Kernel",
      container: "Container",
      kubelet_version: "Kubelet Version",
      kubeProxy_version: "KubeProxy Version",
      allocation: "Allocation"
    },
    pod: {
      image: "Image",
      ready: "Ready",
      restart: "Restart",
      type: "Type",
      reason: "Reason",
      message: "Message",
      lastHeartbeatTime: "Last HeartbeatTime",
      lastTransitionTime: "Last TransitionTime",
      size: "Size",
      resource: "Resource",
      address: "Address",
      watch: "Watch",
      lines: "Lines",
      last_20_lines: "Last 20 lines",
      last_100_lines: "Last 100 lines",
      last_200_lines: "Last 200 lines",
      last_500_lines: "Last 500 lines",
    },
    namespace: {
      namespace: "Namespace",
      description: "Description"
    },
    configuration: {
      configuration: "Configuration",
      data: "Data",
      no_data: "No Data",
      type: "Type",
      basic_auth: "Basic Auth",
      ssh_auth: "SSH Auth",
      tls_auth: "TLS",
      token_auth: "Token Auth",
      cluster_ip: "Cluster IP",
      selector: "Selector",
      registry_domain_name: "Registry Domain Name",
      username: "Username",
      password: "Password",
      authentication: "Authentication",
      certificate: "Certificate",
      publicKey: "PublicKey",
      privateKey: "PrivateKey",
      target: "Target",
      metrics: "Metrics",
      min_replicas: "Min Replicas",
      max_replicas: "Max Replicas",
      current_replicas: "Current Replicas",
      target_reference: "Target Reference",
      source: "Source",
      resource_name: "Resource Name",
      quantity: "Quantity",
      metrics_name: "Metrics Name",
      metrics_selector: "Metrics Selector",
      api_version: "API Version",
      kind: "Kind",
      name: "Name"
    },
    access_control: {
      access_control: "Access Control",
      resource_helper: "Use multiple items, separate"
    },
    custom_resource: {
      custom_resource: "Custom Resource",
      full_name: "Full Name",
      namespaced: "Namespaced",
      version: "Version",
      scope: "Scope",
      names: "Names",
      singular: "Singular",
      plural: "Plural",
      served: "Served",
      storage: "Storage",
      status: "Status"
    },
    user: {
      user_management: "User Management",
      user_group: "User Group",
      nickname: "Nickname",
      email: "Email",
      user_list: "User List",
      user_group_list: "User Group List",
      role_list: "Role List",
      user: "User",
      role: "Role",
    }
  }
}

export default {
  ...el,
  ...fu,
  ...message
}
