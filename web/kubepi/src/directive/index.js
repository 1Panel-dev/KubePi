import Permission from "./permission";

export default {
    install(Vue) {
        Vue.directive('has-permissions', Permission.hasPermissions);
    }
}

