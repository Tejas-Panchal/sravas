import type { User, ApiResponse } from "../types/index.ts";
import client from "./client.ts";

export class AuthService {
  async getCurrentUser(): Promise<User | null> {
    try {
      const res = await client.get<ApiResponse<User>>("/users/get-user");
      return res.data.data;
    } catch {
      try {
        await this.refreshToken();
        const res = await client.get<ApiResponse<User>>("/users/get-user");
        return res.data.data;
      } catch {
        return null;
      }
    }
  }

  async refreshToken(): Promise<void> {
    await client.post("/users/refresh-token");
  }

  async login({ identifier, password }: { identifier: string; password: string }): Promise<User> {
    const res = await client.post<ApiResponse<{ user: User }>>("/users/login", {
      username: identifier,
      email: identifier,
      password,
    });
    return res.data.data.user;
  }

  async register(formData: FormData): Promise<void> {
    await client.post("/users/register", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
  }

  async logout(): Promise<void> {
    await client.post("/users/logout");
  }
}

const authService = new AuthService();

export default authService;
