import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import router from './router';

import {
  Map,
  Layers,
  Sources,
  Interactions,
  MapControls,
} from "vue3-openlayers";

import { OhVueIcon, addIcons } from 'oh-vue-icons';
import * as Icons from './icons';

const icons = Object.values({ ...Icons });
addIcons(...icons);

const app = createApp(App);

app.use(createPinia());
app.use(router);

app.use(Map);
app.use(Layers);
app.use(Sources);
app.use(Interactions);
app.use(MapControls);

app.component('v-icon', OhVueIcon);

app.mount('#app');
