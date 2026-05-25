// d:\gowork\max2\web\src\utils\websocketReconnect.js
import { storage } from '@/utils/storage';

/**
 * WebSocket 自动重连工具
 * 从 localStorage 获取用户信息并重新连接
 */
export const reconnectWebSocket = (websocket) => {
  // 从 localStorage 获取用户信息
  const userData = storage.local.get('user');
  if (!userData) {
    console.log('[WS Reconnect] 未找到用户信息，跳过 WebSocket 重连');
    return;
  }

  try {
    const user = typeof userData === 'string' ? JSON.parse(userData) : userData;
    const { username, token, node } = user;

    if (!node) {
      console.log('[WS Reconnect] 未找到 node 信息，跳过 WebSocket 重连');
      return;
    }

    // 构建 WebSocket URL
    const wsUrl = `ws://${node}/enter?user_name=${encodeURIComponent(username)}&token=${encodeURIComponent(token)}`;
    console.log('[WS Reconnect] 尝试重连 WebSocket:', wsUrl);

    // 如果当前没有连接、连接已关闭或正在关闭，才重新初始化
    // 排除 CONNECTING 状态，避免中断正在建立的连接
    if (!websocket.ws || 
        websocket.ws.readyState === WebSocket.CLOSED || 
        websocket.ws.readyState === WebSocket.CLOSING) {
      websocket.init({
        url: wsUrl,
        reconnectInterval: 5000
      });
      websocket.connect();
    } else {
      console.log('[WS Reconnect] WebSocket 已连接或正在连接');
    }
  } catch (error) {
    console.error('[WS Reconnect] 重连 WebSocket 失败:', error);
  }
};