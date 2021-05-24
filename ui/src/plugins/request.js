import axios from "axios"
import {$alert, $error} from "./message"
import store from "@/store"
import i18n, {getLanguage} from "@/i18n"

const instance = axios.create({
  baseURL: "", // url = base url + request url
  withCredentials: true,
  timeout: 60000 // request timeout, default 1 min
})


instance.interceptors.request.use(
  config => {
    config.headers["lang"] = getLanguage()
    return config
  },
  error => {
    console.log(error) // for debug
    return Promise.reject(error)
  }
)

const checkAuth = response => {
  // 请根据实际需求修改
  if (response.headers["authentication-status"] === "invalid" || response.status === 401) {
    let message = i18n.t("login.expires")
    $alert(message, () => {
      store.dispatch("user/logout").then(() => {
        location.reload()
      })
    })
  }
}

const checkPermission = response => {
  // 请根据实际需求修改
  if (response.status === 403) {
    location.href = "/403"
  }
}

// 请根据实际需求修改
instance.interceptors.response.use(response => {
  checkAuth(response)
  return response
}, error => {
  let msg
  if (error.response) {
    checkAuth(error.response)
    checkPermission(error.response)
    msg = error.response.data.msg || error.response.data
  } else {
    console.log("error: " + error) // for debug
    msg = error.message
  }
  $error(msg)
  return Promise.reject(error)
})

export const request = instance

/* 简化请求方法，统一处理返回结果，并增加loading处理，这里以{success,data,message}格式的返回值为例，具体项目根据实际需求修改 */
const promise = (request, loading = {}) => {
  return new Promise((resolve, reject) => {
    loading.status = true
    request.then(response => {
      // if (response.data.success) {
      resolve(response.data)
      // }
      // else {
      //   reject(response.data)
      // }
      loading.status = false
    }).catch(error => {
      reject(error)
      loading.status = false
    })
  })
}

export const get = (url, data, loading) => {
  return promise(request({ url: url, method: "get", params: data }), loading)
}

export const post = (url, data, loading) => {
  return promise(request({ url: url, method: "post", data }), loading)
}

export const put = (url, data, loading) => {
  return promise(request({ url: url, method: "put", data }), loading)
}

export const del = (url, loading) => {
  return promise(request({ url: url, method: "delete" }), loading)
}

export const patch = (url, data, headers, loading) => {
  if (headers) {
    return promise(request({ url: url, headers: headers, method: "patch", data }), loading)
  }
  return promise(request({ url: url, method: "patch", data }), loading)
}

export default {
  install (Vue) {
    Vue.prototype.$get = get
    Vue.prototype.$post = post
    Vue.prototype.$put = put
    Vue.prototype.$delete = del
    Vue.prototype.$request = request
  }
}
