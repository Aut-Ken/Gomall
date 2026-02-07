import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:3000';

test('User can buy a product immediately', async ({ page }) => {
  // Mock API responses
  await page.route('*/**/api/user/profile', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        code: 0,
        message: 'success',
        data: { id: 1, username: 'testuser', email: 'test@example.com' }
      })
    });
  });

  await page.route('*/**/api/product?*', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        code: 0,
        message: 'success',
        data: {
          list: [
            {
              id: 1,
              name: 'iPhone 15 Pro',
              description: 'Titanium design',
              price: 7999,
              stock: 100,
              category: 'Electronics',
              image_url: 'https://via.placeholder.com/200'
            }
          ],
          total: 1
        }
      })
    });
  });

  await page.route('*/**/api/product/1', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        code: 0,
        message: 'success',
        data: {
          id: 1,
          name: 'iPhone 15 Pro',
          description: 'Titanium design',
          price: 7999,
          stock: 100,
          category: 'Electronics',
          image_url: 'https://via.placeholder.com/200'
        }
      })
    });
  });

  await page.route('*/**/api/order', async route => {
    if (route.request().method() === 'POST') {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          code: 0,
          message: 'success',
          data: { id: 101, order_no: 'ORD12345' }
        })
      });
    } else {
      // GET list
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          code: 0,
          message: 'success',
          data: {
            list: [
              {
                id: 101,
                order_no: 'ORD12345',
                product_id: 1,
                product_name: 'iPhone 15 Pro',
                quantity: 1,
                total_price: 7999,
                status: 0,
                created_at: '2024-01-01T12:00:00Z'
              }
            ],
            total: 1
          }
        })
      });
    }
  });

  // Inject token
  await page.addInitScript(() => {
    localStorage.setItem('token', 'fake-jwt-token');
  });

  // 1. Go to Home Page
  await page.goto(BASE_URL);
  await expect(page).toHaveTitle(/GoMall/);

  // 2. Click on the first product card
  const firstProduct = page.locator('a[href^="/products/"]').first();
  await expect(firstProduct).toBeVisible();
  await firstProduct.click();

  // 3. Verify Product Detail Page
  await expect(page.locator('h1')).toHaveText('iPhone 15 Pro');

  // 4. Click "Buy Now" (立即购买)
  const buyNowBtn = page.getByRole('button', { name: '立即购买' });
  await expect(buyNowBtn).toBeVisible();
  await buyNowBtn.click();

  // 5. Verify Redirect to Orders
  await expect(page).toHaveURL(/.*\/orders/);

  // 6. Verify New Order in List
  await expect(page.locator('.order-item, [class*="order-item"], div:has-text("ORD12345")').first()).toBeVisible();
});
