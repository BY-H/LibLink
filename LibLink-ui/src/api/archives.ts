import request from '@/utils/request'

export function getArchives(params: object) {
    return request({
        url: '/api/archives/list',
        method: 'get',
        params
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

export function batchImportArchives(file: File) {
    const formData = new FormData()
    formData.append('file', file)

    return request({
        url: '/api/archives/batch_import',
        method: 'post',
        headers: {
            'Content-Type': 'multipart/form-data'
        },
        data: formData
    })
}
