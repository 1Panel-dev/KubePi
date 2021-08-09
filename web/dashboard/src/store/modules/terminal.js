const state = {
  terminals: [],
};

const mutations = {
  TERMINALS: (state, value) => {
    state.terminals = value;
  },
};

export default {
  namespaced: true,
  state,
  mutations,
};
