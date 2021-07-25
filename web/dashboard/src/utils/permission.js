import store from '@/store'

export const checkPermissions = function (p) {
    const userClusterRoles = store.getters && store.getters.clusterRoles
    for (const clusterRole of userClusterRoles) {
        if (clusterRole.rules.length > 0) {
            for (const rule of clusterRole.rules) {
                if (rule.apiGroups.includes("*") || rule.apiGroups.includes(p.apiGroup)) {
                    if (rule.resources.includes("*") || rule.resources.includes(p.resource)) {
                        if (rule.verbs.includes("*") || rule.verbs.includes(p.verb)) {
                            return true
                        }
                    }
                }
            }
        }
    }
    return false
}
