const state = {
  terminals: [],
  buttomHeight: "0",
  terminalHeight: "400",
};

const mutations = {
  TERMINALS: (state, value) => {
    state.terminals = value;
  },
  CHANGE_BOTTOM_HEIGHT: (state, value) => {
    state.buttomHeight = value
  },
  CHANGE_TERMINAL_HEIGHT: (state, value) => {
    state.terminalHeight = value
  }
};

export default {
  namespaced: true,
  state,
  mutations,
};
