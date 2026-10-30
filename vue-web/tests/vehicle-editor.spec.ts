import { test, expect } from '@playwright/test'

const vehicle = {
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
  country_code: 'JP',
  country: 'Japan',
  city: 'Tokyo',
  driven_at: '2019-05-01T00:00:00Z',
  notes: 'Bought used.',
  rating: 5,
  sort_order: 1,
  cutout_key: null,
  silhouette_key: null
}

test.describe('Vehicle editor', () => {
  test('lists vehicles with identity and cover state', async ({ page }) => {
    await page.route('http://localhost:8080/vehicles', (route) =>
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify([vehicle]) })
    )

    await page.goto('/editor/vehicles')
    await expect(page.getByText('Honda Civic Type R')).toBeVisible()
    await expect(page.getByText('No cover')).toBeVisible()
  })

  test('offers a retry when the listing fails', async ({ page }) => {
    let attempts = 0
    await page.route('http://localhost:8080/vehicles', (route) => {
      attempts += 1
      if (attempts === 1) {
        return route.fulfill({ status: 500, contentType: 'application/json', body: '{}' })
      }
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([vehicle])
      })
    })

    await page.goto('/editor/vehicles')
    await expect(page.getByText('Could not load the editor.')).toBeVisible()
    await page.getByRole('button', { name: 'Retry' }).click()
    await expect(page.getByText('Honda Civic Type R')).toBeVisible()
  })

  test('blocks saving without a model and identifies the field', async ({ page }) => {
    let puts = 0
    await page.route('http://localhost:8080/vehicles', (route) => {
      if (route.request().method() === 'PUT') {
        puts += 1
        return route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ ...vehicle, id: 'new-id' })
        })
      }
      return route.fulfill({ status: 200, contentType: 'application/json', body: '[]' })
    })

    await page.goto('/editor/vehicles/new')
    await page.locator('#vehicle-make').fill('Honda')
    await page.locator('#vehicle-year').fill('2019')
    await page.getByRole('button', { name: 'Save' }).click()

    await expect(page.getByText('Model is required.')).toBeVisible()
    expect(puts).toBe(0)
  })

  test('saves with only the required fields', async ({ page }) => {
    let payload: Record<string, unknown> | null = null
    await page.route('http://localhost:8080/vehicles', (route) => {
      if (route.request().method() === 'PUT') {
        payload = route.request().postDataJSON() as Record<string, unknown>
        return route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ ...vehicle, ...(payload ?? {}) })
        })
      }
      return route.fulfill({ status: 200, contentType: 'application/json', body: '[]' })
    })
    await page.route('http://localhost:8080/vehicles/new-id', (route) =>
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(vehicle) })
    )

    await page.goto('/editor/vehicles/new')
    await page.locator('#vehicle-make').fill('Honda')
    await page.locator('#vehicle-model').fill('Civic')
    await page.locator('#vehicle-year').fill('2019')
    await page.getByRole('button', { name: 'Save' }).click()

    await expect.poll(() => payload).not.toBeNull()
    expect(payload).toMatchObject({ make: 'Honda', model: 'Civic', year: 2019 })
  })

  test('shows a recoverable error for an unknown vehicle', async ({ page }) => {
    await page.route('http://localhost:8080/vehicles/**', (route) =>
      route.fulfill({ status: 404, contentType: 'application/json', body: '{}' })
    )

    await page.goto('/editor/vehicles/does-not-exist')
    await expect(page.getByText('This vehicle could not be loaded.')).toBeVisible()
    await expect(page.getByRole('link', { name: 'Back to Vehicles' }).first()).toBeVisible()
  })

  test('loads an existing vehicle directly', async ({ page }) => {
    await page.route('http://localhost:8080/vehicles/11111111-1111-1111-1111-111111111111', (route) =>
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(vehicle) })
    )

    await page.goto('/editor/vehicles/11111111-1111-1111-1111-111111111111')
    await expect(page.locator('#vehicle-model')).toHaveValue('Civic Type R')
    await expect(page.getByText('No cover assigned.')).toBeVisible()
    await expect(page.getByText('Unknown')).toBeVisible()
  })

  test('deletes with confirmation', async ({ page }) => {
    let deleted = false

    await page.route('http://localhost:8080/vehicles/**', (route) => {
      if (route.request().method() === 'DELETE') {
        deleted = true
        return route.fulfill({ status: 202, contentType: 'application/json', body: 'null' })
      }
      return route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify([]) })
    })

    await page.route('http://localhost:8080/vehicles', (route) =>
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(deleted ? [] : [vehicle])
      })
    )

    page.on('dialog', (dialog) => dialog.accept())
    await page.goto('/editor/vehicles')
    await expect(page.getByText('Honda Civic Type R')).toBeVisible()
    await page.getByRole('button', { name: 'Delete' }).click()

    await expect(page.getByText('Honda Civic Type R')).toHaveCount(0)
    expect(deleted).toBe(true)
  })

  test('does not delete when the confirmation is dismissed', async ({ page }) => {
    let deletes = 0

    await page.route('http://localhost:8080/vehicles/**', (route) => {
      if (route.request().method() === 'DELETE') deletes += 1
      return route.fulfill({ status: 202, contentType: 'application/json', body: 'null' })
    })

    await page.route('http://localhost:8080/vehicles', (route) =>
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([vehicle])
      })
    )

    page.on('dialog', (dialog) => dialog.dismiss())
    await page.goto('/editor/vehicles')
    await expect(page.getByText('Honda Civic Type R')).toBeVisible()
    await page.getByRole('button', { name: 'Delete' }).click()
    await page.waitForTimeout(300)

    await expect(page.getByText('Honda Civic Type R')).toBeVisible()
    expect(deletes).toBe(0)
  })

  test('previews a stored silhouette and labels it as generated', async ({ page }) => {
    await page.route('http://localhost:8080/vehicles/11111111-1111-1111-1111-111111111111', (route) =>
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ ...vehicle, silhouette_key: 'vehicles/2026/10/cover.png' })
      })
    )

    await page.goto('/editor/vehicles/11111111-1111-1111-1111-111111111111')
    await expect(page.getByText('Generated cover is assigned.')).toBeVisible()
    await expect(page.locator('.silhouette')).toHaveCount(1)
  })
})

