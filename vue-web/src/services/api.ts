import { createAlova } from 'alova';
import GlobalFetch from 'alova/GlobalFetch';
import VueHook from 'alova/vue';

const ASSET_URL = import.meta.env.VITE_ASSETS_URL;
const API_URL = import.meta.env.VITE_API_URL;

export const dataInstance = createAlova({
  requestAdapter: GlobalFetch(),
  statesHook: VueHook,
  responsed: (response) => response.json(),
  baseURL: `${API_URL}`,
  timeout: 10000,
});

export const assetsInstance = createAlova({
  requestAdapter: GlobalFetch(),
  statesHook: VueHook,
  responsed: (response) => response.json(),
  baseURL: `${ASSET_URL}`,
  timeout: 10000,
});