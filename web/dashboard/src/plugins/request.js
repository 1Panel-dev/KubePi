import axios from "axios"
import {$alert, $error} from "./message"
import store from "@/store"
import i18n, {getLanguage} from "@/i18n"

const instance = axios.create({
    baseURL: "/kubepi", // url = base url + request url
    withCredentials: true,
    timeout: 60000 // request timeout, default 1 min
})


instance.interceptors.request.use(
    config => {
        config.headers["lang"] = getLanguage()
        let url = config.url
        const ns = sessionStorage.getItem("namespace")
        if (!config.params) {
            if (ns) {
                config.params = {
                    "namespace": ns
                }
            }
        } else {
            if (ns) {
                config.params["namespace"] = ns
            }
        }
        if (config.params) {
            url += "?"
            let keys = Object.keys(config.params)
            for (let key of keys) {
                url += `${key}=${config.params[key]}&`
            }
            //处理工作空间
            url = url.substring(0, url.length - 1)
            config.params = {}
        }
        config.url = url
        return config
    },
    error => {
        return Promise.reject(error)
    }
)

const checkAuth = response => {
    // 请根据实际需求修改
    if (response.status === 401) {
        let message = i18n.t("commons.login.expires")
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
        location.href = "/kubepi/403"
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
        msg = error.response.data.message || error.response.data
    } else {
        msg = error.message
    }
    if (error.config.url.indexOf("metrics.k8s.io") < 0 ) {
      if (!(error.response.data?.details?.kind === 'secrets' && error.response.data?.code === 404))  {
        $error(msg)
      }
    }
    return Promise.reject(error)
})

export const request = instance

/* 简化请求方法，统一处理返回结果，并增加loading处理，这里以{success,data,message}格式的返回值为例，具体项目根据实际需求修改 */
const promise = (request, loading = {}) => {
    return new Promise((resolve, reject) => {
        loading.status = true
        request.then(response => {
            resolve(response.data)
            loading.status = false
        }).catch(error => {
            reject(error)
            loading.status = false
        })
    })
}

export const get = (url, data, loading) => {
    return promise(request({url: url, method: "get", params: data}), loading)
}

export const post = (url, data, loading) => {
    return promise(request({url: url, method: "post", data}), loading)
}

export const put = (url, data, loading) => {
    return promise(request({url: url, method: "put", data}), loading)
}

export const del = (url, loading) => {
    return promise(request({url: url, method: "delete"}), loading)
}

export const postFile= (url,data, params,loading) => {
  return promise(request({url: url, method: "post",data: data,params:params, headers:{"Content-Type":"application/json"},timeout:600000}), loading)
}

export const patch = (url, data, headers, loading) => {
    if (headers) {
        headers["Content-type"] = "application/merge-patch+json"
        return promise(request({url: url, headers: headers, method: "patch", data}), loading)
    }
    return promise(request({url: url, method: "patch", data}), loading)
}

export default {
    install(Vue) {
        Vue.prototype.$get = get
        Vue.prototype.$post = post
        Vue.prototype.$put = put
        Vue.prototype.$delete = del
        Vue.prototype.$request = request
    }
}
