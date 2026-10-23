export interface Music {
  id?: string
  title: string
  v: string
  artist?: string
  description?: string
  sorting?: number
}

export interface Playlist {
  id?: string
  title: string
  description?: string
  sorting?: number
  songs: Music[]
}
