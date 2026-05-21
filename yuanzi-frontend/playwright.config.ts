import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  timeout: 30_000,
  expect: { timeout: 10_000 },
  
  retries: 0,
  workers: 1,
  
  reporter: [
    ['list'],
    ['json', { outputFile: 'e2e-results.json' }],
  ],

  use: {
    // 默认仍可跑部署站点；功能验收可用 E2E_BASE_URL 指向本地 dev server。
    baseURL: process.env.E2E_BASE_URL || 'https://babygarden.pages.dev',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
});
