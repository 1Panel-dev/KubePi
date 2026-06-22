import axios from "axios"
import {$alert, $error} from "./message"
import store from "@/store"
import i18n, {getLanguage} from "@/i18n"

const instance = axios.create({
    baseURL: "/kubepi", // url = base url + request url
    withCredentials: true,
    timeout: 60000 // request timeout, default 1 min
})

const errorMessageCache = new Map()
const errorMessageInterval = 5000

const getErrorMessage = error => {
    const data = error.response?.data
    if (typeof data === "string" && data) {
        return data
    }
    if (data?.message) {
        return data.message
    }
    if (error.response) {
        return i18n.t("commons.msg.cluster_api_error")
    }
    return error.message || i18n.t("commons.msg.cluster_api_error")
}

const shouldSkipErrorMessage = error => {
    const url = error.config?.url || ""
    const data = error.response?.data
    if (url.indexOf("metrics.k8s.io") >= 0) {
        return true
    }
    return data?.details?.kind === "secrets" && data?.code === 404
}

const shouldShowErrorMessage = message => {
    const key = String(message)
    const now = Date.now()
    const lastShowTime = errorMessageCache.get(key) || 0
    if (now - lastShowTime < errorMessageInterval) {
        return false
    }
    errorMessageCache.set(key, now)
    return true
}


instance.interceptors.request.use(
    config => {
        config.headers["lang"] = getLanguage()
        let url = config.url
        const ns = sessionStorage.getItem("namespace")
        if (!config.params) {
            if (ns) {
                if(
                    url.indexOf("monitoring.coreos.com/v1/prometheuses")==-1
                    &&
                   url.indexOf("api/v1/query_range")==-1
                )
                config.params = {
                    "namespace": ns
                }
            }
        } else {
            if (ns) {
                if(
                    url.indexOf("monitoring.coreos.com/v1/prometheuses")==-1
                   &&
                   url.indexOf("api/v1/query_range")==-1
                )
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
        msg = getErrorMessage(error)
    } else {
        msg = getErrorMessage(error)
    }
    if (!shouldSkipErrorMessage(error) && shouldShowErrorMessage(msg)) {
        $error(msg)
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
export const delWithData = (url,data, loading) => {
    return promise(request({url: url, method: "delete",data:data}), loading)
}

export const postFile= (url,data, params,loading) => {
  return promise(request({url: url, method: "post",data: data,params:params, headers:{"Content-Type":"multipart/form-data"},timeout:600000}), loading)
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
