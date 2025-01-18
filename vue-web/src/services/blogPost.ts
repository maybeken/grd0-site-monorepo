import { useRequest } from 'alova/client';
import { dataInstance } from './api';

import type { Ref } from 'vue';
import type { BlogPost } from '@/interfaces/Blog';

function listBlogPost(): {loading: Ref<boolean, boolean>, data: Ref<BlogPost[]>} {
  try {
    const { loading, data } = useRequest(dataInstance.Get<BlogPost[]>('/blog'));

    return { loading, data };
  } catch(error: unknown) {
    throw error;
  }
}

function getBlogPost(uri: string): {loading: Ref<boolean, boolean>, data: Ref<BlogPost>} | void {
  if (!uri) return;

  try {
    const { loading, data } = useRequest(dataInstance.Get<BlogPost>(`/blog/${uri}`));

    return { loading, data };
  } catch(error: unknown) {
    throw error;
  }
}

export { listBlogPost, getBlogPost };