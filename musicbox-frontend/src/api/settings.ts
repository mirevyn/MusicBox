import request from "../utils/request";
import type { ApiResponse } from '@/types/api'

export interface SystemSettings {
    allowRegister: boolean
}

/**
 * 获取系统设置 (管理员)
 */
export const fetchSystemSettings = (): Promise<ApiResponse<SystemSettings>> => {
    return request.get('/admin/settings')
}

/**
 * 更新系统设置 (管理员)
 */
export const updateSystemSettings = (settings: SystemSettings): Promise<ApiResponse<null>> => {
    return request.put('/admin/settings', settings)
}

/**
 * 获取公开配置 (访客)
 */
export const fetchPublicConfig = (): Promise<ApiResponse<SystemSettings>> => {
    return request.get('/config')
}
