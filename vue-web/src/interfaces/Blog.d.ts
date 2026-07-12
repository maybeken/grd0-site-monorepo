export interface BlogPost {
  uri: string
  author: {
    email: string
    display_name: string
  }
  title: string
  subtitle: string
  content: string
  created_at: string
  updated_at: string
  published_at: string
}

export interface Attachment {
  url: string
  key: string
}
