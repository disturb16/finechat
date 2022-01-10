<template>
  <div class="container row align-items-start">
    <form class="col-md-6 offset-md-3">
      <div class="mb-4">
        <label for="email" class="form-label">Email address</label>
        <input type="email" class="form-control" id="email" v-model="email" />
      </div>
      <div class="mb-4">
        <label for="password" class="form-label">Password</label>
        <input
          type="password"
          class="form-control"
          id="password"
          v-model="password"
        />
      </div>
      <button type="submit" class="btn btn-primary" @click.prevent="signIn">
        Submit
      </button>
      <router-link class="register-link" to="/register"> Register</router-link>
      <div
        v-if="errorMessage != ''"
        class="command-error alert alert-danger"
        role="alert"
      >
        {{ errorMessage }}
      </div>
    </form>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "login",
  data() {
    return {
      email: "",
      password: "",
      errorMessage: "",
    };
  },

  methods: {
    async signIn() {
      const data = {
        email: this.email,
        password: this.password,
      };

      this.errorMessage = "";

      try {
        const response = await axios.post("/api/users/signin", data);
        this.$store.commit("setAuth", response.data.token);

        this.$router.push("/");
      } catch (error) {
        if (error.response.status === 400 || error.response.status === 401) {
          this.errorMessage = "Invalid user credentials";
        } else {
          this.errorMessage = "Something went wrong";
        }
      }
    },
  },
};
</script>

<style scoped>
.register-link {
  margin-left: 1em;
}
</style>