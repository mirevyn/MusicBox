import request from "../utils/request";
import type {
  LoginDTO,
  RegisterDTO,
  ChangePasswordDTO,
  ApiResponse,
  LoginResponse,
  User,
} from "@/types/api";

// 登录
export function login(data: LoginDTO): Promise<ApiResponse<LoginResponse>> {
  return request.post("/auth/login", data);
}

// 注册
export function register(data: RegisterDTO): Promise<ApiResponse<null>> {
  return request.post("/auth/register", data);
}

// 获取用户信息
export function getUserInfo(): Promise<ApiResponse<User>> {
  return request.get("/user/profile");
}

// 更新用户密码
export function updatePassword(
  data: ChangePasswordDTO
): Promise<ApiResponse<null>> {
  return request.put("/user/password", data);
}

// 上传用户头像
export function uploadAvatar(
  formData: FormData
): Promise<ApiResponse<{ avatarUrl: string }>> {
  return request.post("/user/avatar", formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
}
