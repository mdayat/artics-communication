interface UserResponse {
  id: string;
  email: string;
  name: string;
  role: string;
  created_at: string;
}

interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

interface LoginRequest {
  email: string;
  password: string;
}

export type { UserResponse, RegisterRequest, LoginRequest };
