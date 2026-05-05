import request from "../utils/request";

// 获取仪表盘统计数据
export function getDashboardStats() {
  return request.get("/admin/dashboard/stats");
}

// 获取仪表盘图表分析数据
export function getDashboardAnalytics() {
  return request.get("/admin/dashboard/analytics");
}

// 导出管理端仪表盘报告
export function exportDashboardReport() {
  return request.get("/admin/dashboard/export", {
    responseType: "blob",
  });
}

// 获取用户列表（分页）
export function getUsers(params: any) {
  return request.get("/admin/users", { params });
}

// 获取单个用户详情
export function getUser(id: number | string) {
  return request.get(`/admin/users/${id}`);
}

// 创建新用户
export function createUser(data: any) {
  return request.post("/admin/users", data);
}

// 更新用户信息
export function updateUser(id: number | string, data: any) {
  return request.put(`/admin/users/${id}`, data);
}

// 删除用户
export function deleteUser(id: number | string) {
  return request.delete(`/admin/users/${id}`);
}
