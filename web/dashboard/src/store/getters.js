// 根据实际需要修改
const getters = {
    sidebar: state => state.app.sidebar,
    name: state => state.user.name,
    language: state => state.user.language,
    permission_routes: state => state.permission.routes,
    menu: state => state.user.menu,
    roles: state => state.user.roles,
    cluster: state => state.user.cluster,
    clusterRoles: state => state.user.clusterRoles,
    nickName: state => state.user.nickName,
    buttom_height: state => state.terminal.buttomHeight,
    terminals: state => state.terminal.terminals,
    isAdmin: state => state.user.isAdmin,

}
export default getters
