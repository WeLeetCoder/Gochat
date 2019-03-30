import Vue from 'vue'
import Router from 'vue-router'
import login from '@/components/login'
import chatRoom from '@/components/chatRoom'
Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'login',
      component: login
    },
    {
      path: '/chatRoom',
      name: 'chatRoom',
      component: chatRoom
    }
  ]
})
