import request from '/@/utils/request';
export function getUserList(query:Object) {
    return request({
        url: '/api/v1/user/list',
        method: 'get',
        params:query
    })
}

export function updateUserPoint(query:Object) {
    return request({
        url: '/api/v1/user/updatePoint',
        method: 'post',
        data:query
    })
}

export function updateUserPassword(query:Object) {
    return request({
        url: '/api/v1/user/updatePassword',
        method: 'post',
        data:query
    })
}
