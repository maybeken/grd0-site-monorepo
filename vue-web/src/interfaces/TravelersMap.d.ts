export interface MapLocation {
  title: string
  icon: string
  color: string
  pos: [number, number] // Record coordinate by North-East
  subtitle?: string
  display_at?: number
  hide_at?: number
  text_color?: string
}
