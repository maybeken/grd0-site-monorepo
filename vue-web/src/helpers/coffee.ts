import dripper from '@/assets/coffee/equipment/dripper.svg'
import portaFilter from '@/assets/coffee/equipment/porta-filter.svg'
import syphon from '@/assets/coffee/equipment/syphon-coffee.svg'
import frenchPress from '@/assets/coffee/equipment/french-press.svg'
import mokaPot from '@/assets/coffee/equipment/moka-pot.svg'
// import coldBrew from '@/assets/coffee/equipment/cold-brew.svg'
import beans from '@/assets/coffee/equipment/beans.svg'

export const GRIND_SIZES = [
  'Extra Fine',
  'Fine',
  'Medium-Fine',
  'Medium',
  'Medium-Coarse',
  'Coarse',
  'Extra Coarse',
] as const

export const EQUIPMENT_IMAGE_MAP: Record<string, string> = {
  'Pour Over': dripper,
  Espresso: portaFilter,
  'French Press': frenchPress,
  'Moka Pot': mokaPot,
  // 'Cold Brew': coldBrew,
  Syphon: syphon,
  Other: beans,
} as const

export const EQUIPMENT_TYPES = Object.keys(EQUIPMENT_IMAGE_MAP)

export const TASTE_DIMENSIONS = [
  'fruity',
  'sour',
  'sweetness',
  'nutty',
  'spice',
  'floral',
  'green',
] as const

export const DEFAULT_EQUIPMENT_IMAGE = beans
