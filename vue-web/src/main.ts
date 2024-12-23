import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import router from './router';

import { OhVueIcon, addIcons } from 'oh-vue-icons';
import { MdMenu, MdOpeninnewOutlined } from 'oh-vue-icons/icons';
// import * as Icons from 'oh-vue-icons/icons';

// const icons = Object.values({ ...Icons });
// addIcons(...icons);
addIcons(MdMenu, MdOpeninnewOutlined);

const app = createApp(App);

app.use(createPinia());
app.use(router);

app.component('v-icon', OhVueIcon);

app.mount('#app');
