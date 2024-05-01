import { useRequest } from 'alova';
import { alovaInstance } from './api';

import type { Ref } from 'vue';
import type { BlogPost } from '@/interfaces/Blog';

function listBlogPost(): Ref<BlogPost[]> {
  try {
    const { data } = useRequest(alovaInstance.Get<BlogPost[]>('/blogPost.json'));

    return data;
  } catch(error: unknown) {
    throw error;
  }
}

export { listBlogPost };