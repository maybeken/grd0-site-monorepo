import { useRequest } from 'alova/client';
import { dataInstance } from './api';

import type { Ref } from 'vue';
import type { BlogPost } from '@/interfaces/Blog';

function listBlogPostRaw() {
  return dataInstance.Get<BlogPost[]>('/blog');
}

function listBlogPost(): {loading: Ref<boolean, boolean>, data: Ref<BlogPost[]>} {
  try {
    const { loading, data } = useRequest(listBlogPostRaw());

    return { loading, data };
  } catch(error: unknown) {
    throw error;
  }
}

function getBlogPostRaw(uri: string) {
  return dataInstance.Get<BlogPost>(`/blog/${uri}`);
}

function getBlogPost(uri: string): {loading: Ref<boolean, boolean>, data: Ref<BlogPost>} | void {
  if (!uri) return;

  try {
    const { loading, data } = useRequest(getBlogPostRaw(uri));

    return { loading, data };
  } catch(error: unknown) {
    throw error;
  }
}

export { listBlogPostRaw, listBlogPost, getBlogPostRaw, getBlogPost };