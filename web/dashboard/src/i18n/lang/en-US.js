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
}

export default {
  ...el,
  ...fu,
  ...message
};
