import { createAlova } from 'alova';
import adapterFetch from 'alova/fetch';
import VueHook from 'alova/vue';

const ASSET_URL = import.meta.env.VITE_ASSETS_URL;
const API_URL = import.meta.env.VITE_API_URL;

const sessionStorageAdapter = {
  set(key: string, value: unknown) {
    sessionStorage.setItem(key, JSON.stringify(value));
  },
  get(key: string) {
    const data = sessionStorage.getItem(key);
    return data ? JSON.parse(data) : data;
  },
  remove(key: string) {
    sessionStorage.removeItem(key);
  },
  clear() {
    sessionStorage.clear();
  }
};
export const dataInstance = createAlova({
  requestAdapter: adapterFetch(),
  statesHook: VueHook,
  responded: (response) => response.json(),
  baseURL: `${API_URL}`,
  timeout: 10000,
  cacheFor: {
    GET: 5 * 60 * 1000,
  },
  l1Cache: sessionStorageAdapter,
});

export const assetsInstance = createAlova({
  requestAdapter: adapterFetch(),
  statesHook: VueHook,
  responded: (response) => response.json(),
  baseURL: `${ASSET_URL}`,
  timeout: 10000,
});