<template>
  <div>
    {{ message }}
  </div>
</template>

<script>
import amqp from "amqplib/callback_api";

export default {
  name: "amqp",
  data() {
    return {
      message: "Hello AMQP!",
    };
  },

  created() {
    amqp.connect("amqp://localhost", function (error0, connection) {
      if (error0) {
        throw error0;
      }

      connection.createChannel(function (error1, channel) {
        if (error1) {
          throw error1;
        }

        var queue = "hello";

        channel.assertQueue(queue, {
          durable: false,
        });

        channel.consume(
          queue,
          (msg) => {
            console.log(" [x] Received %s", msg.content.toString());
          },
          {
            noAck: true,
          }
        );
      });
    });
  },
};
</script>