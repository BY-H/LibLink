import request from '@/utils/request'

export function login(data: object) {
    return request({
        url: '/login',
        method: 'POST',
        data
    })
}

export function register(data: object) {
    return request({
        url: '/register',
        method: 'POST',
        data
    })
}

export function summary() {
    return request({
        url: '/api/users/summary',
        method: 'GET'
    })
}
