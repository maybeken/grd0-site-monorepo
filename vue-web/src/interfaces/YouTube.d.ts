// Common types
export interface Message {
  event: EventType
  info: InitialDeliveryMessage | InfoDeliveryMessage
  target: string
}

export type EventType = 'onReady' | 'initialDelivery' | 'infoDelivery' | 'apiInfoDelivery'

export interface TimeInfo {
  currentTime: number
  currentTimeLastUpdated_: number
  mediaReferenceTime: number
}

export interface VideoQualityInfo {
  playbackQuality: string
  availableQualityLevels: string[]
}

export interface VideoDataMessage {
  videoData: VideoData
  videoStartBytes: number
  videoBytesTotal: number
  videoLoadedFraction: number
  videoUrl: string
  videoEmbedCode: string
  videoContentRect: VideoContentRect
}

export interface PlaylistInfo {
  playlist: string | null
  playlistIndex: number
  playlistId: string | null
}

export interface VolumeInfo {
  muted: boolean
  volume: number
}

// InitialDeliveryMessage now extends the common interfaces
export interface InitialDeliveryMessage
  extends TimeInfo, VideoDataMessage, PlaylistInfo, VolumeInfo {
  apiInterface: ApiInterface[]
  videoBytesLoaded: number
  debugText: string
  duration: number
  playerMode: Record<string, unknown>
  playerState: number
  availablePlaybackRates: number[]
  size: Size
  videoInfoVisible: boolean
  playbackRate: number
  sphericalProperties: Record<string, unknown>
}

export type InfoDeliveryMessage =
  | StateChangeMessage
  | PlaybackRateChangeMessage
  | PlaybackQualityChangeMessage
  | VolumeChangeMessage

// StateChangeMessage extends the common interfaces
export interface StateChangeMessage extends TimeInfo, VideoDataMessage, PlaylistInfo {
  playerState: PlayState
  duration: number
  playbackRate: number
  progressState: ProgressState
}

// PlaybackRateChangeMessage extends TimeInfo
export interface PlaybackRateChangeMessage extends TimeInfo {
  videoBytesLoaded: number
  videoLoadedFraction: number
  playbackRate: number
  progressState: ProgressState
}

// PlaybackQualityChangeMessage extends VideoQualityInfo
export interface PlaybackQualityChangeMessage extends VideoQualityInfo {
  preferredQuality: string
}

// VolumeChangeMessage remains unchanged
export interface VolumeChangeMessage extends VolumeInfo {}

export enum PlayState {
  UNSTARTED = -1,
  ENDED = 0,
  PLAYING = 1,
  PAUSED = 2,
  BUFFERING = 3,
  CUED = 4
}

type ApiActions =
  | 'cueVideoById'
  | 'loadVideoById'
  | 'cueVideoByUrl'
  | 'loadVideoByUrl'
  | 'playVideo'
  | 'pauseVideo'
  | 'stopVideo'
  | 'clearVideo'
  | 'getVideoBytesLoaded'
  | 'getVideoBytesTotal'
  | 'getVideoLoadedFraction'
  | 'getVideoStartBytes'
  | 'cuePlaylist'
  | 'loadPlaylist'
  | 'nextVideo'
  | 'previousVideo'
  | 'playVideoAt'
  | 'setShuffle'
  | 'setLoop'
  | 'getPlaylist'
  | 'getPlaylistIndex'
  | 'getPlaylistId'
  | 'loadModule'
  | 'unloadModule'
  | 'setOption'
  | 'getOption'
  | 'getOptions'
  | 'mute'
  | 'unMute'
  | 'isMuted'
  | 'setVolume'
  | 'getVolume'
  | 'seekTo'
  | 'getPlayerMode'
  | 'getPlayerState'
  | 'getAvailablePlaybackRates'
  | 'getPlaybackQuality'
  | 'setPlaybackQuality'
  | 'getAvailableQualityLevels'
  | 'getCurrentTime'
  | 'getDuration'
  | 'addEventListener'
  | 'removeEventListener'
  | 'getDebugText'
  | 'getVideoData'
  | 'addCueRange'
  | 'removeCueRange'
  | 'setSize'
  | 'getApiInterface'
  | 'destroy'
  | 'mutedAutoplay'
  | 'getVideoEmbedCode'
  | 'getVideoUrl'
  | 'getMediaReferenceTime'
  | 'getSize'
  | 'setFauxFullscreen'
  | 'logImaAdEvent'
  | 'preloadVideoById'
  | 'wakeUpControls'
  | 'showVideoInfo'
  | 'hideVideoInfo'
  | 'isVideoInfoVisible'
  | 'getPlaybackRate'
  | 'setPlaybackRate'
  | 'getSphericalProperties'
  | 'setSphericalProperties'

interface Size {
  width: number
  height: number
}

interface VideoData {
  video_id: string
  author: string
  title: string
  isPlayable: boolean
  errorCode: string | null
  backgroundable: boolean
  cpn: string
  isLive: boolean
  isWindowedLive: boolean
  isManifestless: boolean
  allowLiveDvr: boolean
  isListed: boolean
  isMultiChannelAudio: boolean
  hasProgressBarBoundaries: boolean
  isPremiere: boolean
  playerResponseCpn: string
  progressBarStartPositionUtcTimeMillis: number | null
  progressBarEndPositionUtcTimeMillis: number | null
  paidContentOverlayDurationMs: number
}

interface VideoContentRect {
  left: number
  top: number
  width: number
  height: number
}

interface ProgressState {
  airingStart: number
  airingEnd: number
  allowSeeking: boolean
  clipEnd: number | null
  clipStart: number
  current: number
  displayedStart: number
  duration: number
  ingestionTime: number | null
  isAtLiveHead: boolean
  loaded: number
  seekableStart: number
  seekableEnd: number
  offset: number
  viewerLivestreamJoinMediaTime: number
}
