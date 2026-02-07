import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:3000';

test('User can seckill a product successfully', async ({ page }) => {
    // Mock API responses

    // 1. Mock Auth (User Profile)
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

    // 2. Mock Seckill Action (POST /api/seckill)
    await page.route('*/**/api/seckill', async route => {
        if (route.request().method() === 'POST') {
            await route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify({
                    code: 0, // Success
                    message: 'success',
                    data: {
                        order_no: 'SECKILL_ORDER_123',
                        success: true,
                        message: 'Queueing'
                    }
                })
            });
        } else {
            route.continue();
        }
    });

    // 3. Mock Order List (to verify redirection)
    await page.route('*/**/api/order', async route => {
        await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({
                code: 0,
                message: 'success',
                data: {
                    list: [
                        {
                            id: 201,
                            order_no: 'SECKILL_ORDER_123',
                            product_id: 1,
                            product_name: 'iPhone 15 Pro Max',
                            quantity: 1,
                            total_price: 5999,
                            status: 0,
                            created_at: new Date().toISOString()
                        }
                    ],
                    total: 1
                }
            })
        });
    });

    // Inject token
    await page.addInitScript(() => {
        localStorage.setItem('token', 'fake-jwt-token');
    });

    // 1. Go to Home Page
    await page.goto(BASE_URL);

    // 2. Navigate to Seckill Page
    // Assuming there is a link to Seckill page, or we can go directly
    const seckillLink = page.locator('a[href="/seckill"]');
    if (await seckillLink.count() > 0) {
        await seckillLink.first().click();
    } else {
        await page.goto(`${BASE_URL}/seckill`);
    }

    // 3. Verify Seckill Page
    await expect(page.locator('h1')).toContainText('限时秒杀');

    // 4. Find the first product with "Immediate Seckill" button (立即秒杀)
    // The button text varies: '立即秒杀' or '正在秒杀...' or '已售罄'
    // We want to click '立即秒杀'
    const seckillBtn = page.locator('button:has-text("立即秒杀")').first();
    await expect(seckillBtn).toBeVisible();

    // 5. Click Seckill Button
    await seckillBtn.click();

    // 6. Verify Redirect to Orders
    await expect(page).toHaveURL(/.*\/orders/);

    // 7. Verify Order is present
    await expect(page.locator('div:has-text("SECKILL_ORDER_123")').first()).toBeVisible();
});
