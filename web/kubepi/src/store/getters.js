// 根据实际需要修改
const getters = {
    sidebar: state => state.app.sidebar,
    name: state => state.user.name,
    nickName: state => state.user.nickName,
    language: state => state.user.language,
    permission_routes: state => state.permission.routes,
    permissions: state => state.user.permissions,
    isAdmin: state => state.user.isAdmin,
}
export default getters
