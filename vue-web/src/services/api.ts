import { createAlova } from 'alova';
import adapterFetch from 'alova/fetch';
import VueHook from 'alova/vue';

const API_URL = import.meta.env.VITE_API_URL;
const IS_DEV = import.meta.env.DEV;

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
  beforeRequest: (method) => {
    method.config.headers.Authorization = `Bearer ${sessionStorage.getItem("jwt_token")}`;
  },
  responded: (response) => response.json(),
  baseURL: `${API_URL}`,
  timeout: 10000,
  cacheFor: {
    GET: 5 * 60 * 1000,
  },
  l1Cache: IS_DEV ? undefined : sessionStorageAdapter,
});