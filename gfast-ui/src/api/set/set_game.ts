import request from '/@/utils/request';
export function getSetGameList(query:Object) {
    return request({
        url: '/api/v1/set/setGame/list',
        method: 'get',
        params:query
    })
}

export function updateSetGame(query:Object) {
    return request({
        url: '/api/v1/set/setGame/update',
        method: 'post',
        data:query
    })
}
