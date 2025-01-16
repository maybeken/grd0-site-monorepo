import { useRequest } from 'alova';
import { dataInstance } from './api';

import type { Ref } from 'vue';
import type { BlogPost } from '@/interfaces/Blog';

function listBlogPost(): Ref<BlogPost[]> {
  try {
    const { data } = useRequest(dataInstance.Get<BlogPost[]>('/blog'));

    return data;
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