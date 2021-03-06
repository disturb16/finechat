import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

Vue.config.productionTip = false

router.beforeEach((to, from, next) => {
  if (to.meta.protected && !store.state.auth.token ) {
    next('/login')
  }

  next()
})

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app')
