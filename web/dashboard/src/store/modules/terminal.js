const state = {
  terminals: [],
  buttomHeight: "0",
};

const mutations = {
  TERMINALS: (state, value) => {
    state.terminals = value;
  },
  CHANGE_BOTTOM_HEIGHT: (state, value) => {
    state.buttomHeight = value
  },
};

export default {
  namespaced: true,
  state,
  mutations,
};
