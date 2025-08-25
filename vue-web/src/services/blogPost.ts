import { useRequest } from 'alova/client';
import { dataInstance, adminInstance } from './api';

import type { Ref } from 'vue';
import type { BlogPost, Attachment } from '@/interfaces/Blog';

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

function listBlogPostAdminRaw() {
  return adminInstance.Get<BlogPost[]>('/blog/all');
}

function listBlogPostAdmin(): {loading: Ref<boolean, boolean>, data: Ref<BlogPost[]>} {
  try {
    const { loading, data } = useRequest(listBlogPostAdminRaw());

    return { loading, data };
  } catch(error: unknown) {
    throw error;
  }
}

function upsertBlogPostRaw(url: string, blog_post: BlogPost) {
  return adminInstance.Put<BlogPost>(`/blog/${url}`, blog_post);
}

function requestBlogAttachmentRaw(path: string) {
  return adminInstance.Put<Attachment>(`/blog/attachment/${path}`);
}

function deleteBlogPostRaw(url: string) {
  return adminInstance.Delete(`/blog/${url}`);
}

export { listBlogPostRaw, listBlogPostAdmin, listBlogPost, getBlogPostRaw, getBlogPost, upsertBlogPostRaw, requestBlogAttachmentRaw, deleteBlogPostRaw };