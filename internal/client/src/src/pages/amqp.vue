<template>
  <div>
    {{ message }}

    <button @click.prevent="socket.send('Hello there')">Send Message</button>
  </div>
</template>

<script>
export default {
  name: "amqp",
  props: ["chatRoomId"],
  data() {
    return {
      socket: null,
      message: "Hello AMQP!",
    };
  },

  created() {
    const loc = window.location;
    let uri = loc.protocol === "https:" ? "wss:" : "ws:";

    uri += `//${loc.host}/ws/${this.chatRoomId}`;
    this.socket = new WebSocket(uri);

    // Connection opened
    this.socket.addEventListener("open", (event) => {
      console.log(event);
      this.socket.send("Hello Server!");
    });

    // Listen for messages
    this.socket.addEventListener("message", (event) => {
      console.log("Message from server ", event.data);
    });
  },
};
</script>