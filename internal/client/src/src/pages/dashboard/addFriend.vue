<template>
  <div>
    <h3>Keep in touch with your friends</h3>
    <div class="mb-4">
      <label for="email" class="form-label">Email</label>
      <input
        type="email"
        class="form-control"
        id="email"
        v-model="friendEmail"
      />
    </div>

    <button type="button" class="btn btn-primary" @click.prevent="addFriend">
      Add
    </button>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "AddFriend",
  data() {
    return {
      friendEmail: "",
    };
  },
  methods: {
    async addFriend() {
      if (this.friendEmail == "") {
        return;
      }

      const { email } = this.$store.getters.tokenClaims;

      const data = {
        friend_email: this.friendEmail,
      };

      try {
        const url = `/api/users/${email}/friends`;
        await axios.post(url, data);
        this.$store.dispatch("fetchFriends");
        this.$router.push("/");
      } catch (error) {
        console.error(error);
      }
    },
  },
};
</script>