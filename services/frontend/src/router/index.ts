import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import {client} from "../client/client";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Home",
    component: () => import("../views/Home.vue")
  },
  {
    path: "/login",
    name: "Login",
    component: () => import("../views/Login.vue")
  },
  {
    path: "/register",
    name: "Register",
    component: () => import("../views/Register.vue")
  }
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
});


router.beforeEach((to, from, next) => {
  if (to.name === 'Home' && !client.IsLoggedIn() ) {
    return next({ name: 'Login' })
  }
  if (to.name === 'Login' && client.IsLoggedIn() ) {
    return next({ name: 'Home' })
  }

  next()
})


export default router;
