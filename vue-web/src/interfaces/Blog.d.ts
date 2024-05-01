export interface BlogPost {
  uri: string,
  author: {
    email: string,
    displayName: string,
  },
  title: string,
  subtitle: string,
  content: string,
  created_at: string,
  updated_at: string,
};