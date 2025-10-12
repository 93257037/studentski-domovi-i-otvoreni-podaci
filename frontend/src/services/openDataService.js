// API Gateway URL - all services accessible through one gateway
const API_BASE_URL = process.env.REACT_APP_API_GATEWAY_URL || 'http://localhost/api/v1';

// Service for communication with open_data_service API
// Handles public data access, statistics, and analytics
class OpenDataService {
  constructor() {
    this.baseURL = API_BASE_URL;
  }

  // Makes HTTP request without automatically adding token (public API)
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

  // ====================
  // 1. Public Statistics Dashboard
  // ====================

  /**
   * Get comprehensive public statistics about all dorms
   * @returns {Promise} Statistics object with dorm data, occupancy, applications, etc.
   */
  async getPublicStatistics() {
    return this.makeRequest('/open-data/statistics');
  }

  // ====================
  // 2. Room Availability Search
  // ====================

  /**
   * Search for available rooms with advanced filters
   * @param {Object} filters - Filter options
   * @param {string} filters.dormId - Filter by specific dorm ID
   * @param {number} filters.minCapacity - Minimum bed capacity
   * @param {number} filters.maxCapacity - Maximum bed capacity
   * @param {Array<string>} filters.amenities - Required amenities
   * @param {boolean} filters.onlyAvailable - Only show available rooms
   * @param {number} filters.limit - Number of results (default: 50)
   * @param {number} filters.offset - Pagination offset
   * @returns {Promise} List of rooms with availability info
   */
  async searchAvailableRooms(filters = {}) {
    const params = new URLSearchParams();
    
    if (filters.dormId) params.append('dorm_id', filters.dormId);
    if (filters.minCapacity) params.append('min_capacity', filters.minCapacity);
    if (filters.maxCapacity) params.append('max_capacity', filters.maxCapacity);
    if (filters.amenities && filters.amenities.length > 0) {
      params.append('amenities', filters.amenities.join(','));
    }
    if (filters.onlyAvailable !== undefined) {
      params.append('only_available', filters.onlyAvailable.toString());
    }
    if (filters.limit) params.append('limit', filters.limit);
    if (filters.offset) params.append('offset', filters.offset);

    return this.makeRequest(`/open-data/rooms/search?${params}`);
  }

  // ====================
  // 3. Dorm Comparison Tool
  // ====================

  /**
   * Compare multiple dorms side-by-side
   * @param {Array<string>} dormIds - Array of dorm IDs to compare (max 10)
   * @returns {Promise} Comparison data for all dorms
   */
  async compareDorms(dormIds) {
    if (!dormIds || dormIds.length === 0) {
      throw new Error('At least one dorm ID is required');
    }
    if (dormIds.length > 10) {
      throw new Error('Maximum 10 dorms can be compared at once');
    }

    const params = new URLSearchParams();
    params.append('dorm_ids', dormIds.join(','));

    return this.makeRequest(`/open-data/dorms/compare?${params}`);
  }

  /**
   * Get a simple list of all dorms (for dropdown/selection)
   * @returns {Promise} List of dorms with basic info
   */
  async getDormList() {
    return this.makeRequest('/open-data/dorms/list');
  }

  // ====================
  // 4. Application Trends Analysis
  // ====================

  /**
   * Get historical trends of applications by academic year
   * @returns {Promise} Trends data with yearly and per-dorm statistics
   */
  async getApplicationTrends() {
    return this.makeRequest('/open-data/trends/applications');
  }

  // ====================
  // 5. Real-time Occupancy Heatmap
  // ====================

  /**
   * Get real-time occupancy data for visualization
   * @returns {Promise} Heatmap data with occupancy for each dorm
   */
  async getOccupancyHeatmap() {
    return this.makeRequest('/open-data/occupancy/heatmap');
  }

  // ====================
  // 6. Open Data Export (CSV/JSON)
  // ====================

  /**
   * Export data in CSV or JSON format
   * @param {string} dataset - Type of data: 'dorms', 'rooms', or 'statistics'
   * @param {string} format - Export format: 'json' or 'csv'
   * @returns {Promise} Exported data
   */
  async exportData(dataset, format = 'json') {
    const params = new URLSearchParams();
    params.append('dataset', dataset);
    params.append('format', format);

    if (format === 'csv') {
      // For CSV, we need to handle the response differently
      const fullURL = `${this.baseURL}/open-data/export?${params}`;
      const response = await fetch(fullURL);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      return response.text(); // Return CSV as text
    }

    return this.makeRequest(`/open-data/export?${params}`);
  }

