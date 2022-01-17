<template>
  <div>
    <aside>
      <nav>
        <div class="accordion" id="accordionPanelsStayOpenExample">
          <div class="accordion-item">
            <h2 class="accordion-header" id="panelsStayOpen-headingOne">
              <button
                class="accordion-button"
                type="button"
                data-bs-toggle="collapse"
                data-bs-target="#panelsStayOpen-collapseOne"
                aria-expanded="true"
                aria-controls="panelsStayOpen-collapseOne"
              >
                Chatrooms
              </button>
            </h2>
            <div
              id="panelsStayOpen-collapseOne"
              class="accordion-collapse collapse show"
              aria-labelledby="panelsStayOpen-headingOne"
            >
              <div class="accordion-body">
                <div class="row">
                  <div class="col col-md6">
                    <router-link class="btn btn-primary" to="/chatrooms/create"
                      >New</router-link
                    >
                  </div>
                  <div class="col col-md6">
                    <a
                      href="#"
                      class="btn btn-primary"
                      @click.prevent="reloadChatRooms"
                    >
                      refresh
                    </a>
                  </div>
                </div>

                <ul class="list-item">
                  <li v-for="cr in chatRooms" :key="cr.id">
                    <router-link :to="chatRoomUrl(cr.id)">{{
                      cr.name
                    }}</router-link>
                  </li>
                </ul>
              </div>
            </div>
          </div>
          <!-- <div class="accordion-item">
            <h2 class="accordion-header" id="friends-panel">
              <button
                class="accordion-button"
                type="button"
                data-bs-toggle="collapse"
                data-bs-target="#friends-panel-collapseOne"
                aria-controls="friends-panel-collapseOne"
              >
                Friends
              </button>
            </h2>
            <div
              id="friends-panel-collapseOne"
              class="accordion-collapse collapse"
              aria-labelledby="friends-panel"
            >
              <div class="accordion-body">
                <router-link class="btn btn-primary" to="/friends/add"
                  >Add</router-link
                >

                <ul>
                  <li class="list-item" v-for="f in friends" :key="f.email">
                    {{ f.name }}
                  </li>
                </ul>
              </div>
            </div>
          </div> -->
        </div>
      </nav>
    </aside>

    <section>
      <router-view></router-view>
    </section>
  </div>
</template>

<script>
export default {
  name: "Dashboard",

  async created() {
    this.$store.dispatch("fetchChatRooms");
    this.$store.dispatch("fetchFriends");
  },

  methods: {
    chatRoomUrl(id) {
      return `/chatrooms/${id}`;
    },

    reloadChatRooms() {
      this.$store.dispatch("fetchChatRooms");
    },
  },

  computed: {
    chatRooms() {
      return this.$store.state.chatRooms;
    },

    friends() {
      return this.$store.state.friends;
    },
  },
};
</script>

<style scoped>
aside,
section {
  display: inline-block;
  vertical-align: top;
  width: 25%;
}

section {
  padding: 1em;
  margin-left: 1em;
  width: 70%;
}

.list-item {
  margin-top: 0.5em;
  padding: 1em;
}
</style>