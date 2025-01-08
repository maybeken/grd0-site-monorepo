import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import router from './router';

import OpenLayersMap from "vue3-openlayers";

import { OhVueIcon, addIcons } from 'oh-vue-icons';
// Import on MD icons
import * as MdIcons from 'oh-vue-icons/icons/md';

const icons = Object.values({ ...MdIcons });
addIcons(...icons);

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(OpenLayersMap);

app.component('v-icon', OhVueIcon);

app.mount('#app');
