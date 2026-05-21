import { test, expect, type Page, type Route } from '@playwright/test';

const token = 'e2e-token';
const babyId = '400e8400-e29b-41d4-a716-446655440000';
const familyId = '200e8400-e29b-41d4-a716-446655440000';

type RecordItem = {
  id: string;
  baby_id: string;
  type: string;
  started_at: string;
  ended_at?: string;
  content: Record<string, unknown>;
  note?: string;
};

const seedRecords: RecordItem[] = [
  { id: 'record-feed-seed', baby_id: babyId, type: 'feeding', started_at: new Date().toISOString(), content: { type: 'formula', amount: 120 }, note: '配方奶' },
  { id: 'record-sleep-seed', baby_id: babyId, type: 'sleep', started_at: new Date(Date.now() - 7200000).toISOString(), ended_at: new Date(Date.now() - 3600000).toISOString(), content: { quality: 'good', location: 'crib' }, note: '睡眠' },
  { id: 'record-excretion-seed', baby_id: babyId, type: 'excretion', started_at: new Date().toISOString(), content: { type: 'poop', amount: 'normal' }, note: '排泄' },
  { id: 'record-temp-seed', baby_id: babyId, type: 'temperature', started_at: new Date().toISOString(), content: { value: 36.7 }, note: '测温' },
];

function apiResponse(data: unknown) {
  return { code: 0, msg: 'success', data };
}

function listResponse<T>(list: T[]) {
  return apiResponse({ list, pagination: { page: 1, page_size: 20, total: list.length, total_pages: 1 } });
}

