import { test, expect } from '@playwright/test';

const BASE_URL = process.env.E2E_BASE_URL || 'https://babygarden.pages.dev';
const ADMIN_PHONE = process.env.E2E_ADMIN_PHONE || '13800000001';
const ADMIN_PASSWORD = process.env.E2E_ADMIN_PASSWORD || 'admin123';

async function loginAsAdmin(page: any) {
  await page.goto(`${BASE_URL}/admin/login`);
  await page.fill('input[placeholder="手机号"]', ADMIN_PHONE);
  await page.fill('input[placeholder="密码"]', ADMIN_PASSWORD);
  await page.click('button:has-text("登录")');
  await page.waitForURL('**/admin/dashboard', { timeout: 15000 });
}

test.describe('Issue #38 - 后台管理新增/编辑功能', () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
  });

  test('用户管理页面 - 新增按钮存在', async ({ page }) => {
    await page.click('text=用户管理');
    await page.waitForURL('**/admin/users');
    
    // 验证新增按钮存在
    await expect(page.locator('button:has-text("新增")')).toBeVisible({ timeout: 10000 });
  });

  test('宝宝管理页面 - 新增按钮存在', async ({ page }) => {
    await page.click('text=宝宝管理');
    await page.waitForURL('**/admin/babies');
    
    // 验证新增按钮存在
    await expect(page.locator('button:has-text("新增")')).toBeVisible({ timeout: 10000 });
  });

  test('家庭管理页面 - 新增按钮存在', async ({ page }) => {
    await page.click('text=家庭管理');
    await page.waitForURL('**/admin/families');
    
    // 验证新增按钮存在
    await expect(page.locator('button:has-text("新增")')).toBeVisible({ timeout: 10000 });
  });

  test('记录管理页面 - 新增按钮存在', async ({ page }) => {
    await page.click('text=记录管理');
    await page.waitForURL('**/admin/records');
    
    // 验证新增按钮存在
    await expect(page.locator('button:has-text("新增")')).toBeVisible({ timeout: 10000 });
  });
});

test.describe('Issue #39 - 照片上传/下载功能', () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
  });

  test('照片管理页面 - 上传按钮存在', async ({ page }) => {
    await page.click('text=照片管理');
    await page.waitForURL('**/admin/photos');
    
    // 验证上传按钮存在
    await expect(page.locator('button:has-text("上传")')).toBeVisible({ timeout: 10000 });
  });

  test('照片管理页面 - 批量下载按钮存在', async ({ page }) => {
    await page.click('text=照片管理');
    await page.waitForURL('**/admin/photos');
    
    // 等待表格加载
    await page.waitForSelector('.ant-table-row', { timeout: 10000 });
    
    // 选择第一张照片
    await page.click('.ant-table-row:first-child .ant-checkbox-input');
    
    // 验证批量下载按钮可用
    await expect(page.locator('button:has-text("批量下载")')).toBeEnabled();
  });

  test('照片管理页面 - 单张下载按钮存在', async ({ page }) => {
    await page.click('text=照片管理');
    await page.waitForURL('**/admin/photos');
    
    // 等待表格加载
    await page.waitForSelector('.ant-table-row', { timeout: 10000 });
    
    // 验证下载按钮存在
    await expect(page.locator('button:has-text("下载")').first()).toBeVisible();
  });
});
