// API Gateway svi servisi su dostupni kroz jedan gateway
const API_BASE_URL = process.env.REACT_APP_API_GATEWAY_URL || 'http://localhost/api/v1';

//komunikacija sa st_dom_service
//zahtevi vezani za domove, sobe, aplikacije i placanja
class StDomService {
  constructor() {
    this.baseURL = API_BASE_URL;
  }

  //HTTP zahtev sa automatskim dodavanjem JWT tokena
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

  // dobija sve studentske domove (javno)
  async getAllStDoms() {
    return this.makeRequest('/st_doms/');
  }

  // dobija odredjeni studentski dom po ID-u (javno)
  async getStDom(id) {
    return this.makeRequest(`/st_doms/${id}`);
  }

  // dobija sobe za odredjeni dom (javno)
  async getStDomRooms(id) {
    return this.makeRequest(`/st_doms/${id}/rooms`);
  }

  // kreira novi studentski dom (samo admin)
  async createStDom(stDomData) {
    return this.makeRequest('/st_doms/', {
      method: 'POST',
      body: JSON.stringify(stDomData),
    });
  }

  // azurira studentski dom (samo admin)
  async updateStDom(id, stDomData) {
    return this.makeRequest(`/st_doms/${id}`, {
      method: 'PUT',
      body: JSON.stringify(stDomData),
    });
  }

  // brise studentski dom (samo admin)
  async deleteStDom(id) {
    return this.makeRequest(`/st_doms/${id}`, {
      method: 'DELETE',
    });
  }

  // dobija sve sobe (javno)
  async getAllRooms() {
    return this.makeRequest('/sobas/');
  }

  // dobija odredjenu sobu po ID-u (javno)
  async getRoom(id) {
    return this.makeRequest(`/sobas/${id}`);
  }

  // kreira novu sobu (samo admin)
  async createRoom(roomData) {
    return this.makeRequest('/sobas/', {
      method: 'POST',
      body: JSON.stringify(roomData),
    });
  }

  // azurira sobu (samo admin)
  async updateRoom(id, roomData) {
    return this.makeRequest(`/sobas/${id}`, {
      method: 'PUT',
      body: JSON.stringify(roomData),
    });
  }

  // brise sobu (samo admin)
  async deleteRoom(id) {
    return this.makeRequest(`/sobas/${id}`, {
      method: 'DELETE',
    });
  }

  // endpoint-i za aplikacije

  // kreira novu aplikaciju za sobu (samo korisnik)
  async createAplikacija(aplikacijaData) {
    return this.makeRequest('/aplikacije/', {
      method: 'POST',
      body: JSON.stringify(aplikacijaData),
    });
  }

  // dobija sve aplikacije trenutnog korisnika
  async getMyAplikacije() {
    return this.makeRequest('/aplikacije/my');
  }

  // dobija odredjenu aplikaciju po ID-u
  async getAplikacija(id) {
    return this.makeRequest(`/aplikacije/${id}`);
  }

  // azurira aplikaciju (korisnik moze azurirati svoju)
  async updateAplikacija(id, aplikacijaData) {
    return this.makeRequest(`/aplikacije/${id}`, {
      method: 'PUT',
      body: JSON.stringify(aplikacijaData),
    });
  }

  // brise aplikaciju (korisnik moze brisati svoju)
  async deleteAplikacija(id) {
    return this.makeRequest(`/aplikacije/${id}`, {
      method: 'DELETE',
    });
  }

  // upravljanje aplikacijama za administratore

  // dobija sve aplikacije (samo admin)
  async getAllAplikacije() {
    return this.makeRequest('/aplikacije/');
  }

  // odobrava aplikaciju (samo admin)
  async approveAplikacija(aplikacijaId, academicYear) {
    return this.makeRequest('/prihvacene_aplikacije/approve', {
      method: 'POST',
      body: JSON.stringify({
        aplikacija_id: aplikacijaId,
        academic_year: academicYear
      }),
    });
  }

  // odbacuje/brise aplikaciju (samo admin)
  async rejectAplikacija(id) {
    return this.makeRequest(`/aplikacije/${id}`, {
      method: 'DELETE',
    });
  }

  // upravljanje sobama i prihvacenim aplikacijama
  // dobija sve prihvacene aplikacije (samo admin)
  async getAllPrihvaceneAplikacije() {
    return this.makeRequest('/prihvacene_aplikacije/');
  }

  // izbacuje studenta iz sobe (samo admin)
  async evictStudent(userId, reason) {
    return this.makeRequest('/prihvacene_aplikacije/evict', {
      method: 'POST',
      body: JSON.stringify({
        user_id: userId,
        reason: reason
      }),
    });
  }

  // dobija sva placanja sa opcionalnim filterom statusa (samo admin)
  async getAllPayments(status = null) {
    const url = status ? `/payments/?status=${status}` : '/payments/';
    return this.makeRequest(url);
  }

  // kreira novo placanje (samo admin)
  async createPayment(paymentData) {
    return this.makeRequest('/payments/', {
      method: 'POST',
      body: JSON.stringify(paymentData),
    });
  }

  // pretrazuje placanja po indeksu studenta (samo admin)
  async searchPaymentsByIndex(indexPattern, status = null) {
    const url = status 
      ? `/payments/search?index=${indexPattern}&status=${status}`
      : `/payments/search?index=${indexPattern}`;
    return this.makeRequest(url);
  }

  // dobija placanja po sobi (samo admin)
  async getPaymentsByRoom(sobaId) {
    return this.makeRequest(`/payments/room/${sobaId}`);
  }

  // dobija placanja po korisniku (samo admin)
  async getPaymentsByUser(userId) {
    return this.makeRequest(`/payments/user/${userId}`);
  }

  // oznacava placanje kao placeno (samo admin)
  async markPaymentAsPaid(paymentId, paidAt = null) {
    return this.makeRequest(`/payments/${paymentId}/mark-paid`, {
      method: 'PATCH',
      body: JSON.stringify(paidAt ? { paid_at: paidAt } : {}),
    });
  }

  // oznacava placanje kao neplaceno (samo admin)
  async markPaymentAsUnpaid(paymentId) {
    return this.makeRequest(`/payments/${paymentId}/mark-unpaid`, {
      method: 'PATCH',
    });
  }

  // endpoint-i specificni za korisnika

  // dobija moje prihvacene aplikacije (trenutni korisnik)
  async getMyPrihvaceneAplikacije() {
    return this.makeRequest('/prihvacene_aplikacije/my');
  }

  // odjavljuje se iz sobe (trenutni korisnik)
  async checkoutFromRoom() {
    return this.makeRequest('/prihvacene_aplikacije/checkout', {
      method: 'POST',
    });
  }

  // dobija placanja po ID-u aplikacije
  async getPaymentsByAplikacija(aplikacijaId) {
    return this.makeRequest(`/payments/aplikacija/${aplikacijaId}`);
  }
}

export const stDomService = new StDomService();
