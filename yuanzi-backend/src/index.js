/**
 * YuanziBackend — Cloudflare Container
 *
 * Go 后端运行在 Cloudflare Containers 上，
 * Worker 作为反向代理，将 HTTP 请求转发到容器。
 */
import { Container } from '@cloudflare/containers';

export class YuanziBackend extends Container {
  // 容器监听端口（对应 Dockerfile EXPOSE 8080）
  defaultPort = 8080;

  /**
   * 处理所有进入容器的 HTTP 请求
   * containerFetch 自动发现容器的 EXPOSE 端口并转发
   */
  async fetch(request) {
    const rayId = request.headers.get('cf-ray') || 'unknown';
    console.log(`[${rayId}] ${request.method} ${request.url}`);

    try {
      return await this.containerFetch(request);
    } catch (err) {
      console.error(`Container fetch error:`, err);
      return new Response(JSON.stringify({
        error: 'Service temporarily unavailable',
        rayId,
      }), {
        status: 503,
        headers: { 'Content-Type': 'application/json' },
      });
    }
  }
}

/**
 * Worker entry — 必须 default export 才能以 ES Module 格式运行
 * 这样 DO binding 和 [[containers]] 才能被 wrangler 4.x 正确识别
 */
export default {
  async fetch(request, env) {
    // Cloudflare Containers don't expose [vars] as OS env vars.
    // We relay all Worker vars as a single base64-encoded JSON header.
    const vars = {};
    for (const [k, v] of Object.entries(env)) {
      if (typeof v === 'string' && v.length > 0) vars[k] = v;
    }
    const headers = new Headers(request.headers);
    if (Object.keys(vars).length > 0) {
      headers.set('X-Worker-Vars', btoa(JSON.stringify(vars)));
    }
    const modifiedRequest = new Request(request, { headers });

    const id = env.YUANZI_BACKEND.idFromName('default');
    const stub = env.YUANZI_BACKEND.get(id);
    return stub.fetch(modifiedRequest);
  },
};
