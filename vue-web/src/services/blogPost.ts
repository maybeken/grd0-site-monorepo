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

function getBlogPost(uri: string): Ref<BlogPost> | void {
  if (!uri) return;

  try {
    const { data } = useRequest(dataInstance.Get<BlogPost>(`/blog/${uri}`));

    return data;
  } catch(error: unknown) {
    throw error;
  }
}

export { listBlogPost, getBlogPost };