import request from '@/utils/request'

export function ping() {
    return request({
        url: '/ping_without_login',
        method: 'GET'
    })
}
