// API Gateway URL - svi servisi su dostupni kroz jedan gateway
const API_BASE_URL = process.env.REACT_APP_API_GATEWAY_URL || 'http://localhost/api/v1';

// servis za komunikaciju sa open_data_service API-jem
// rukuje pretragom domova, soba i statistikama
class OpenDataService {
  constructor() {
    this.baseURL = API_BASE_URL;
  }

  // pravi HTTP zahtev bez automatskog dodavanja tokena (javni API)
  async makeRequest(endpoint, options = {}) {
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    const fullURL = `${this.baseURL}${endpoint}`;
    const response = await fetch(fullURL, config);
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  // pretrazuje studentske domove po imenu
  async searchStDomsByIme(imePattern) {
    const params = new URLSearchParams({ ime: imePattern });
    return this.makeRequest(`/st-doms/search-by-ime?${params}`);
  }

  // pretrazuje studentske domove po adresi
  async searchStDomsByAddress(addressPattern) {
    const params = new URLSearchParams({ address: addressPattern });
    return this.makeRequest(`/st-doms/search-by-address?${params}`);
  }

  // dobija sve studentske domove
  async getAllStDoms() {
    return this.makeRequest('/st-doms');
  }

  // dobija sve sobe
  async getAllRooms() {
    return this.makeRequest('/rooms');
  }

  // filtrira sobe po luksuznim sadrzajima
  async filterRoomsByLuksuz(luksuzi) {
    const params = new URLSearchParams({ luksuzi: luksuzi.join(',') });
    return this.makeRequest(`/rooms/filter-by-luksuz?${params}`);
  }

  // filtrira sobe po broju kreveta
  async filterRoomsByKrevetnost(exact, min, max) {
    const params = new URLSearchParams();
    if (exact) params.append('exact', exact);
    if (min) params.append('min', min);
    if (max) params.append('max', max);
    return this.makeRequest(`/rooms/filter-by-krevetnost?${params}`);
  }

  // napredna pretraga soba sa vise kriterijuma
  async advancedFilterRooms(luksuzi, stDomId, address, exact, min, max) {
    const params = new URLSearchParams();
    if (luksuzi && luksuzi.length > 0) params.append('luksuzi', luksuzi.join(','));
    if (stDomId) params.append('st_dom_id', stDomId);
    if (address) params.append('address', address);
    if (exact) params.append('exact', exact);
    if (min) params.append('min', min);
    if (max) params.append('max', max);
    return this.makeRequest(`/rooms/advanced-filter?${params}`);
  }

  // dobija statistike
  async getTopFullStDoms() {
    return this.makeRequest('/statistics/top-full-st-doms');
  }

  async getTopEmptyStDoms() {
    return this.makeRequest('/statistics/top-empty-st-doms');
  }

  async getStDomWithMostApplications() {
    return this.makeRequest('/statistics/st-dom-most-applications');
  }

  async getStDomWithHighestAverageProsek() {
    return this.makeRequest('/statistics/st-dom-highest-average-prosek');
  }

  // metode za komunikaciju izmedju servisa
  async getPrihvaceneAplikacijeForAcademicYear(academicYear, token) {
    const headers = {};
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
    
    const params = new URLSearchParams({ academic_year: academicYear });
    const endpoint = `/inter-service/prihvacene-aplikacije/academic-year?${params}`;
    
    return this.makeRequest(endpoint, {
      headers
    });
  }

  // dobija prihvacene aplikacije za odredjenu sobu
  async getPrihvaceneAplikacijeForRoom(roomId, token) {
    const headers = {};
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
    
    const endpoint = `/inter-service/prihvacene-aplikacije/room/${roomId}`;
    
    return this.makeRequest(endpoint, {
      headers
    });
  }

  // dobija sve prihvacene aplikacije
  async getPrihvaceneAplikacije(token) {
    const headers = {};
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
    
    const endpoint = `/inter-service/prihvacene-aplikacije`;
    
    return this.makeRequest(endpoint, {
      headers
    });
  }

  // dobija prihvacene aplikacije za odredjenog korisnika
  async getPrihvaceneAplikacijeForUser(userId, token) {
    const headers = {};
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
    
    const endpoint = `/inter-service/prihvacene-aplikacije/user/${userId}`;
    
    return this.makeRequest(endpoint, {
      headers
    });
  }
}

export const openDataService = new OpenDataService();
