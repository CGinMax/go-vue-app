import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '@/components/HelloWorld'
import MyVuetify from '@/components/MyVuetify'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'HelloWorld',
      component: MyVuetify
    },
    {
      path: '/vue',
      name: 'MyVuetify',
      component: HelloWorld

    }
  ]
})
