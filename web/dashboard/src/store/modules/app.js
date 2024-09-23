const get = () => {
  return localStorage.getItem('sidebarStatus')
}
const set = value => {
  localStorage.setItem('sidebarStatus', value)
}
const state = {
  sidebar: {
    opened: get() ? !!+get() : true
  },
  device: 'desktop',
  /*定时刷新频率，单位为秒，-1为手动刷新*/
  autorefresh: '-1'
}

const mutations = {
  TOGGLE_SIDEBAR: state => {
    state.sidebar.opened = !state.sidebar.opened
    if (state.sidebar.opened) {
      set(1)
    } else {
      set(0)
    }
  },
  OPEN_SIDEBAR: (state) => {
    set('sidebarStatus', 1)
    state.sidebar.opened = true
  },
  CLOSE_SIDEBAR: (state) => {
    set('sidebarStatus', 0)
    state.sidebar.opened = false
  },
  SET_AUTOREFRESH: (state, value) => {
    localStorage.setItem('autorefresh',value)
    state.autorefresh = value
  }
}

const actions = {
  toggleSideBar({commit}) {
    commit('TOGGLE_SIDEBAR')
  },
  openSideBar({commit}) {
    commit('OPEN_SIDEBAR')
  },
  closeSideBar({commit}) {
    commit('CLOSE_SIDEBAR')
  },
  setAutorefresh({commit},value) {
    commit('SET_AUTOREFRESH',value)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
