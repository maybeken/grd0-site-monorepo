import dripper from '@/assets/coffee/equipment/dripper.svg'
import portaFilter from '@/assets/coffee/equipment/porta-filter.svg'
import syphon from '@/assets/coffee/equipment/syphon-coffee.svg'
import frenchPress from '@/assets/coffee/equipment/french-press.svg'
import mokaPot from '@/assets/coffee/equipment/moka-pot.svg'
// import coldBrew from '@/assets/coffee/equipment/cold-brew.svg'
import beans from '@/assets/coffee/equipment/beans.svg'

export const GRIND_SIZES = [
  { value: 'Extra Fine', label: 'coffee.grindSizes.extraFine' },
  { value: 'Fine', label: 'coffee.grindSizes.fine' },
  { value: 'Medium-Fine', label: 'coffee.grindSizes.mediumFine' },
  { value: 'Medium', label: 'coffee.grindSizes.medium' },
  { value: 'Medium-Coarse', label: 'coffee.grindSizes.mediumCoarse' },
  { value: 'Coarse', label: 'coffee.grindSizes.coarse' },
  { value: 'Extra Coarse', label: 'coffee.grindSizes.extraCoarse' }
] as const

export const EQUIPMENT_IMAGE_MAP: Record<string, string> = {
  Grinder: beans,
  'Pour Over': dripper,
  Espresso: portaFilter,
  'French Press': frenchPress,
  'Moka Pot': mokaPot,
  // 'Cold Brew': coldBrew,
  Syphon: syphon,
  Other: beans,
} as const

export const EQUIPMENT_TYPES = [
  { value: 'Grinder', label: 'coffee.equipmentTypes.grinder' },
  { value: 'Pour Over', label: 'coffee.equipmentTypes.pourOver' },
  { value: 'Espresso', label: 'coffee.equipmentTypes.espresso' },
  { value: 'French Press', label: 'coffee.equipmentTypes.frenchPress' },
  { value: 'Moka Pot', label: 'coffee.equipmentTypes.mokaPot' },
  { value: 'Syphon', label: 'coffee.equipmentTypes.syphon' },
  { value: 'Other', label: 'coffee.equipmentTypes.other' }
] as const

export const TASTE_DIMENSIONS = [
  { value: 'fruity', label: 'coffee.tasteDimensions.fruity' },
  { value: 'sour', label: 'coffee.tasteDimensions.sour' },
  { value: 'fermented', label: 'coffee.tasteDimensions.fermented' },
  { value: 'nutty', label: 'coffee.tasteDimensions.nutty' },
  { value: 'spice', label: 'coffee.tasteDimensions.spice' },
  { value: 'sweetness', label: 'coffee.tasteDimensions.sweetness' },
  { value: 'floral', label: 'coffee.tasteDimensions.floral' },
  { value: 'green', label: 'coffee.tasteDimensions.green' },
  { value: 'tobacco', label: 'coffee.tasteDimensions.tobacco' },
  { value: 'bitter', label: 'coffee.tasteDimensions.bitter' },
] as const

export const DEFAULT_EQUIPMENT_IMAGE = beans
