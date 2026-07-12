export interface CoffeeBean {
  id?: string
  name: string
  origin?: string | null
  roaster?: string | null
  roast_date?: string | null
  variety?: string | null
  process?: string | null
  altitude?: string | null
  description?: string | null
  created_at?: string
  updated_at?: string
}

export interface BrewEquipment {
  id?: string
  name: string
  type?: string | null
  description?: string | null
  created_at?: string
  updated_at?: string
}

export interface TastingNote {
  id?: string
  bean: CoffeeBean
  equipment: BrewEquipment
  grinder?: BrewEquipment
  grind_size?: string | null
  grind_setting?: number | null
  coffee_dose?: number | null
  water_in?: number | null
  coffee_out?: number | null
  ratio?: number | null
  brew_time?: number | null
  water_temperature?: number | null
  taste_fruity?: number | null
  taste_sour?: number | null
  taste_fermented?: number | null
  taste_nutty?: number | null
  taste_spice?: number | null
  taste_sweetness?: number | null
  taste_floral?: number | null
  taste_green?: number | null
  taste_tobacco?: number | null
  taste_bitter?: number | null
  overall_notes?: string | null
  rating?: number | null
  tasted_at: string
  pinned: boolean
  created_at?: string
  updated_at?: string
}
