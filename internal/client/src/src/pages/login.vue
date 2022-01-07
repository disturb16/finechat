<template>
  <div class="container row align-items-start">
    <forma class="col-md-6 offset-md-3">
      <div class="mb-4">
        <label for="email" class="form-label">Email address</label>
        <input type="email" class="form-control" id="email"/>
      </div>
      <div class="mb-4">
        <label for="password" class="form-label">Password</label>
        <input type="password" class="form-control" id="password" />
      </div>
      <button type="submit" class="btn btn-primary">Submit</button>
      <router-link class="register-link" to="/register"> Register</router-link>
    </forma>
  </div>
</template>

<style scoped>
.register-link{
    margin-left: 1em;
}
</style>