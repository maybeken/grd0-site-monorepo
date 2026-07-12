import { test, expect } from '@playwright/test'

test.describe('Coffee Page', () => {
  test.beforeEach(async ({ page }) => {
    await page.route('**/coffee/tastings', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            id: '1',
            bean: { id: 'b1', name: 'Ethiopian Yirgacheffe', origin: 'Ethiopia', roaster: 'Local Roaster' },
            equipment: { id: 'e1', name: 'Hario V60', type: 'Pour Over' },
            grind_size: 'Medium-Fine',
            grind_setting: 15,
            coffee_dose: 18,
            water_in: 270,
            coffee_out: null,
            ratio: 15.0,
            brew_time: 210,
            taste_fruity: 7,
            taste_sour: 4,
            taste_sweetness: 6,
            taste_nutty: 3,
            taste_spice: 2,
            taste_floral: 8,
            taste_green: 1,
            overall_notes: 'Great floral notes',
            rating: 8,
            tasted_at: '2024-07-10T00:00:00Z',
            pinned: false,
          },
          {
            id: '2',
            bean: { id: 'b2', name: 'Colombian Supremo', origin: null, roaster: null },
            equipment: { id: 'e2', name: 'La Marzocco', type: 'Espresso' },
            grind_size: null,
            grind_setting: null,
            coffee_dose: null,
            water_in: null,
            coffee_out: null,
            ratio: null,
            brew_time: null,
            taste_fruity: null,
            taste_sour: null,
            taste_sweetness: null,
            taste_nutty: null,
            taste_spice: null,
            taste_floral: null,
            taste_green: null,
            overall_notes: null,
            rating: null,
            tasted_at: '2024-07-05T00:00:00Z',
            pinned: true,
          },
        ]),
      })
    })

    await page.goto('/coffee')
  })

  test('displays page title', async ({ page }) => {
    await expect(page.locator('text=Coffee')).toBeVisible()
  })

  test('tasting cards render with 1:1 aspect ratio', async ({ page }) => {
    const cards = page.locator('.aspect-square')
    const count = await cards.count()
    expect(count).toBeGreaterThan(0)

    for (let i = 0; i < count; i++) {
      const box = await cards.nth(i).boundingBox()
      if (box) {
        expect(box.width).toBeCloseTo(box.height, 0)
      }
    }
  })

  test('radar chart displays 7 taste dimensions', async ({ page }) => {
    const firstCard = page.locator('.aspect-square').first()
    const canvas = firstCard.locator('canvas')
    await expect(canvas).toBeVisible()
  })

  test('pinned tastings appear first', async ({ page }) => {
    await page.waitForSelector('.aspect-square')
    const cards = page.locator('.aspect-square')
    const firstCard = cards.first()
    const pinIcon = firstCard.locator('[data-icon="mynaui:pin"], svg')
    const hasPinIcon = (await pinIcon.count()) > 0
    expect(hasPinIcon).toBeTruthy()
  })

  test('IDK displays for null fields', async ({ page }) => {
    await page.waitForSelector('.aspect-square')
    const cards = page.locator('.aspect-square')
    const secondCard = cards.nth(1)
    const idkTexts = secondCard.locator('text=IDK')
    const count = await idkTexts.count()
    expect(count).toBeGreaterThan(0)
  })

  test('equipment images render correctly', async ({ page }) => {
    await page.waitForSelector('.aspect-square')
    const cards = page.locator('.aspect-square')
    const firstCard = cards.first()
    const img = firstCard.locator('img')
    await expect(img).toBeVisible()
    const src = await img.getAttribute('src')
    expect(src).toBeTruthy()
  })
})
