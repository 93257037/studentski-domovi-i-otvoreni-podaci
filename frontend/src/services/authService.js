import axios from 'axios';

// API Gatewaysvi servisi su dostupni kroz jedan gateway
const API_BASE_URL = process.env.REACT_APP_API_GATEWAY_URL || 'http://localhost';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const authService = {
  // prijavljuje korisnika sa email-om i lozinkom
  async login(email, password) {
    const response = await api.post('/api/v1/auth/login', {
      email,
      password,
    });
    return response.data;
  },

  // registruje novog korisnika
  async register(userData) {
    const response = await api.post('/api/v1/auth/register', userData);
    return response.data;
  },

  // dobija profil korisnika na osnovu tokena
  async getProfile(token) {
    const response = await api.get('/api/v1/profile', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data.user;
  },

  // brise nalog korisnika
  async deleteAccount(token) {
    const response = await api.delete('/api/v1/account', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return response.data;
  },
};
