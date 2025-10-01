import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_SSO_SERVICE_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const authService = {
  async login(email, password) {
    const response = await api.post('/api/v1/auth/login', {
      email,
      password,
    });
    return response.data;
  },

  async register(userData) {
    const response = await api.post('/api/v1/auth/register', userData);
    return response.data;
  },

  async getProfile(token) {
    const response = await api.get('/api/v1/profile', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data.user;
  },

  async deleteAccount(token) {
    const response = await api.delete('/api/v1/account', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  },
};
