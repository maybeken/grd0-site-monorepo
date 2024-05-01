import { createAlova } from 'alova';
import GlobalFetch from 'alova/GlobalFetch';
import VueHook from 'alova/vue';

export const alovaInstance = createAlova({
  requestAdapter: GlobalFetch(),
  statesHook: VueHook,
  responsed: (response) => response.json(),
  baseURL: '/data',
  timeout: 10000,
});