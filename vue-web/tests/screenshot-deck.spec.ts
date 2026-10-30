import { test } from '@playwright/test'

const car = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 260"><path fill="#fff" d="M28 188l6-30q6-20 32-24l88-6 58-52q12-12 34-12h162q22 0 36 14l50 50 84 8q24 4 30 24l4 18q2 12-12 12H44q-14 0-16-2z"/><circle fill="#fff" cx="156" cy="190" r="38"/><circle fill="#fff" cx="484" cy="190" r="38"/></svg>`

const base = {
  engine: '2.0L',
  power: '200 hp',
  transmission: 'manual',
  drivetrain: 'FWD',
  fuel: 'petrol',
  country_code: 'JP',
  country: 'Japan',
  city: 'Tokyo',
  driven_at: '2019-05-01T00:00:00Z',
  rating: 4
}

const vehicles = [
  { ...base, id: '1', make: 'Honda', model: 'Civic Type R', year: 2019, sort_order: 1, silhouette_key: 'covers/1.svg' },
  { ...base, id: '2', make: 'Toyota', model: 'Corolla', year: 2012, sort_order: 2, silhouette_key: 'covers/2.svg' },
  { ...base, id: '3', make: 'Mazda', model: 'MX-5', year: 2016, sort_order: 3, silhouette_key: 'covers/3.svg' },
  { ...base, id: '4', make: 'Tesla', model: 'Model 3', year: 2021, sort_order: 4, silhouette_key: 'covers/4.svg' },
  { ...base, id: '5', make: 'Land Rover', model: 'Defender', year: 1997, sort_order: 5, silhouette_key: 'covers/5.svg' }
]

test('screenshot deck', async ({ page }) => {
  const dir = process.env.SHOT_DIR || '.'
  await page.route('http://localhost:8080/vehicles', (route) =>
    route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(vehicles) })
  )
  await page.route('https://assets.grd0.net/covers/**', (route) =>
    route.fulfill({ status: 200, contentType: 'image/svg+xml', body: car })
  )

  await page.setViewportSize({ width: 1280, height: 800 })
  await page.goto('/drive')
  await page.locator('article').first().waitFor()
  await page.waitForTimeout(900)
  await page.screenshot({ path: `${dir}/deck-1.png` })

  await page.getByLabel('Next vehicle').click()
  await page.waitForTimeout(900)
  await page.screenshot({ path: `${dir}/deck-2.png` })
})
