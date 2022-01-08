import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    auth: { token: "" },
  },
  mutations: {
    setAuth(state, token) {
      state.auth.token = token;
    },
  },
});

export default store;
