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
    // 测试目标：已部署的前端
    baseURL: 'https://babygarden.pages.dev',
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
