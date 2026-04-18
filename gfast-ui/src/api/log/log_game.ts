import request from '/@/utils/request';
export function getLogGameList(query:Object) {
    return request({
        url: '/api/v1/log/logGame/list',
        method: 'get',
        params:query
    })
}
