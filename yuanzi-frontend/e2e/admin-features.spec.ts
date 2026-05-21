import { test, expect } from '@playwright/test';

// Test against deployed environment
const BASE_URL = 'https://babygarden.pages.dev';
const ADMIN_PHONE = '13800138000';
const ADMIN_PASSWORD = 'yuanzi123';

// Skip tests if deployed env is not ready
test.beforeEach(async ({ page }) => {
  const response = await page.goto(`${BASE_URL}/admin/login`);
  if (!response || response.status() !== 200) {
    test.skip('Deployed environment not available');
  }
});

test.describe('Admin Pages Load (Issue #38)', () => {
  test('admin login page loads', async ({ page }) => {
    await expect(page.locator('input[placeholder*="手机号"]')).toBeVisible();
    await expect(page.locator('input[placeholder*="密码"]')).toBeVisible();
    await expect(page.locator('button:has-text("登录")')).toBeVisible();
  });

  test('admin dashboard loads after login', async ({ page }) => {
    // Login with test credentials
    await page.fill('input[placeholder*="手机号"]', ADMIN_PHONE);
    await page.fill('input[placeholder*="密码"]', ADMIN_PASSWORD);
    await page.click('button:has-text("登录")');
    
    // Wait for dashboard
    await page.waitForURL('**/admin/dashboard', { timeout: 10000 });
    
    // Verify dashboard elements
    await expect(page.locator('text=总用户数')).toBeVisible();
    await expect(page.locator('text=总家庭数')).toBeVisible();
  });

  test('user management page has add button', async ({ page }) => {
    // Login first
    await page.fill('input[placeholder*="手机号"]', ADMIN_PHONE);
    await page.fill('input[placeholder*="密码"]', ADMIN_PASSWORD);
    await page.click('button:has-text("登录")');
    await page.waitForURL('**/admin/dashboard');
    
    // Navigate to users page
    await page.click('text=用户管理');
    await page.waitForURL('**/admin/users');
    
    // Verify add button exists
    await expect(page.locator('button:has-text("新增")')).toBeVisible();
  });
});

test.describe('Photo Management Page (Issue #39)', () => {
  test('photo page loads with upload button', async ({ page }) => {
    // Login
    await page.fill('input[placeholder*="手机号"]', ADMIN_PHONE);
    await page.fill('input[placeholder*="密码"]', ADMIN_PASSWORD);
    await page.click('button:has-text("登录")');
    await page.waitForURL('**/admin/dashboard');
    
    // Navigate to photos
    await page.click('text=照片管理');
    await page.waitForURL('**/admin/photos');
    
    // Verify upload button
    await expect(page.locator('button:has-text("上传")')).toBeVisible();
  });
});
