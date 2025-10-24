// src/utils/storage.js
// 定义统一前缀（可根据项目名或功能自定义）
const PREFIX = 'max2_' // 例如：项目名为 max2，前缀设为 max2_

export const storage = {
    // localStorage 操作（带前缀）
    local: {
        set(key, value) {
            const prefixedKey = `${PREFIX}${key}` // 拼接前缀
            localStorage.setItem(prefixedKey, JSON.stringify(value))
        },
        get(key) {
            const prefixedKey = `${PREFIX}${key}`
            const data = localStorage.getItem(prefixedKey)
            return data ? JSON.parse(data) : null
        },
        remove(key) {
            const prefixedKey = `${PREFIX}${key}`
            localStorage.removeItem(prefixedKey)
        },
        clear() {
            // 清除时只删除带当前前缀的 key（避免误删其他数据）
            Object.keys(localStorage).forEach(key => {
                if (key.startsWith(PREFIX)) {
                    localStorage.removeItem(key)
                }
            })
        }
    },
    // sessionStorage 同理（带前缀）
    session: {
        set(key, value) {
            const prefixedKey = `${PREFIX}${key}`
            sessionStorage.setItem(prefixedKey, JSON.stringify(value))
        },
        get(key) {
            const prefixedKey = `${PREFIX}${key}`
            const data = sessionStorage.getItem(prefixedKey)
            return data ? JSON.parse(data) : null
        },
        remove(key) {
            const prefixedKey = `${PREFIX}${key}`
            sessionStorage.removeItem(prefixedKey)
        },
        clear() {
            Object.keys(sessionStorage).forEach(key => {
                if (key.startsWith(PREFIX)) {
                    sessionStorage.removeItem(key)
                }
            })
        }
    }
}