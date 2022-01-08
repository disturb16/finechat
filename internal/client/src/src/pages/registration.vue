<template>
  <div class="container">
    <form class="col-md-6 offset-md-3">
      <div class="mb-4">
        <label for="fname" class="form-label">First Name</label>
        <input
          type="text"
          class="form-control"
          id="fname"
          v-model="firstName"
        />
      </div>

      <div class="mb-4">
        <label for="lname" class="form-label">Last Name</label>
        <input type="text" class="form-control" id="lname" v-model="lastName" />
      </div>

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

      <div class="mb-4">
        <label for="password" class="form-label">Confirm Password</label>
        <input
          type="password"
          class="form-control"
          id="password"
          v-model="confirmPassword"
        />
      </div>

      <button type="submit" class="btn btn-primary" @click.prevent="register">
        Submit
      </button>
    </form>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "RegisterPage",
  data() {
    return {
      firstName: "",
      lastName: "",
      email: "",
      password: "",
      confirmPassword: "",
    };
  },
  methods: {
    async register() {
      if (this.password !== this.confirmPassword) {
        console.error("passwords don't match");
        return;
      }

      const data = {
        first_name: this.firstName,
        last_name: this.lastName,
        email: this.email,
        password: this.password,
      };

      try {
        await axios.post("/api/users", data);
        this.$router.push({ path: "/login" });
      } catch (error) {
        console.error(error);
      }
    },
  },
};
</script>