  /**
   * Download exported data as a file
   * @param {string} dataset - Type of data: 'dorms', 'rooms', or 'statistics'
   * @param {string} format - Export format: 'json' or 'csv'
   */
  async downloadExport(dataset, format = 'json') {
    const data = await this.exportData(dataset, format);
    
    let blob;
    let filename;
    
    if (format === 'csv') {
      blob = new Blob([data], { type: 'text/csv' });
      filename = `${dataset}.csv`;
    } else {
      blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
      filename = `${dataset}.json`;
    }

    // Create download link
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
  }

  // ====================
  // Helper Endpoints
  // ====================

  /**
   * Get list of all available amenities
   * @returns {Promise} List of amenity names
   */
  async getAvailableAmenities() {
    return this.makeRequest('/open-data/amenities');
  }

  // ====================
  // Legacy/Backward Compatibility Methods
  // ====================
  // These methods map old API calls to new endpoints for backward compatibility

  /**
   * @deprecated Use getDormList() instead
   */
  async getAllStDoms() {
    return this.getDormList();
  }

  /**
   * @deprecated Use searchAvailableRooms() instead
   */
  async getAllRooms() {
    return this.searchAvailableRooms({ limit: 1000 });
  }

  /**
   * @deprecated Use searchAvailableRooms() with amenities filter instead
   */
  async filterRoomsByLuksuz(luksuzi) {
    return this.searchAvailableRooms({ amenities: luksuzi, limit: 1000 });
  }

  /**
   * @deprecated Use searchAvailableRooms() with capacity filters instead
   */
  async filterRoomsByKrevetnost(exact, min, max) {
    const filters = { limit: 1000 };
    if (exact) {
      filters.minCapacity = exact;
      filters.maxCapacity = exact;
    } else {
      if (min) filters.minCapacity = min;
      if (max) filters.maxCapacity = max;
    }
    return this.searchAvailableRooms(filters);
  }

  /**
   * @deprecated Use searchAvailableRooms() instead
   */
  async advancedFilterRooms(luksuzi, stDomId, address, exact, min, max) {
    const filters = { limit: 1000 };
    if (luksuzi && luksuzi.length > 0) filters.amenities = luksuzi;
    if (stDomId) filters.dormId = stDomId;
    if (exact) {
      filters.minCapacity = exact;
      filters.maxCapacity = exact;
    } else {
      if (min) filters.minCapacity = min;
      if (max) filters.maxCapacity = max;
    }
    return this.searchAvailableRooms(filters);
  }

  /**
   * @deprecated Use getPublicStatistics() and filter the results
   */
  async getTopFullStDoms() {
    const stats = await this.getPublicStatistics();
    return {
      st_doms: stats.statistics.dorm_statistics
        .sort((a, b) => b.occupancy_rate - a.occupancy_rate)
        .slice(0, 5)
    };
  }

  /**
   * @deprecated Use getPublicStatistics() and filter the results
   */
  async getTopEmptyStDoms() {
    const stats = await this.getPublicStatistics();
    return {
      st_doms: stats.statistics.dorm_statistics
        .sort((a, b) => a.occupancy_rate - b.occupancy_rate)
        .slice(0, 5)
    };
  }

  /**
   * @deprecated Use getApplicationTrends() instead
   */
  async getStDomWithMostApplications() {
    const trends = await this.getApplicationTrends();
    const sortedDorms = [...trends.trends.dorm_trends]
      .sort((a, b) => b.total_applications - a.total_applications);
    
    return {
      st_dom: sortedDorms.length > 0 ? {
        id: sortedDorms[0].dorm_id,
        name: sortedDorms[0].dorm_name,
        application_count: sortedDorms[0].total_applications
      } : null
    };
  }

  /**
   * @deprecated Use getPublicStatistics() and analyze application_statistics
   */
  async getStDomWithHighestAverageProsek() {
    const stats = await this.getPublicStatistics();
    return {
      average_prosek: stats.statistics.application_statistics.average_grade_of_accepted
    };
  }
}

export const openDataService = new OpenDataService();
