import Vue from "vue";
import VueRouter from "vue-router";

Vue.use(VueRouter);

const registerPage = () => import("./pages/registration.vue");
const dashboardPage = () => import("./pages/dashboard.vue");
const loginPage = () => import("./pages/login.vue");

const routes = [
  {
    path: "/",
    component: dashboardPage,
    meta: {
      protected: true,
    },
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
