import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

function decodeTokenClaims(token) {
  const tokenDecodablePart = token.split(".")[1];
  const decoded = Buffer.from(tokenDecodablePart, "base64").toString();
  return JSON.parse(decoded);
}

const store = new Vuex.Store({
  state: {
    auth: { token: "" },
  },
  mutations: {
    setAuth(state, token) {
      state.auth.token = token;
    },
  },
  getters: {
    tokenClaims: (state) => {
      return decodeTokenClaims(state.auth.token);
    },

    isAuthenticated: (state) => {
      const token = state.auth.token;
      return typeof token != "undefined" && token != "";
    },
  },
});

export default store;
