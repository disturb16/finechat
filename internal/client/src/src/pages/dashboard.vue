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
                <button
                  type="button"
                  class="btn btn-primary"
                  data-bs-toggle="modal"
                  data-bs-target="#crModal"
                >
                  create chatroom
                </button>
                <ul>
                  <li v-for="cr in chatrooms" :key="cr.id">
                    <router-link :to="chatRoomUrl(cr.id)">{{
                      cr.name
                    }}</router-link>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </nav>
    </aside>

    <section>
      <router-view></router-view>
    </section>

    <!-- Modal -->
    <div
      class="modal fade"
      id="crModal"
      tabindex="-1"
      aria-labelledby="createChatRoomModalLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="createChatRoomModalLabel">
              New Chatroom
            </h5>
            <button
              type="button"
              class="btn-close"
              data-bs-dismiss="modal"
              aria-label="Close"
            ></button>
          </div>
          <div class="modal-body">
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
            </form>
          </div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-secondary"
              data-bs-dismiss="modal"
            >
              Close
            </button>
            <button
              type="button"
              class="btn btn-primary"
              @click.prevent="createChatRoom"
            >
              Save
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import { Modal } from "bootstrap";

export default {
  name: "Dashboard",
  data() {
    return {
      chatrooms: [],
      chatroomName: "",
      modal: null,
    };
  },

  async created() {
    this.getChatRooms();
  },

  methods: {
    chatRoomUrl(id) {
      return `/chatroom/${id}`;
    },

    async getChatRooms() {
      const { email } = this.$store.getters.tokenClaims;

      try {
        const url = `/api/users/${email}/chatrooms`;
        const response = await axios.get(url);
        this.chatrooms = response.data;
      } catch (error) {
        console.error(error);
      }
    },

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
        await axios.post("/api/chatrooms", data);
        this.getChatRooms();

        const modal = Modal.getInstance(document.getElementById("crModal"));
        modal.toggle();
      } catch (error) {
        console.error(error);
      }
    },
  },
};
</script>

<style scoped>
aside,
section {
  display: inline-block;
  vertical-align: top;
  width: 30%;
}

section {
  width: 65%;
}
</style>