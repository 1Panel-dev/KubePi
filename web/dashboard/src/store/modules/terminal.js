const state = {
  terminals: [],
  bottomHeight: "0",
};

const mutations = {
  TERMINALS: (state, value) => {
    state.terminals = value;
  },
  CHANGE_BOTTOM_HEIGHT: (state, value) => {
    state.bottomHeight = value
  },
};

export default {
  namespaced: true,
  state,
  mutations,
};
