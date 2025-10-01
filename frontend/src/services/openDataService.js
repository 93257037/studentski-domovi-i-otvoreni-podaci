const API_BASE_URL = process.env.REACT_APP_OPEN_DATA_SERVICE_URL || 'http://localhost:8082/api/v1';

class OpenDataService {
  constructor() {
    this.baseURL = API_BASE_URL;
  }

  async makeRequest(endpoint, options = {}) {
    const config = {
      headers: {
        'Content-Type': 'application/json',
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

  // Search student dormitories by name (ime)
  async searchStDomsByIme(imePattern) {
    const params = new URLSearchParams({ ime: imePattern });
    return this.makeRequest(`/st-doms/search-by-ime?${params}`);
  }

  // Search student dormitories by address
  async searchStDomsByAddress(addressPattern) {
    const params = new URLSearchParams({ address: addressPattern });
    return this.makeRequest(`/st-doms/search-by-address?${params}`);
  }

  // Get all student dormitories
  async getAllStDoms() {
    return this.makeRequest('/st-doms');
  }

  // Get all rooms
  async getAllRooms() {
    return this.makeRequest('/rooms');
  }

  // Filter rooms by luxury amenities
  async filterRoomsByLuksuz(luksuzi) {
    const params = new URLSearchParams({ luksuzi: luksuzi.join(',') });
    return this.makeRequest(`/rooms/filter-by-luksuz?${params}`);
  }

  // Filter rooms by bed capacity
  async filterRoomsByKrevetnost(exact, min, max) {
    const params = new URLSearchParams();
    if (exact) params.append('exact', exact);
    if (min) params.append('min', min);
    if (max) params.append('max', max);
    return this.makeRequest(`/rooms/filter-by-krevetnost?${params}`);
  }

  // Advanced room filtering
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

  // Get statistics
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
}

export const openDataService = new OpenDataService();
