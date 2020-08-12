import axios from './api'

const Role = {
    /**
     * 全部角色
     */
    GetAll() {
        return axios.get(`/v1/role`);
    },

    /**
     * 添加
     * @param {*} data 添加角色信息
     */
    Add(data){
        return axios.post('/v1/role', data)
    },

    /**
     * 删除角色
     * @param {*} id 
     */
    Del(id){
        return axios.delete(`/v1/role?id=${id}`)
    },

    /**
     * 修改角色
     * @param {*} data 
     */
    Save(data){
        return axios.put('/v1/role', data)
    },

}

export {
    Role
}
