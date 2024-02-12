import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import PhotosView from '../views/PhotosView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/',   redirect: '/login'},
        {path: '/login',  component: LoginView},
        {path: '/homepage',   component: HomeView},
		{path: '/users/:username/profile', component: ProfileView},
		{path: '/photos',   component: PhotosView}
	]
})

export default router
