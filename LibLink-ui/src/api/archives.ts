import request from '@/utils/request'

export function getArchives() {
    return request({
        url: '/api/archives/list',
        method: 'get'
    })
}

export function addArchive(data: object) {
    return request({
        url: '/api/archives/add',
        method: 'post',
        data
    })
}

export function borrowArchive(params: object) {
    return request({
        url: '/api/archives/borrow',
        method: 'patch',
        params
    })
}

export function returnArchive(params: object) {
    return request({
        url: '/api/archives/return',
        method: 'patch',
        params
    })
}