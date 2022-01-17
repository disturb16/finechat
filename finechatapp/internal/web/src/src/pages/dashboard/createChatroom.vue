<template>
  <div>
    <h1>New Chatroom</h1>
    <form>
      <div class="mb-4">
        <label for="name" class="form-label">Name</label>
        <input
          type="text"
          class="form-control"
          id="name"
          v-model="chatroomName"
        />
      </div>

      <button
        type="submit"
        class="btn btn-primary"
        @click.prevent="createChatRoom"
      >
        Save
      </button>
    </form>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "CreateChatroom",
  data() {
    return {
      chatroomName: "",
    };
  },
  methods: {
    async createChatRoom() {
      if (this.chatroomName == "") {
        return;
      }

      const { email } = this.$store.getters.tokenClaims;

      const data = {
        email,
        name: this.chatroomName,
      };

      try {
        const reqConfig = {
          headers: {
            Authorization: this.$store.state.auth.token,
          },
        };

        await axios.post("/api/chatrooms", data, reqConfig);
        this.$store.dispatch("fetchChatRooms");
        this.$router.push("/");
      } catch (error) {
        console.error(error);
      }
    },
  },
};
</script>
