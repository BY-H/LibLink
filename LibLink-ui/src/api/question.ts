import request from '@/utils/request'

export function summary() {
    return request({
        url: '/api/questions/summary',
        method: 'GET'
    })
}

export function getQuestions(params: any) {
    return request({
        url: '/api/questions/list',
        method: 'GET',
        params
    })
}

export function addQuestion(data: any) {
    return request({
        url: '/api/questions/add',
        method: 'POST',
        data
    })
}

export function editQuestion(data: any) {
    return request({
        url: '/api/questions/edit',
        method: 'PATCH',
        data
    })
}

export function deleteQuestion(data: any) {
    return request({
        url: '/api/questions/delete',
        method: 'DELETE',
        data
    })
}

export function getTags() {
    return request({
        url: '/api/questions/tags',
        method: 'GET'
    })
}