async function mockApi(page: Page) {
  let records = [...seedRecords];
  let photos = [{
    id: 'photo-seed',
    url: 'data:image/svg+xml,%3Csvg xmlns=%22http://www.w3.org/2000/svg%22 width=%22200%22 height=%22200%22%3E%3Crect width=%22200%22 height=%22200%22 fill=%22%23ff9a8b%22/%3E%3Ctext x=%2240%22 y=%22105%22 fill=%22white%22%3E小园子%3C/text%3E%3C/svg%3E',
    thumb_url: '',
    like_count: 1,
    comment_count: 1,
    liked_by_me: false,
    taken_at: new Date().toISOString(),
    description: '小园子照片',
  }];
  let comments = [{ id: 'comment-seed', photo_id: 'photo-seed', nickname: '奶奶', content: '真可爱', created_at: new Date().toISOString() }];

  await page.route('**/api/v1/**', async (route: Route) => {
    const request = route.request();
    const url = new URL(request.url());
    const path = url.pathname.replace('/api/v1', '');
    const method = request.method();

    if (path === '/auth/password-login' && method === 'POST') {
      await route.fulfill({ json: apiResponse({ access_token: token, refresh_token: 'refresh-token' }) });
      return;
    }
    if ((path === '/auth/login' || path === '/auth/send-code') && method === 'POST') {
      await route.fulfill({ json: path === '/auth/login' ? apiResponse({ access_token: token, refresh_token: 'refresh-token' }) : apiResponse({ expires_in: 300 }) });
      return;
    }
    if (path === '/user/profile') {
      await route.fulfill({ json: apiResponse({ id: '100e8400-e29b-41d4-a716-446655440000', phone: '13800138000', username: 'mom', nickname: '妈妈' }) });
      return;
    }
    if (path === '/baby') {
      await route.fulfill({ json: apiResponse([{ id: babyId, family_id: familyId, name: '小园子', birthday: '2024-01-01', gender: 2 }]) });
      return;
    }
    if (path === '/stats/daily') {
      await route.fulfill({ json: apiResponse({ feeding: { count: 3, total_amount: 360 }, sleep: { count: 2, total_hours: 10.5 }, diaper: { count: 2 }, temperature: { count: 1, latest: 36.7 } }) });
      return;
    }
    if (path === '/stats/summary') {
      await route.fulfill({ json: apiResponse({ range: url.searchParams.get('range') || 'week', dates: ['2026-05-15', '2026-05-16'], daily_avg_sleep_hours: [10, 11], daytime_single_sleep_hours: [1.5, 2], daily_avg_milk_amount: [360, 420], summary: { avg_daily_sleep_hours: 10.5, avg_daytime_single_sleep_hours: 1.75, avg_daily_milk_amount: 390 } }) });
      return;
    }
    if (path === '/record' && method === 'GET') {
      const type = url.searchParams.get('type');
      await route.fulfill({ json: listResponse(type ? records.filter((record) => record.type === type) : records) });
      return;
    }
    if (path === '/record' && method === 'POST') {
      const body = request.postDataJSON() as Partial<RecordItem>;
      const record = { id: `record-${Date.now()}`, baby_id: babyId, started_at: new Date().toISOString(), content: {}, ...body } as RecordItem;
      records = [record, ...records];
      await route.fulfill({ json: apiResponse(record) });
      return;
    }
    if (path.startsWith('/record/') && method === 'PUT') {
      const id = path.split('/').pop();
      const body = request.postDataJSON() as Partial<RecordItem>;
      records = records.map((record) => record.id === id ? { ...record, ...body } : record);
      await route.fulfill({ json: apiResponse(records.find((record) => record.id === id)) });
      return;
    }
    if (path.startsWith('/record/') && method === 'DELETE') {
      const id = path.split('/').pop();
      records = records.filter((record) => record.id !== id);
      await route.fulfill({ json: apiResponse({ ok: true }) });
      return;
    }
    if (path.startsWith('/record/') && method === 'GET') {
      const id = path.split('/').pop();
      await route.fulfill({ json: apiResponse(records.find((record) => record.id === id)) });
      return;
    }
    if (path === `/family/${familyId}`) {
      await route.fulfill({ json: apiResponse({ id: familyId, name: '小园子的家', invite_code: 'ABC12345', is_paid: false, storage_limit: 1073741824, storage_used: 1024 }) });
      return;
    }
    if (path === `/family/${familyId}/members`) {
      await route.fulfill({ json: apiResponse([{ user_id: 'mom', nickname: '妈妈', role: 'admin' }, { user_id: 'dad', nickname: '爸爸', role: 'member' }]) });
      return;
    }
    if (path === '/photo' && method === 'GET') {
      await route.fulfill({ json: listResponse(photos) });
      return;
    }
    if (path === '/photo/upload-url' && method === 'POST') {
      await route.fulfill({ json: apiResponse({ upload_url: 'https://upload.local/photo', photo_id: 'photo-uploaded', upload_headers: {} }) });
      return;
    }
    if (path === '/photo/confirm' && method === 'POST') {
      photos = [{ ...photos[0], id: 'photo-uploaded', description: '上传照片' }, ...photos];
      await route.fulfill({ json: apiResponse({ ok: true }) });
      return;
    }
    if (path.endsWith('/like') && method === 'POST') {
      photos = photos.map((photo) => ({ ...photo, liked_by_me: true, like_count: photo.like_count + 1 }));
      await route.fulfill({ json: apiResponse({ liked_by_me: true, like_count: photos[0].like_count, comment_count: photos[0].comment_count }) });
      return;
    }
    if (path.endsWith('/comments') && method === 'GET') {
      await route.fulfill({ json: listResponse(comments) });
      return;
    }
    if (path.endsWith('/comments') && method === 'POST') {
      const body = request.postDataJSON() as { content: string };
      comments = [{ id: `comment-${Date.now()}`, photo_id: 'photo-seed', nickname: '妈妈', content: body.content, created_at: new Date().toISOString() }, ...comments];
      await route.fulfill({ json: apiResponse(comments[0]) });
      return;
    }
    if (path === '/ai/chats' || path === '/ai/history') {
      await route.fulfill({ json: listResponse([{ id: 'ai-seed', question: '近一周趋势如何？', answer: '睡眠平稳，奶量正常，排泄规律。', created_at: new Date().toISOString() }]) });
      return;
    }
    if (path === '/ai/chat/stream') {
      await route.fulfill({
        contentType: 'text/event-stream',
        body: [
          'event: delta\ndata: {"delta":"睡眠平稳，"}\n\n',
          'event: delta\ndata: {"delta":"奶量正常，排泄规律。"}\n\n',
          'event: done\ndata: {"id":"ai-new","answer":"睡眠平稳，奶量正常，排泄规律。"}\n\n',
        ].join(''),
      });
      return;
    }
    await route.fulfill({ status: 404, json: apiResponse({ path, method }) });
  });

  await page.route('https://upload.local/photo', async (route) => {
    await route.fulfill({ status: 200, body: 'ok' });
  });
}

