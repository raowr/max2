import axios from 'axios';
import { ElMessage, ElMessageBox } from 'element-plus';

// 配置新建一个 axios 实例
const service = axios.create({
  baseURL: getRandomBaseUrl(),
  timeout: 50000,
  headers: { 'Content-Type': 'application/json' },
  paramsSerializer: {
    serialize(params) {
      return new URLSearchParams(params).toString();
    },
  },
});

// 添加请求拦截器
service.interceptors.request.use(
  (config) => {
    // 在发送请求之前做些什么 - 添加 token
    const token = sessionStorage.getItem('token');
    if (token) {
      config.headers = config.headers || {};
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    // 对请求错误做些什么
    console.error('请求错误:', error);
    return Promise.reject(error);
  }
);

// 添加响应拦截器
service.interceptors.response.use(
  (response) => {
    // 对响应数据做点什么
    const res = response.data;
    const code = res.code;
    
    if (code === 401) {
      // 登录状态已过期
      ElMessageBox.alert('登录状态已过期，请重新登录', '提示', {
        confirmButtonText: '确定'
      }).then(() => {
        sessionStorage.clear();
        window.location.href = '/';
      }).catch(() => {});
      return Promise.reject(new Error('登录状态已过期'));
    } else if (code !== 0 && code !== 200) {
      // 业务错误
      ElMessage.error(res.message || '请求失败');
      return res;
      // return Promise.reject(new Error(res.message || '请求失败'));
    } else {
      return res;
    }
  },
  (error) => {
    // 对响应错误做点什么
    if (error.message && error.message.indexOf('timeout') !== -1) {
      ElMessage.error('网络超时');
    } else if (error.message === 'Network Error') {
      ElMessage.error('网络连接错误');
    } else {
      if (error.response && error.response.data) {
        ElMessage.error(error.response.statusText || '请求失败');
      } else {
        ElMessage.error('接口路径找不到');
      }
    }
    return Promise.reject(error);
  }
);

// 从分号分隔的多个 URL 中随机选择一个
function getRandomBaseUrl() {
    const urls = import.meta.env.VITE_USER_URL;
    if (!urls) {
        console.warn('VITE_USER_URL 环境变量未配置');
        return '/api';
    }
    // 用分号分割 URL
    const urlList = urls.split(';').map(url => url.trim()).filter(url => url);
    if (urlList.length === 0) {
        console.warn('VITE_USER_URL 没有有效 URL');
        return '/api';
    }
    // 随机选择一个
    const randomIndex = Math.floor(Math.random() * urlList.length);
    return urlList[randomIndex];
}

// 导出 axios 实例
export default service;