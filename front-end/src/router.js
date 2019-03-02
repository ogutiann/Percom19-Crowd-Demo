import Vue from "vue";
import Router from "vue-router";

Vue.use(Router);

export default new Router({
  routes: [
      {
          path: "/admin",
          name: "admin",
          // route level code-splitting
          // this generates a separate chunk (about.[hash].js) for this route
          // which is lazy-loaded when the route is visited.
          component: () =>
              import(/* webpackChunkName: "admin" */ "./views/admin.vue")
      },
    {
      path: "",
      name: "home",
      component: ()=> import("./views/Home.vue")
    },
  ]
});
