import { useRequest } from 'alova';
import { dataInstance } from './api';

import type { Ref } from 'vue';
import type { BlogPost } from '@/interfaces/Blog';

function listBlogPost(): Ref<BlogPost[]> {
  try {
    const { data } = useRequest(dataInstance.Get<BlogPost[]>('/blogPost.json'));

    return data;
  } catch(error: unknown) {
    throw error;
  }
}

export { listBlogPost };