export interface User {
  id: string;
  email: string;
}

export interface Bookmark {
  id: string;
  url: string;
  title: string;
  description: string;
  tags: string[];
  createdAt: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface ApiError {
  message: string;
  status: number;
}