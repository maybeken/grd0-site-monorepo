import { test, expect } from '@playwright/test'

function createMockAssets(startIndex: number, count: number) {
  return Array.from({ length: count }, (_, i) => ({
    filename: `photo_${startIndex + i}.jpg`,
    exif: {
      datetime: `2024:01:${String(startIndex + i + 1).padStart(2, '0')} 12:00:00`,
      shutter: '1/250',
      fstop: '2.8',
      iso: 100,
      focal: 50.0,
      equipment: { camera: 'Test Camera', lens: 'Test Lens' },
    },
    collection: '/gallery/test-collection',
  }))
}

test.describe('Gallery Pagination', () => {
  test.beforeEach(async ({ page }) => {
    await page.route('**/v2/gallery/all*', (route) => {
      const url = new URL(route.request().url())
      const pageParam = parseInt(url.searchParams.get('page') || '1')
      const pageSize = parseInt(url.searchParams.get('page_size') || '40')

      const totalItems = 100
      const totalPages = Math.ceil(totalItems / pageSize)
      const startIndex = (pageParam - 1) * pageSize
      const remainingItems = Math.max(0, totalItems - startIndex)
      const itemsInPage = Math.min(pageSize, remainingItems)

      const response = {
        data: createMockAssets(startIndex, itemsInPage),
        total: totalItems,
        page: pageParam,
        page_size: pageSize,
        total_pages: totalPages,
      }

      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(response),
      })
    })

    await page.route('**/v2/gallery/details/*', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: [],
          total: 0,
          page: 1,
          page_size: 40,
          total_pages: 0,
        }),
      })
    })

    await page.route('**/gallery/collection', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          'test-collection': { title: 'Test Collection', cover: 'cover.jpg' },
        }),
      })
    })

    await page.goto('/gallery')
  })

  test('displays gallery title', async ({ page }) => {
    await expect(page.locator('text=Captures The Moment')).toBeVisible()
  })

  test('loads first page with 40 items by default', async ({ page }) => {
    const cards = page.locator('button[title="Download"]')
    await expect(cards).toHaveCount(40)
  })

  test('pagination response includes metadata', async ({ page }) => {
    const responsePromise = page.waitForResponse((response) =>
      response.url().includes('/v2/gallery/all')
    )

    await page.reload()
    const response = await responsePromise
    const body = await response.json()

    expect(body).toHaveProperty('data')
    expect(body).toHaveProperty('total')
    expect(body).toHaveProperty('page')
    expect(body).toHaveProperty('page_size')
    expect(body).toHaveProperty('total_pages')

    expect(body.total).toBe(100)
    expect(body.page).toBe(1)
    expect(body.page_size).toBe(40)
    expect(body.total_pages).toBe(3)
  })

  test('auto-loads more items when scrolling to bottom', async ({ page }) => {
    const cards = page.locator('button[title="Download"]')
    await expect(cards).toHaveCount(40)

    await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight))

    await page.waitForFunction(
      () => document.querySelectorAll('button[title="Download"]').length > 40,
      { timeout: 5000 }
    )

    await expect(cards).toHaveCount(80)
  })

  test('loads all items through multiple pages', async ({ page }) => {
    const cards = page.locator('button[title="Download"]')
    await expect(cards).toHaveCount(40)

    await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight))
    await page.waitForFunction(
      () => document.querySelectorAll('button[title="Download"]').length > 40,
      { timeout: 5000 }
    )

    await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight))
    await page.waitForFunction(
      () => document.querySelectorAll('button[title="Download"]').length > 80,
      { timeout: 5000 }
    )

    await expect(cards).toHaveCount(100)
  })

  test('does not load more when all items are loaded', async ({ page }) => {
    const cards = page.locator('button[title="Download"]')

    await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight))
    await page.waitForFunction(
      () => document.querySelectorAll('button[title="Download"]').length > 40,
      { timeout: 5000 }
    )

    await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight))
    await page.waitForFunction(
      () => document.querySelectorAll('button[title="Download"]').length > 80,
      { timeout: 5000 }
    )

    await expect(cards).toHaveCount(100)

    await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight))
    await page.waitForTimeout(1000)

    await expect(cards).toHaveCount(100)
  })

  test('resets pagination when collection changes', async ({ page }) => {
    const cards = page.locator('button[title="Download"]')

    await page.evaluate(() => window.scrollTo(0, document.body.scrollHeight))
    await page.waitForFunction(
      () => document.querySelectorAll('button[title="Download"]').length > 40,
      { timeout: 5000 }
    )

    const dropdownToggle = page.locator('.relative.h-12 .cursor-default').first()
    await dropdownToggle.click()

    const testCollectionOption = page.locator('text=Test Collection').first()
    await testCollectionOption.click()

    await page.waitForFunction(
      () => document.querySelectorAll('button[title="Download"]').length === 40,
      { timeout: 5000 }
    )

    await expect(cards).toHaveCount(40)
  })
})
