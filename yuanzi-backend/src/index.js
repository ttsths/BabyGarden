/**
 * YuanziBackend — Cloudflare Container
 *
 * Go 后端运行在 Cloudflare Containers 上，
 * Worker 作为反向代理，将 HTTP 请求转发到容器。
 */
import { Container } from '@cloudflare/containers';

export class YuanziBackend extends Container {
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