test.describe('BabyGarden 多端登录与小园子真实链路', () => {
  test.beforeEach(async ({ page }) => {
    await mockApi(page);
  });

  test('用户名密码登录 PC、APP、管理端', async ({ page }) => {
    await page.goto('/login');
    await page.getByPlaceholder('请输入手机号或用户名').fill('mom');
    await page.getByPlaceholder('请输入密码').fill('yuanzi123');
    await page.getByRole('button', { name: '登录' }).click();
    await expect(page).toHaveURL(/\/$/);
    await expect(page.getByText('数据源：后端 API')).toBeVisible();
    await expect(page.getByText('配方奶')).toBeVisible();
    await expect(page.getByText('小园子照片')).toBeVisible();
    await expect(page.getByText('本地示例')).toHaveCount(0);

    await page.goto('/app/login');
    await page.getByRole('button', { name: '账号密码' }).click();
    await page.getByPlaceholder('请输入手机号或用户名').fill('13800138000');
    await page.getByPlaceholder('请输入密码').fill('yuanzi123');
    await page.getByRole('button', { name: '登录' }).click();
    await expect(page).toHaveURL(/\/app$/);

    await page.route('**/api/v1/admin/login', async (route) => {
      await route.fulfill({ json: apiResponse({ token: 'admin-token', expires_in: 86400, user: { id: 'admin', phone: '13800138000', nickname: '妈妈', is_admin: 1 } }) });
    });
    await page.route('**/api/v1/admin/stats/**', async (route) => route.fulfill({ json: apiResponse({ users: 6, families: 1, babies: 1, photos: 1, records: 4 }) }));
    await page.goto('/admin/login');
    await page.getByPlaceholder('请输入管理员手机号').fill('13800138000');
    await page.getByPlaceholder('请输入密码').fill('yuanzi123');
    await page.locator('button[type="submit"]').click();
    await expect(page).toHaveURL(/\/admin\/dashboard/);
  });

  test('首页统计、记录 CRUD、AI 问答和照片互动', async ({ page }) => {
    await page.goto('/app/login');
    await page.getByRole('button', { name: '账号密码' }).click();
    await page.getByPlaceholder('请输入手机号或用户名').fill('mom');
    await page.getByPlaceholder('请输入密码').fill('yuanzi123');
    await page.getByRole('button', { name: '登录' }).click();
    await expect(page.getByText('今日喝奶')).toBeVisible();

    for (const [label, value] of [['喂养', 'feeding'], ['睡眠', 'sleep'], ['排泄', 'excretion'], ['测温', 'temperature']] as const) {
      await page.getByRole('button', { name: '记录', exact: true }).click();
      await page.locator('select').nth(1).selectOption(value);
      await page.getByPlaceholder('备注').fill(`${label} E2E 记录`);
      await page.getByRole('button', { name: '保存记录' }).click();
      await expect(page.getByText(`${label} E2E 记录`)).toBeVisible();
      await page.waitForTimeout(500);
      await page.getByRole('button', { name: '修改' }).first().click();
      await page.getByPlaceholder('备注').fill(`${label} E2E 修改`);
      await page.getByRole('button', { name: '保存修改' }).click();
      await expect(page.getByText(`${label} E2E 修改`)).toBeVisible();
      await page.getByRole('button', { name: '删除' }).first().click();
    }

    await page.getByRole('button', { name: '统计' }).click();
    await page.getByRole('button', { name: '月' }).click();
    await expect(page.getByText('日均喝奶量')).toBeVisible();
    await page.getByRole('button', { name: '自定义' }).click();
    await page.getByRole('button', { name: '更新区间' }).click();

    await page.getByRole('button', { name: 'AI' }).click();
    await page.getByRole('button', { name: '分析近一周趋势' }).click();
    await expect(page.getByText('奶量正常')).toBeVisible();
    await page.getByPlaceholder('输入你的问题').fill('今晚睡前怎么安排？');
    await page.getByRole('button', { name: '提问' }).click();
    await expect(page.getByText('历史会话')).toBeVisible();

    await page.getByRole('button', { name: '照片' }).click();
    await page.locator('input[type="file"]').setInputFiles({ name: 'yuanzi.jpg', mimeType: 'image/jpeg', buffer: Buffer.from('fake-image') });
    await expect(page.getByText('上传成功')).toBeVisible();
    await page.getByRole('button', { name: '点赞' }).first().click();
    await page.getByRole('button', { name: '查看评论' }).first().click();
    await page.getByPlaceholder('写评论').first().fill('照片评论 E2E');
    await page.getByRole('button', { name: '发送' }).first().click();
    await expect(page.getByText('照片评论 E2E')).toBeVisible();
  });
});
