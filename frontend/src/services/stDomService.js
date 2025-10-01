const API_BASE_URL = process.env.REACT_APP_ST_DOM_SERVICE_URL || 'http://localhost:8081/api/v1';

class StDomService {
  constructor() {
    this.baseURL = API_BASE_URL;
  }

  async makeRequest(endpoint, options = {}) {
    const token = localStorage.getItem('token');
    
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...(token && { Authorization: `Bearer ${token}` }),
        ...options.headers,
      },
      ...options,
    };

    const response = await fetch(`${this.baseURL}${endpoint}`, config);
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  // Get all student dormitories (public)
  async getAllStDoms() {
    return this.makeRequest('/st_doms/');
  }

  // Get specific student dormitory (public)
  async getStDom(id) {
    return this.makeRequest(`/st_doms/${id}`);
  }

  // Get rooms for a specific dormitory (public)
  async getStDomRooms(id) {
    return this.makeRequest(`/st_doms/${id}/rooms`);
  }

  // Create new student dormitory (admin only)
  async createStDom(stDomData) {
    return this.makeRequest('/st_doms/', {
      method: 'POST',
      body: JSON.stringify(stDomData),
    });
  }

  // Update student dormitory (admin only)
  async updateStDom(id, stDomData) {
    return this.makeRequest(`/st_doms/${id}`, {
      method: 'PUT',
      body: JSON.stringify(stDomData),
    });
  }

  // Delete student dormitory (admin only)
  async deleteStDom(id) {
    return this.makeRequest(`/st_doms/${id}`, {
      method: 'DELETE',
    });
  }

  // Get all rooms (public)
  async getAllRooms() {
    return this.makeRequest('/sobas/');
  }

  // Get specific room (public)
  async getRoom(id) {
    return this.makeRequest(`/sobas/${id}`);
  }

  // Create new room (admin only)
  async createRoom(roomData) {
    return this.makeRequest('/sobas/', {
      method: 'POST',
      body: JSON.stringify(roomData),
    });
  }

  // Update room (admin only)
  async updateRoom(id, roomData) {
    return this.makeRequest(`/sobas/${id}`, {
      method: 'PUT',
      body: JSON.stringify(roomData),
    });
  }

  // Delete room (admin only)
  async deleteRoom(id) {
    return this.makeRequest(`/sobas/${id}`, {
      method: 'DELETE',
    });
  }
}

export const stDomService = new StDomService();
