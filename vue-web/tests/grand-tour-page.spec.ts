import { test, expect } from '@playwright/test'

const vehicles = [
  {
    id: '11111111-1111-1111-1111-111111111111',
    make: 'Honda',
    model: 'Civic Type R',
    year: 2019,
    body_type: 'hatchback',
    engine: '2.0L turbo',
    power: '320 hp',
    transmission: '6-speed manual',
    drivetrain: 'FWD',
    fuel: 'petrol',
    country_code: 'jp',
    country: 'Japan',
    city: 'Tokyo',
    driven_at: '2019-05-01T00:00:00Z',
    notes: 'Bought used with 42,000 km.',
    rating: 5,
    sort_order: 1,
    cutout_key: null,
    silhouette_key: null
  },
  {
    id: '22222222-2222-2222-2222-222222222222',
    make: 'Toyota',
    model: 'Corolla',
    year: 2012,
    body_type: null,
    country_code: null,
    country: 'Australia',
    city: null,
    driven_at: '2013-01-01T00:00:00Z',
    notes: null,
    rating: null,
    sort_order: null,
    silhouette_key: null
  }
]

test.describe('Grand Tour page', () => {
  test('renders carousel, flags and slide content', async ({ page }) => {
    const external: string[] = []
    page.on('request', (request) => {
      if (!request.url().startsWith('http://localhost')) external.push(request.url())
    })

    await page.route('**/vehicles', (route) =>
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(vehicles) })
    )

    await page.goto('/drive')

    await expect(page).toHaveTitle(/Grand Tour/)
    await expect(page.getByRole('heading', { name: 'The Grand Tour' })).toHaveCount(1)
    await expect(page.getByRole('link', { name: 'Grand Tour' })).toHaveCount(0)
    await expect(page.locator('article').first()).toContainText('Civic Type R')
    await expect(page.locator('article').first()).toContainText('Japan · Tokyo')
    await expect(page.locator('article').first()).toContainText('320 hp')
    await expect(page.locator('article').first()).toContainText('Bought used with 42,000 km.')

    const svgs = await page.locator('article').first().locator('svg').count()
    expect(svgs).toBeGreaterThan(0)

    expect(external.filter((url) => /flag|iconify|circle/.test(url))).toEqual([])
  })

  test('stacks neighbours behind the active vehicle and advances with keys', async ({ page }) => {
    await page.route('**/vehicles', (route) =>
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(vehicles) })
    )

    await page.goto('/drive')
    await page.locator('article').first().waitFor()

    const active = page.locator('[data-active="true"]')
    await expect(active).toContainText('Civic Type R')

    const slides = page.locator('.carousel-slide')
    await expect(slides).toHaveCount(2)

    const neighbour = slides.nth(1)
    await expect(neighbour).toHaveAttribute('aria-hidden', 'true')

    const transform = await neighbour.evaluate((element) => getComputedStyle(element).transform)
    expect(transform).not.toBe('none')

    await page.getByLabel('Next vehicle').click()
    await expect(active).toContainText('Corolla')

    await page.keyboard.press('ArrowLeft')
    await expect(active).toContainText('Civic Type R')
  })

  test('wraps from the first vehicle to the last', async ({ page }) => {
    await page.route('**/vehicles', (route) =>
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(vehicles) })
    )

    await page.goto('/drive')
    await page.locator('article').first().waitFor()

    const active = page.locator('[data-active="true"]')
    await expect(active).toContainText('Civic Type R')

    await page.keyboard.press('ArrowLeft')
    await expect(active).toContainText('Corolla')
  })

  test('fits the viewport without vertical scrolling', async ({ page }) => {
    await page.route('**/vehicles', (route) =>
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(vehicles) })
    )

    await page.goto('/drive')
    await page.locator('article').first().waitFor()

    const overflow = await page.evaluate(() => ({
      scroll: document.documentElement.scrollHeight,
      viewport: window.innerHeight
    }))

    expect(overflow.scroll).toBeLessThanOrEqual(overflow.viewport + 1)
  })

  test('shows the empty collection message', async ({ page }) => {
    await page.route('**/vehicles', (route) =>
      route.fulfill({ status: 200, contentType: 'application/json', body: '[]' })
    )

    await page.goto('/drive')
    await expect(page.getByText('No vehicles parked here yet.', { exact: false })).toBeVisible()
  })

  test('shows placeholder slides while loading', async ({ page }) => {
    await page.route('**/vehicles', async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 1500))
      await route.fulfill({ status: 200, contentType: 'application/json', body: '[]' })
    })

    await page.goto('/drive')
    await expect(page.getByText('Warming up the garage…')).toBeAttached()
  })

  test('shows an Unknown placeholder when a vehicle has no silhouette', async ({ page }) => {
    await page.route('**/vehicles', (route) =>
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([vehicles[1]])
      })
    )

    await page.goto('/drive')
    await expect(page.locator('article').getByText('Unknown')).toBeVisible()
    await expect(page.locator('article .silhouette')).toHaveCount(0)
  })

  test('advances exactly one vehicle per button press and wraps', async ({ page }) => {
    const deck = ['Alpha', 'Bravo', 'Charlie', 'Delta', 'Echo'].map((model, index) => ({
      id: String(index + 1),
      make: 'Make',
      model,
      year: 2000 + index,
      silhouette_key: null,
      sort_order: index
    }))

    await page.route('**/vehicles', (route) =>
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(deck) })
    )

    await page.goto('/drive')
    const active = page.locator('[data-active="true"] h2')
    await expect(active).toContainText('Alpha')

    for (const model of ['Bravo', 'Charlie', 'Delta', 'Echo', 'Alpha']) {
      await page.getByLabel('Next vehicle').click()
      await expect(active).toContainText(model)
    }

    await page.getByLabel('Previous vehicle').click()
    await expect(active).toContainText('Echo')

    await page.getByLabel('Previous vehicle').click()
    await expect(active).toContainText('Delta')
  })
})
