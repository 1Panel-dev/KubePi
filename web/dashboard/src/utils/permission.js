import store from '@/store'

export const checkPermissions = function (p) {
    const userClusterRoles = store.getters && store.getters.clusterRoles
    for (const clusterRole of userClusterRoles) {
        const scope = clusterRole.metadata.labels["kubeoperator.io/role-type"]
        if (clusterRole.rules.length > 0) {
            for (const rule of clusterRole.rules) {
                if (((scope === p.scope || scope === 'cluster') && rule.apiGroups.includes("*")) || rule.apiGroups.includes(p.apiGroup)) {
                    if (((scope === p.scope || scope === 'cluster') && rule.resources.includes("*")) || rule.resources.includes(p.resource)) {
                        if (((scope === p.scope || scope === 'cluster') && rule.verbs.includes("*")) || rule.verbs.includes(p.verb)) {
                            return true
                        }
                    }
                }
            }
        }
    }
    return false
}
