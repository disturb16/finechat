<template>
  <div class="container">
    <form class="col-md-6 offset-md-3">
      <div class="mb-4">
        <label for="fname" class="form-label">First Name</label>
        <input type="text" class="form-control" id="fname" />
      </div>

      <div class="mb-4">
        <label for="lname" class="form-label">Last Name</label>
        <input type="text" class="form-control" id="lname" />
      </div>

      <div class="mb-4">
        <label for="email" class="form-label">Email address</label>
        <input type="email" class="form-control" id="email" />
      </div>

      <div class="mb-4">
        <label for="password" class="form-label">Password</label>
        <input type="password" class="form-control" id="password" />
      </div>

      <button type="submit" class="btn btn-primary">Submit</button>
    </form>
  </div>
</template>