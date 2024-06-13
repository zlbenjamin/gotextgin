import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import TextView from '../views/TextView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: {
        title: 'Home | Hi, text!'
      }
    },
    {
      path: '/text/:id',
      name: 'text',
      component: TextView,
      meta: {
        title: 'Detail of Text'
      }
    },
  ]
})

router.beforeEach((to, from, next) => {
  // modify title
  if (to.meta.title) {
    document.title = to.meta.title
  }

  next()
})

export default router
