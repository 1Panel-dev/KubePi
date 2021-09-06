import store from '@/store'

export const checkPermissions = function (...ps) {
    const userResourcePermissions = store.getters && store.getters.permissions
    const isAdmin = store.getters && store.getters.isAdmin
    if (isAdmin) {
        return true
    }
    for (const p of ps) {
        if (userResourcePermissions["*"] || userResourcePermissions[p.resource]) {
            if (userResourcePermissions["*"]) {
                if (userResourcePermissions["*"].includes("*") || userResourcePermissions["*"].includes(p.verb)) {
                    if (ps.indexOf(p) !== ps.length - 1) {
                        continue
                    }
                    return true
                }
            }
            if (userResourcePermissions[p.resource]) {
                if (userResourcePermissions[p.resource].includes(p.verb) || userResourcePermissions[p.resource].includes("*")) {
                    if (ps.indexOf(p) !== ps.length - 1) {
                        continue
                    }
                    return true
                }
            }
        }
        return false
    }
}
