import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";

Vue.use(Vuex);

function decodeTokenClaims(token) {
  if (!token || token === "") {
    return {};
  }
  const tokenDecodablePart = token.split(".")[1];
  const decoded = Buffer.from(tokenDecodablePart, "base64").toString();
  return JSON.parse(decoded);
}

const store = new Vuex.Store({
  state: {
    auth: { token: "" },
    chatRooms: [],
    friends: [],
  },
  mutations: {
    setAuth(state, token) {
      state.auth.token = token;
    },
    setChatRooms(state, chatRooms) {
      state.chatRooms = chatRooms;
    },
  },
  actions: {
    async fetchChatRooms({ commit, state }) {
      const { email } = decodeTokenClaims(state.auth.token);

      try {
        const url = `/api/users/${email}/chatrooms`;
        const reqConfig = { headers: { Authorization: state.auth.token } };
        const response = await axios.get(url, reqConfig);
        commit("setChatRooms", response.data);
      } catch (error) {
        console.error(error);
      }
    },

    async fetchFriends({ state }) {
      const { email } = decodeTokenClaims(state.auth.token);

      try {
        const url = `/api/users/${email}/friends`;
        const reqConfig = { headers: { Authorization: state.auth.token } };
        const response = await axios.get(url, reqConfig);
        state.friends = response.data;
      } catch (error) {
        console.error(error);
      }
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
