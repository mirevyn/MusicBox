import { createRouter, createWebHistory } from "vue-router";
import ClientLayout from "@/layout/ClientLayout.vue";
import Home from "@/views/Home.vue";
import Favorites from "@/views/Favorites.vue";
import Auth from "@/views/Auth.vue";
import AdminLayout from "@/layout/AdminLayout.vue";
import { useUserStore } from "@/stores/user";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // 客户端路由组
    {
      path: "/",
      component: ClientLayout,
      children: [
        {
          path: "", // 默认子路由
          name: "Home",
          component: Home, // 首页内容
        },
        {
          path: "playlists/hot",
          name: "HotPlaylists",
          component: () => import("@/views/playlist/HotPlaylists.vue"),
        },
        {
          path: "favorites",
          name: "Favorites",
          component: Favorites,
          meta: { requiresAuth: true },
        },
        {
          path: "daily",
          name: "DailyRecommend",
          component: () => import("@/views/DailyRecommend.vue"),
          meta: { requiresAuth: true },
        },
        // 歌单详情页路由
        // :id 是动态参数
        {
          path: "playlist/:id",
          name: "PlaylistDetail",
          // 使用路由懒加载，当访问该路由时才加载组件，提升首屏速度
          component: () => import("@/views/PlaylistDetail.vue"),
        },
        {
          path: "library",
          name: "Library",
          component: () => import("@/views/Library.vue"),
          meta: { requiresAuth: true },
        },
        {
          path: "search",
          name: "Search",
          component: () => import("@/views/Search.vue"),
        },
        {
          path: "profile",
          name: "Profile",
          component: () => import("@/views/Profile.vue"),
          meta: { requiresAuth: true },
        },
      ],
    },

    // 登录页
    {
      path: "/auth",
      name: "Auth",
      component: Auth,
      meta: { guestOnly: true },
    },

    // 管理端路由
    {
      path: "/admin",
      component: AdminLayout,
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: "",
          name: "AdminDashboard",
          component: () => import("@/components/admin/Dashboard.vue"),
        },
        {
          path: "songs",
          name: "AdminSongs",
          component: () => import("@/components/admin/Songs.vue"),
        },
        {
          path: "users",
          name: "AdminUsers",
          component: () => import("@/components/admin/Users.vue"),
        },
        {
          path: "playlists",
          name: "AdminPlaylists",
          component: () => import("@/components/admin/Playlists.vue"),
        },
        {
          path: "settings",
          name: "AdminSettings",
          component: () => import("@/components/admin/Settings.vue"),
        },
      ],
    },
  ],
});

// 路由守卫
router.beforeEach(async (to, _from, next) => {
  const token = localStorage.getItem("token");

  // 已登录用户不能访问登录页
  if (to.meta.guestOnly && token) {
    next("/");
    return;
  }

  // 需要认证的页面
  if (to.meta.requiresAuth && !token) {
    next("/auth");
    return;
  }

  // 需要管理员权限的页面
  if (to.meta.requiresAdmin) {
    const userStore = useUserStore();

    // 如果有 token 但没有用户信息，尝试获取
    if (token && !userStore.user) {
      await userStore.fetchProfile();
    }

    if (!userStore.isAdmin) {
      // 不是管理员，重定向到首页
      next("/");
      return;
    }
  }

  next();
});

export default router;