test.describe('Country field', () => {
  test('filters by code, shows the flag and stores label and code together', async ({ page }) => {
    let payload: Record<string, unknown> | null = null

    await page.route('http://localhost:8080/vehicles', (route) => {
      if (route.request().method() === 'PUT') {
        payload = route.request().postDataJSON() as Record<string, unknown>
        return route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ ...vehicle, ...(payload ?? {}) })
        })
      }
      return route.fulfill({ status: 200, contentType: 'application/json', body: '[]' })
    })
    await page.route('http://localhost:8080/vehicles/new-id', (route) =>
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(vehicle) })
    )

    await page.goto('/editor/vehicles/new')
    await page.locator('#vehicle-make').fill('Honda')
    await page.locator('#vehicle-model').fill('Civic')
    await page.locator('#vehicle-year').fill('2019')

    const field = page.locator('#vehicle-country')
    await field.fill('jp')
    await expect(page.getByRole('option', { name: /Japan/ })).toBeVisible()
    await page.getByRole('option', { name: /Japan/ }).click()
    await expect(field).toHaveValue('Japan')

    const wrapper = field.locator('..')
    expect(await wrapper.locator('svg').count()).toBeGreaterThan(0)

    await page.getByRole('button', { name: 'Save' }).click()
    await expect.poll(() => payload).not.toBeNull()
    expect(payload).toMatchObject({ country: 'Japan', country_code: 'JP' })
  })

  test('clears both the label and the code', async ({ page }) => {
    let payload: Record<string, unknown> | null = null

    await page.route('http://localhost:8080/vehicles', (route) => {
      if (route.request().method() !== 'PUT') {
        return route.fulfill({ status: 200, contentType: 'application/json', body: '[]' })
      }
      payload = route.request().postDataJSON() as Record<string, unknown>
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ ...vehicle, ...(payload ?? {}) })
      })
    })
    await page.route('http://localhost:8080/vehicles/11111111-1111-1111-1111-111111111111', (route) =>
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(vehicle) })
    )

    await page.goto('/editor/vehicles/11111111-1111-1111-1111-111111111111')
    const field = page.locator('#vehicle-country')
    await expect(field).toHaveValue('Japan')

    await page.getByRole('button', { name: 'Clear country' }).click()
    await expect(field).toHaveValue('')

    await page.getByRole('button', { name: 'Save' }).click()
    await expect.poll(() => payload).not.toBeNull()
    expect(payload).toMatchObject({ country: null, country_code: null })
  })

  test('reverts unmatched text on blur', async ({ page }) => {
    await page.route('http://localhost:8080/vehicles/11111111-1111-1111-1111-111111111111', (route) =>
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(vehicle) })
    )

    await page.goto('/editor/vehicles/11111111-1111-1111-1111-111111111111')
    const field = page.locator('#vehicle-country')
    await expect(field).toHaveValue('Japan')

    await field.fill('zzzz')
    await field.press('Tab')

    await expect(field).toHaveValue('Japan')
    await expect(page.getByRole('button', { name: 'Save' })).toBeDisabled()
  })

  test('previews the cutout when both cover assets exist', async ({ page }) => {
    await page.route('http://localhost:8080/vehicles/11111111-1111-1111-1111-111111111111', (route) =>
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          ...vehicle,
          cutout_key: 'vehicles/2026/10/cutout.png',
          silhouette_key: 'vehicles/2026/10/silhouette.png'
        })
      })
    )

    await page.goto('/editor/vehicles/11111111-1111-1111-1111-111111111111')
    await expect(page.getByText('Generated cover is assigned.')).toBeVisible()
    await expect(page.locator('img[alt="Honda Civic Type R"]')).toHaveCount(1)
    await expect(page.locator('.silhouette')).toHaveCount(0)
  })

  test('shows a cover thumbnail on each listing row', async ({ page }) => {
    const cover = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 60"><rect width="100" height="60" fill="#888"/></svg>`

    await page.route('https://assets.grd0.net/**', (route) =>
      route.fulfill({ status: 200, contentType: 'image/svg+xml', body: cover })
    )

    await page.route('http://localhost:8080/vehicles', (route) =>
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          { ...vehicle, cutout_key: 'vehicles/cutout.png', silhouette_key: 'vehicles/sil.png' },
          { ...vehicle, id: '22222222-2222-2222-2222-222222222222', cutout_key: null, silhouette_key: 'vehicles/sil.png' },
          { ...vehicle, id: '33333333-3333-3333-3333-333333333333', cutout_key: null, silhouette_key: null }
        ])
      })
    )

    await page.goto('/editor/vehicles')
    await expect(page.getByText('Honda Civic Type R').first()).toBeVisible()
    await expect(page.locator('main img')).toHaveCount(1)
    await expect(page.locator('main .silhouette')).toHaveCount(1)
    await expect(page.getByText('Unknown')).toHaveCount(1)

    const thumbnail = (await page.locator('main img').first().boundingBox())!
    const title = (await page.getByText('Honda Civic Type R').first().boundingBox())!
    expect(title.x - (thumbnail.x + thumbnail.width)).toBeLessThan(24)
  })
})
