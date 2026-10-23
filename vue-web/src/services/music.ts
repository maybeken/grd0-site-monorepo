import { useRequest } from 'alova/client'
import { dataInstance, adminInstance } from './api'

import type { Ref } from 'vue'
import type { Music, Playlist } from '@/interfaces/Music'

function getMusicRaw() {
  return dataInstance.Get<Music[]>('/music')
}

function getMusic(): { loading: Ref<boolean, boolean>; data: Ref<Music[]> } {
  try {
    const { loading, data } = useRequest(getMusicRaw())

    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

function upsertMusicRaw(music: Music) {
  return adminInstance.Put<Music>('/music', music)
}

function deleteMusicRaw(v: string) {
  return adminInstance.Delete(`/music/${v}`)
}

function getPlaylistsRaw() {
  return dataInstance.Get<Playlist[]>('/music/playlists')
}

async function getPlaylistRaw(id: string): Promise<Playlist | null> {
  const value = await dataInstance.Get<Playlist>(`/music/playlist/${id}`)

  return value && typeof value === 'object' && 'id' in value ? value : null
}

function getPlaylists(): { loading: Ref<boolean, boolean>; data: Ref<Playlist[]> } {
  try {
    const { loading, data } = useRequest(getPlaylistsRaw())

    return { loading, data }
  } catch (error: unknown) {
    throw error
  }
}

function upsertPlaylistRaw(playlist: Playlist) {
  return adminInstance.Put<Playlist>('/music/playlist', playlist)
}

function deletePlaylistRaw(id: string) {
  return adminInstance.Delete(`/music/playlist/${id}`)
}

export {
  getMusicRaw,
  getMusic,
  upsertMusicRaw,
  deleteMusicRaw,
  getPlaylistsRaw,
  getPlaylists,
  getPlaylistRaw,
  upsertPlaylistRaw,
  deletePlaylistRaw
}
