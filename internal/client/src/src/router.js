import Vue from "vue";
import VueRouter from "vue-router";

Vue.use(VueRouter);

const registerPage = () => import("./pages/registration.vue");
const dashboardPage = () => import("./pages/dashboard.vue");
const loginPage = () => import("./pages/login.vue");
const createChatroomPage = () => import("./pages/dashboard/createChatroom.vue");
const chatroomPage = () => import("./pages/dashboard/chatroom.vue");
const addFriendPage = () => import("./pages/dashboard/addFriend.vue");

const routes = [
  {
    path: "/",
    component: dashboardPage,
    meta: {
      protected: true,
    },
    children: [
      {
        path: "/chatrooms/create",
        component: createChatroomPage,
        meta: {
          protected: true,
        },
      },
      {
        path: "/chatrooms/:chatRoomId",
        component: chatroomPage,
        props: true,
        meta: {
          protected: true,
        },
      },
      {
        path: "/friends/add",
        component: addFriendPage,
        meta: {
          protected: true,
        },
      },
    ],
  },
  {
    path: "/register",
    component: registerPage,
  },
  {
    path: "/login",
    component: loginPage,
  },
];

const router = new VueRouter({
  routes,
});

export default router;
