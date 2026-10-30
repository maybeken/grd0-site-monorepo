export interface Vehicle {
  id?: string
  make: string
  model: string
  year: number
  body_type?: string | null
  engine?: string | null
  power?: string | null
  transmission?: string | null
  drivetrain?: string | null
  fuel?: string | null
  country_code?: string | null
  country?: string | null
  city?: string | null
  driven_at?: string | null
  notes?: string | null
  rating?: number | null
  sort_order?: number | null
  cutout_key?: string | null
  silhouette_key?: string | null
  created_at?: string
  updated_at?: string
}

export interface VehicleUpload {
  url: string
  key: string
}

export interface GeneratedCover {
  cutout: Blob
  silhouette_light: Blob
  silhouette_dark: Blob
  cutout_url: string
  silhouette_light_url: string
  silhouette_dark_url: string
}
