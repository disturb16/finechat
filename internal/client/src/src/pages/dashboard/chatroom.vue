<template>
  <div>
    <section>
      <div id="wrap">
        <div id="main" class="container clear-top">
          <div class="message" v-for="(m, i) in messages" :key="i">
            <div class="message-owner">{{ m.user }}</div>
            <div class="message-text">{{ m.message }}</div>
          </div>
        </div>
      </div>
      <footer class="footer">
        <div class="input-group">
          <textarea
            v-model="messageContent"
            class="form-control"
            aria-label="With textarea"
            @keypress.enter.prevent="sendMessage"
          ></textarea>
          <button
            class="btn btn-primary"
            type="button"
            @click.prevent="sendMessage"
          >
            Send
          </button>
        </div>
      </footer>
    </section>
    <aside>
      <h2>
        Participants
        <button
          type="button"
          class="btn btn-primary"
          @click.prevent="() => modal.show()"
        >
          Add
        </button>
      </h2>

      <ul>
        <li v-for="u in chatRoomUsers" :key="u.email">
          {{ u.name }}
          <a href="#" @click.prevent="removeParticipant(u.email)"> remove</a>
        </li>
      </ul>
    </aside>

    <div class="modal fade" id="addFriendModal" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Add new participant</h5>
            <button
              type="button"
              class="btn-close"
              data-bs-dismiss="modal"
              aria-label="Close"
            ></button>
          </div>
          <div class="modal-body">
            <div class="mb-4">
              <label for="newParticipantEmail" class="form-label">Email</label>
              <input
                type="text"
                class="form-control"
                id="newParticipantEmail"
                v-model="newParticipantEmail"
              />
            </div>
          </div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-primary"
              @click.prevent="saveNewParticipant"
            >
              Add
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
  name: "Chatroom",
  props: ["chatRoomId"],
  data() {
    return {
      socket: null,
      modal: null,
      messages: [],
      messageContent: "",
      chatRoomUsers: [{ email: "ss", name: "ss" }],
      newParticipantEmail: "",
    };
  },

  created() {
    this.setupWebsocket();
    this.getMessages(this.chatRoomId);
    this.fetcheChatRoomUsers(this.chatRoomId);
  },

  mounted() {
    if (this.modal == null) {
      this.modal = new Modal(document.getElementById("addFriendModal"), {});
    }
  },

  methods: {
    async sendMessage() {
      if (this.messageContent == "") {
        return;
      }

      try {
        const { name, email } = this.$store.getters.tokenClaims;

        const url = `/api/chatrooms/${this.chatRoomId}/messages`;
        const data = {
          message: this.messageContent,
          email,
          created_date: new Date(),
        };

        await axios.post(url, data);

        this.messages.push({
          user: name,
          message: data.message,
          createdDate: data.created_date,
        });

        this.messages.sort((a, b) => {
          const aDate = new Date(a.created_date);
          const bDate = new Date(b.created_date);
          return aDate - bDate;
        });

        this.messageContent = "";
      } catch (error) {
        console.error(error);
      }
    },

    async getMessages(chatRoomId) {
      try {
        const url = `/api/chatrooms/${chatRoomId}/messages`;
        const response = await axios.get(url);

        this.messages = response.data;
      } catch (error) {
        console.error(error);
      }
    },

    async fetcheChatRoomUsers(chatRoomId) {
      try {
        const url = `/api/chatrooms/${chatRoomId}/users`;
        const response = await axios.get(url);

        this.chatRoomUsers = response.data;
      } catch (error) {
        console.error(error);
      }
    },

    async saveNewParticipant() {
      try {
        const url = `/api/chatrooms/${this.chatRoomId}/users`;
        const data = {
          email: this.newParticipantEmail,
        };

        await axios.post(url, data);

        this.fetcheChatRoomUsers(this.chatRoomId);

        this.newParticipantEmail = "";
        this.modal.hide();
      } catch (error) {
        console.error(error);
      }
    },

    async removeParticipant(email) {
      try {
        const url = `/api/chatrooms/${this.chatRoomId}/users/${email}`;

        await axios.delete(url);

        this.fetcheChatRoomUsers(this.chatRoomId);
      } catch (error) {
        console.error(error);
      }
    },

    setupWebsocket() {
      const loc = window.location;
      let uri = loc.protocol === "https:" ? "wss:" : "ws:";
      const { email } = this.$store.getters.tokenClaims;
      uri += `//${loc.host}/ws/${this.chatRoomId}/email/${email}`;

      if (this.socket != null) {
        this.socket.close();
        this.socket = null;
      }

      this.socket = new WebSocket(uri);

      // Connection opened
      this.socket.addEventListener("open", () => {
        console.log("Connected to chatroom stream");
      });

      // Listen for messages
      this.socket.addEventListener("message", (event) => {
        const data = JSON.parse(event.data);

        switch (data.type) {
          case "reload":
            this.getMessages(this.chatRoomId);
            break;

          case "show_stock":
            // this.showStock(data.stock);
            break;

          default:
            break;
        }
      });
    },
  },
  watch: {
    chatRoomId(newChatRoomId) {
      this.getMessages(newChatRoomId);
      this.fetcheChatRoomUsers(newChatRoomId);
      this.setupWebsocket();
    },
  },
};
</script>

<style scoped>
html,
body {
  height: 100%;
}

section,
aside {
  display: inline-block;
  vertical-align: top;
  width: 65%;
}

aside {
  width: 30%;
  margin-left: 1em;
}

#wrap {
  min-height: 90vh;
}

#main {
  overflow: auto;
  padding-bottom: 150px; /* this needs to be bigger than footer height*/
  height: 60vh;
  overflow-y: scroll;
  scrollbar-color: rebeccapurple green;
  scrollbar-width: thin;
}

.footer {
  position: relative;
  margin-top: -150px; /* negative value of footer height */
  height: 150px;
  clear: both;
  padding-top: 20px;
}

.message {
  margin-bottom: 1em;
}

.message-owner {
  font-weight: bold;
}

.message-text,
.message-owner {
  padding: 0.5em;
  border-radius: 0.5em;
  background-color: #f5f5f5;
}
</style>