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

  // Application endpoints

  // Create new application (user only)
  async createAplikacija(aplikacijaData) {
    return this.makeRequest('/aplikacije/', {
      method: 'POST',
      body: JSON.stringify(aplikacijaData),
    });
  }

  // Get all applications for current user
  async getMyAplikacije() {
    return this.makeRequest('/aplikacije/my');
  }

  // Get specific application
  async getAplikacija(id) {
    return this.makeRequest(`/aplikacije/${id}`);
  }

  // Update application (user can update their own)
  async updateAplikacija(id, aplikacijaData) {
    return this.makeRequest(`/aplikacije/${id}`, {
      method: 'PUT',
      body: JSON.stringify(aplikacijaData),
    });
  }

  // Delete application (user can delete their own)
  async deleteAplikacija(id) {
    return this.makeRequest(`/aplikacije/${id}`, {
      method: 'DELETE',
    });
  }

  // Admin application management

  // Get all applications (admin only)
  async getAllAplikacije() {
    return this.makeRequest('/aplikacije/');
  }

  // Approve application (admin only)
  async approveAplikacija(aplikacijaId, academicYear) {
    return this.makeRequest('/prihvacene_aplikacije/approve', {
      method: 'POST',
      body: JSON.stringify({
        aplikacija_id: aplikacijaId,
        academic_year: academicYear
      }),
    });
  }

  // Reject/Delete application (admin only)
  async rejectAplikacija(id) {
    return this.makeRequest(`/aplikacije/${id}`, {
      method: 'DELETE',
    });
  }

  // Room/Accepted Applications Management

  // Get all accepted applications (admin only)
  async getAllPrihvaceneAplikacije() {
    return this.makeRequest('/prihvacene_aplikacije/');
  }

  // Evict student from room (admin only)
  async evictStudent(userId, reason) {
    return this.makeRequest('/prihvacene_aplikacije/evict', {
      method: 'POST',
      body: JSON.stringify({
        user_id: userId,
        reason: reason
      }),
    });
  }

  // Get all payments with optional status filter (admin only)
  async getAllPayments(status = null) {
    const url = status ? `/payments/?status=${status}` : '/payments/';
    return this.makeRequest(url);
  }

  // Create payment (admin only)
  async createPayment(paymentData) {
    return this.makeRequest('/payments/', {
      method: 'POST',
      body: JSON.stringify(paymentData),
    });
  }

  // Search payments by student index (admin only)
  async searchPaymentsByIndex(indexPattern, status = null) {
    const url = status 
      ? `/payments/search?index=${indexPattern}&status=${status}`
      : `/payments/search?index=${indexPattern}`;
    return this.makeRequest(url);
  }

  // Get payments by room (admin only)
  async getPaymentsByRoom(sobaId) {
    return this.makeRequest(`/payments/room/${sobaId}`);
  }

  // Get payments by user (admin only)
  async getPaymentsByUser(userId) {
    return this.makeRequest(`/payments/user/${userId}`);
  }

  // Mark payment as paid (admin only)
  async markPaymentAsPaid(paymentId, paidAt = null) {
    return this.makeRequest(`/payments/${paymentId}/mark-paid`, {
      method: 'PATCH',
      body: JSON.stringify(paidAt ? { paid_at: paidAt } : {}),
    });
  }

  // Mark payment as unpaid (admin only)
  async markPaymentAsUnpaid(paymentId) {
    return this.makeRequest(`/payments/${paymentId}/mark-unpaid`, {
      method: 'PATCH',
    });
  }

  // User-specific endpoints

  // Get my accepted applications (current user)
  async getMyPrihvaceneAplikacije() {
    return this.makeRequest('/prihvacene_aplikacije/my');
  }

  // Checkout from room (current user)
  async checkoutFromRoom() {
    return this.makeRequest('/prihvacene_aplikacije/checkout', {
      method: 'POST',
    });
  }

  // Get payments by aplikacija ID
  async getPaymentsByAplikacija(aplikacijaId) {
    return this.makeRequest(`/payments/aplikacija/${aplikacijaId}`);
  }
}

export const stDomService = new StDomService();
