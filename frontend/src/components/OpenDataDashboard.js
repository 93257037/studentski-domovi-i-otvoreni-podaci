import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { openDataService } from '../services/openDataService';
import './OpenDataDashboard.css';

// Open Data Dashboard - showcases all 6 open data functionalities
const OpenDataDashboard = () => {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState('statistics');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // Tab 1: Statistics
  const [statistics, setStatistics] = useState(null);

  // Tab 2: Room Search
  const [roomSearchFilters, setRoomSearchFilters] = useState({
    dormId: '',
    minCapacity: '',
    maxCapacity: '',
    amenities: [],
    onlyAvailable: true,
    limit: 20
  });
  const [roomResults, setRoomResults] = useState([]);
  const [availableAmenities, setAvailableAmenities] = useState([]);

  // Tab 3: Dorm Comparison
  const [dormList, setDormList] = useState([]);
  const [selectedDorms, setSelectedDorms] = useState([]);
  const [comparisonData, setComparisonData] = useState(null);

  // Tab 4: Application Trends
  const [trendsData, setTrendsData] = useState(null);

  // Tab 5: Occupancy Heatmap
  const [heatmapData, setHeatmapData] = useState(null);

  // Load initial data
  useEffect(() => {
    loadInitialData();
  }, []);

  const loadInitialData = async () => {
    try {
      // Load dorm list for comparison
      const dormsResponse = await openDataService.getDormList();
      setDormList(dormsResponse.dorms || []);

      // Load available amenities
      const amenitiesResponse = await openDataService.getAvailableAmenities();
      setAvailableAmenities(amenitiesResponse.amenities || []);
    } catch (err) {
      console.error('Error loading initial data:', err);
    }
  };

  // Load data when tab changes
  useEffect(() => {
    loadTabData();
  }, [activeTab]);

  const loadTabData = async () => {
    setLoading(true);
    setError('');

    try {
      switch (activeTab) {
        case 'statistics':
          if (!statistics) {
            const data = await openDataService.getPublicStatistics();
            setStatistics(data.statistics);
          }
          break;
        case 'trends':
          if (!trendsData) {
            const data = await openDataService.getApplicationTrends();
            setTrendsData(data.trends);
          }
          break;
        case 'heatmap':
          if (!heatmapData) {
            const data = await openDataService.getOccupancyHeatmap();
            setHeatmapData(data.heatmap);
          }
          break;
        default:
          break;
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  // Room Search handlers
  const handleRoomSearch = async () => {
    setLoading(true);
    setError('');
    try {
      const filters = {
        ...roomSearchFilters,
        minCapacity: roomSearchFilters.minCapacity ? parseInt(roomSearchFilters.minCapacity) : undefined,
        maxCapacity: roomSearchFilters.maxCapacity ? parseInt(roomSearchFilters.maxCapacity) : undefined,
      };
      const data = await openDataService.searchAvailableRooms(filters);
      setRoomResults(data.rooms || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const toggleAmenity = (amenity) => {
    setRoomSearchFilters(prev => ({
      ...prev,
      amenities: prev.amenities.includes(amenity)
        ? prev.amenities.filter(a => a !== amenity)
        : [...prev.amenities, amenity]
    }));
  };

  // Dorm Comparison handlers
  const toggleDormSelection = (dormId) => {
    setSelectedDorms(prev => {
      if (prev.includes(dormId)) {
        return prev.filter(id => id !== dormId);
      } else if (prev.length < 10) {
        return [...prev, dormId];
      } else {
        alert('Maximum 10 dorms can be compared at once');
        return prev;
      }
    });
  };

  const handleCompareDorms = async () => {
    if (selectedDorms.length === 0) {
      alert('Please select at least one dorm to compare');
      return;
    }

    setLoading(true);
    setError('');
    try {
      const data = await openDataService.compareDorms(selectedDorms);
      setComparisonData(data.comparison);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  // Export handlers
  const handleExport = async (dataset, format) => {
    try {
      await openDataService.downloadExport(dataset, format);
    } catch (err) {
      alert('Export failed: ' + err.message);
    }
  };

  // Render functions for each tab
  const renderStatistics = () => {
    if (!statistics) return null;

    return (
      <div className="statistics-container">
        <div className="stats-overview">
          <h2>Pregled Sistema</h2>
          <div className="stats-grid">
            <div className="stat-card">
              <h3>{statistics.total_dorms}</h3>
              <p>Ukupno Domova</p>
            </div>
            <div className="stat-card">
              <h3>{statistics.total_rooms}</h3>
              <p>Ukupno Soba</p>
            </div>
            <div className="stat-card">
              <h3>{statistics.total_capacity}</h3>
              <p>Ukupan Kapacitet</p>
            </div>
            <div className="stat-card">
              <h3>{statistics.total_occupied}</h3>
              <p>Popunjenih Mesta</p>
            </div>
            <div className="stat-card highlight">
              <h3>{statistics.occupancy_rate.toFixed(1)}%</h3>
              <p>Stopa Popunjenosti</p>
            </div>
          </div>
        </div>

        <div className="stats-section">
          <h2>Statistika Prijava</h2>
          <div className="stats-grid">
            <div className="stat-card">
              <h3>{statistics.application_statistics.total_applications}</h3>
              <p>Ukupno Prijava</p>
            </div>
            <div className="stat-card">
              <h3>{statistics.application_statistics.active_applications}</h3>
              <p>Aktivnih Prijava</p>
            </div>
            <div className="stat-card">
              <h3>{statistics.application_statistics.accepted_applications}</h3>
              <p>Prihvaćeno</p>
            </div>
            <div className="stat-card">
              <h3>{statistics.application_statistics.acceptance_rate.toFixed(1)}%</h3>
              <p>Stopa Prihvatanja</p>
            </div>
            <div className="stat-card">
              <h3>{statistics.application_statistics.average_grade_of_accepted.toFixed(2)}</h3>
              <p>Prosečan Prosek (Prihvaćeni)</p>
            </div>
          </div>
        </div>

        <div className="stats-section">
          <h2>Statistika po Domovima</h2>
          <div className="table-container">
            <table className="dorm-stats-table">
              <thead>
                <tr>
                  <th>Ime Doma</th>
                  <th>Adresa</th>
                  <th>Sobe</th>
                  <th>Kapacitet</th>
                  <th>Popunjeno</th>
                  <th>Dostupno</th>
                  <th>Popunjenost</th>
                </tr>
              </thead>
              <tbody>
                {statistics.dorm_statistics.map(dorm => (
                  <tr key={dorm.dorm_id}>
                    <td>{dorm.dorm_name}</td>
                    <td>{dorm.address}</td>
                    <td>{dorm.total_rooms}</td>
                    <td>{dorm.total_capacity}</td>
                    <td>{dorm.occupied_spots}</td>
                    <td>{dorm.available_spots}</td>
                    <td>
                      <span className={`occupancy-badge occupancy-${getOccupancyLevel(dorm.occupancy_rate)}`}>
                        {dorm.occupancy_rate.toFixed(1)}%
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        <div className="export-section">
          <h3>Izvoz Podataka</h3>
          <div className="export-buttons">
            <button onClick={() => handleExport('statistics', 'json')} className="btn btn-secondary">
              Izvezi kao JSON
            </button>
            <button onClick={() => handleExport('statistics', 'csv')} className="btn btn-secondary">
              Izvezi kao CSV
            </button>
          </div>
        </div>
      </div>
    );
  };

  const renderRoomSearch = () => {
    return (
      <div className="room-search-container">
        <h2>Pretraga Dostupnosti Soba</h2>
        
        <div className="search-filters">
          <div className="filter-group">
            <label>Dom</label>
            <select
              value={roomSearchFilters.dormId}
              onChange={(e) => setRoomSearchFilters(prev => ({ ...prev, dormId: e.target.value }))}
            >
              <option value="">Svi Domovi</option>
              {dormList.map(dorm => (
                <option key={dorm.id} value={dorm.id}>{dorm.name}</option>
              ))}
            </select>
          </div>

          <div className="filter-group">
            <label>Min. Kapacitet</label>
            <input
              type="number"
              min="1"
              value={roomSearchFilters.minCapacity}
              onChange={(e) => setRoomSearchFilters(prev => ({ ...prev, minCapacity: e.target.value }))}
              placeholder="Min. kreveta"
            />
          </div>

          <div className="filter-group">
            <label>Maks. Kapacitet</label>
            <input
              type="number"
              min="1"
              value={roomSearchFilters.maxCapacity}
              onChange={(e) => setRoomSearchFilters(prev => ({ ...prev, maxCapacity: e.target.value }))}
              placeholder="Maks. kreveta"
            />
          </div>

          <div className="filter-group">
            <label>
              <input
                type="checkbox"
                checked={roomSearchFilters.onlyAvailable}
                onChange={(e) => setRoomSearchFilters(prev => ({ ...prev, onlyAvailable: e.target.checked }))}
              />
              Samo Dostupne Sobe
            </label>
          </div>
        </div>

        <div className="amenities-filter">
          <label>Potrebne Pogodnosti:</label>
          <div className="amenities-checkboxes">
            {availableAmenities.map(amenity => (
              <label key={amenity} className="amenity-checkbox">
                <input
                  type="checkbox"
                  checked={roomSearchFilters.amenities.includes(amenity)}
                  onChange={() => toggleAmenity(amenity)}
                />
                {amenity}
              </label>
            ))}
          </div>
        </div>

        <button onClick={handleRoomSearch} className="btn btn-primary" disabled={loading}>
          {loading ? 'Pretraga...' : 'Pretraži Sobe'}
        </button>

        {roomResults.length > 0 && (
          <div className="results-section">
            <h3>Rezultati ({roomResults.length} soba)</h3>
            <div className="table-container">
              <table className="room-results-table">
                <thead>
                  <tr>
                    <th>Dom</th>
                    <th>Adresa</th>
                    <th>Kapacitet</th>
                    <th>Popunjeno</th>
                    <th>Dostupno</th>
                    <th>Pogodnosti</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  {roomResults.map(room => (
                    <tr key={room.room_id}>
                      <td>{room.dorm_name}</td>
                      <td>{room.dorm_address}</td>
                      <td>{room.capacity}</td>
                      <td>{room.occupied}</td>
                      <td>{room.available_spots}</td>
                      <td>{room.amenities.join(', ') || 'Bez pogodnosti'}</td>
                      <td>
                        <span className={`status-badge ${room.is_available ? 'available' : 'full'}`}>
                          {room.is_available ? 'Dostupno' : 'Popunjeno'}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    );
  };

  const renderDormComparison = () => {
    return (
      <div className="dorm-comparison-container">
        <h2>Alatka za Poređenje Domova</h2>
        
        <div className="dorm-selection">
          <h3>Izaberite Domove za Poređenje (maks. 10)</h3>
          <p className="selection-count">Izabrano: {selectedDorms.length}</p>
          <div className="dorm-checkboxes">
            {dormList.map(dorm => (
              <label key={dorm.id} className="dorm-checkbox">
                <input
                  type="checkbox"
                  checked={selectedDorms.includes(dorm.id)}
                  onChange={() => toggleDormSelection(dorm.id)}
                  disabled={!selectedDorms.includes(dorm.id) && selectedDorms.length >= 10}
                />
                <span>{dorm.name}</span>
                <small>{dorm.address}</small>
              </label>
            ))}
          </div>
          <button onClick={handleCompareDorms} className="btn btn-primary" disabled={loading || selectedDorms.length === 0}>
            {loading ? 'Poređenje...' : 'Uporedi Izabrane Domove'}
          </button>
        </div>

        {comparisonData && (
          <div className="comparison-results">
            <h3>Rezultati Poređenja</h3>
            <div className="comparison-grid">
              {comparisonData.dorms.map(dorm => (
                <div key={dorm.dorm_id} className="comparison-card">
                  <h4>{dorm.dorm_name}</h4>
                  <p className="address">{dorm.address}</p>
                  
                  <div className="contact-info">
                    <p><strong>Telefon:</strong> {dorm.contact_info.phone}</p>
                    <p><strong>Email:</strong> {dorm.contact_info.email}</p>
                  </div>

                  <div className="capacity-info">
                    <h5>Kapacitet</h5>
                    <p>Sobe: {dorm.capacity.total_rooms}</p>
                    <p>Ukupan Kapacitet: {dorm.capacity.total_capacity}</p>
                    <p>Popunjeno: {dorm.capacity.occupied_spots}</p>
                    <p>Dostupno: {dorm.capacity.available_spots}</p>
                    <p className="occupancy">
                      Popunjenost: <strong>{dorm.capacity.occupancy_rate.toFixed(1)}%</strong>
                    </p>
                  </div>

                  <div className="application-info">
                    <h5>Prijave</h5>
                    <p>Ukupno: {dorm.application_metrics.total_applications}</p>
                    <p>Prihvaćeno: {dorm.application_metrics.accepted_applications}</p>
                    <p>Stopa Prihvatanja: {dorm.application_metrics.acceptance_rate.toFixed(1)}%</p>
                    <p>Prosečan Prosek: {dorm.application_metrics.average_grade.toFixed(2)}</p>
                  </div>

                  <div className="amenities-info">
                    <h5>Pogodnosti</h5>
                    {Object.entries(dorm.amenities_offered).map(([amenity, count]) => (
                      <p key={amenity}>{amenity}: {count} soba</p>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    );
  };

  const renderTrends = () => {
    if (!trendsData) return null;

    return (
      <div className="trends-container">
        <h2>Analiza Trendova Prijava</h2>

        <div className="trends-overview">
          <h3>Ukupna Statistika</h3>
          <div className="stats-grid">
            <div className="stat-card">
              <h3>{trendsData.overall_metrics?.total_years || 0}</h3>
              <p>Praćenih Godina</p>
            </div>
            <div className="stat-card">
              <h3>{trendsData.overall_metrics?.average_applications_per_year || 0}</h3>
              <p>Prosečno Prijava/Godinu</p>
            </div>
            <div className="stat-card">
              <h3 className={`trend-${trendsData.overall_metrics?.trend_direction || 'stable'}`}>
                {trendsData.overall_metrics?.trend_direction === 'increasing' ? 'Raste' : 
                 trendsData.overall_metrics?.trend_direction === 'decreasing' ? 'Opada' :
                 trendsData.overall_metrics?.trend_direction === 'stable' ? 'Stabilno' : 'N/A'}
              </h3>
              <p>Pravac Trenda</p>
            </div>
          </div>
        </div>

        <div className="yearly-trends">
          <h3>Godišnji Trendovi</h3>
          <div className="table-container">
            <table className="trends-table">
              <thead>
                <tr>
                  <th>Školska Godina</th>
                  <th>Prijave</th>
                  <th>Prihvaćeno</th>
                  <th>Stopa Prihvatanja</th>
                  <th>Prosečan Prosek</th>
                  <th>Min. Prosek</th>
                  <th>Maks. Prosek</th>
                </tr>
              </thead>
              <tbody>
                {trendsData.yearly_trends && trendsData.yearly_trends.length > 0 ? (
                  trendsData.yearly_trends.map(year => (
                    <tr key={year.academic_year}>
                      <td>{year.academic_year}</td>
                      <td>{year.total_applications}</td>
                      <td>{year.accepted_applications}</td>
                      <td>{year.acceptance_rate?.toFixed(1) || 0}%</td>
                      <td>{year.average_grade?.toFixed(2) || 0}</td>
                      <td>{year.min_grade}</td>
                      <td>{year.max_grade}</td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan="7" style={{textAlign: 'center', padding: '20px'}}>
                      Nema dostupnih podataka o trendovima
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>

        <div className="dorm-trends">
          <h3>Trendovi Prijava po Domovima</h3>
          <div className="table-container">
            <table className="trends-table">
              <thead>
                <tr>
                  <th>Dom</th>
                  <th>Ukupno Prijava</th>
                  <th>Prihvaćeno</th>
                  <th>Stopa Prihvatanja</th>
                </tr>
              </thead>
              <tbody>
                {trendsData.dorm_trends && trendsData.dorm_trends.length > 0 ? (
                  trendsData.dorm_trends
                    .sort((a, b) => b.total_applications - a.total_applications)
                    .map(dorm => (
                      <tr key={dorm.dorm_id}>
                        <td>{dorm.dorm_name}</td>
                        <td>{dorm.total_applications}</td>
                        <td>{dorm.accepted_applications}</td>
                        <td>{dorm.acceptance_rate?.toFixed(1) || 0}%</td>
                      </tr>
                    ))
                ) : (
                  <tr>
                    <td colSpan="4" style={{textAlign: 'center', padding: '20px'}}>
                      Nema dostupnih podataka o trendovima domova
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    );
  };

  const renderHeatmap = () => {
    if (!heatmapData) return null;

    return (
      <div className="heatmap-container">
        <h2>Mapa Popunjenosti u Realnom Vremenu</h2>

        <div className="heatmap-summary">
          <h3>Pregled</h3>
          <div className="stats-grid">
            <div className="stat-card">
              <h3>{heatmapData.summary.average_occupancy.toFixed(1)}%</h3>
              <p>Prosečna Popunjenost</p>
            </div>
            <div className="stat-card">
              <h3>{heatmapData.summary.highest_occupancy.toFixed(1)}%</h3>
              <p>Najveća Popunjenost</p>
            </div>
            <div className="stat-card">
              <h3>{heatmapData.summary.lowest_occupancy.toFixed(1)}%</h3>
              <p>Najniža Popunjenost</p>
            </div>
            <div className="stat-card">
              <h3>{heatmapData.summary.full_dorms}</h3>
              <p>Popunjenih Domova</p>
            </div>
            <div className="stat-card">
              <h3>{heatmapData.summary.empty_dorms}</h3>
              <p>Praznih Domova</p>
            </div>
          </div>
        </div>

        <div className="heatmap-grid">
          {heatmapData.dorms
            .sort((a, b) => b.occupancy_rate - a.occupancy_rate)
            .map(dorm => (
              <div key={dorm.dorm_id} className={`heatmap-card occupancy-status-${dorm.status}`}>
                <h4>{dorm.dorm_name}</h4>
                <p className="address">{dorm.address}</p>
                <div className="occupancy-meter">
                  <div 
                    className="occupancy-fill"
                    style={{ width: `${dorm.occupancy_rate}%` }}
                  />
                </div>
                <p className="occupancy-rate">{dorm.occupancy_rate.toFixed(1)}%</p>
                <p className="capacity-info">
                  {dorm.occupied_spots} / {dorm.total_capacity} popunjeno
                </p>
                <p className="available-info">
                  {dorm.available_spots} mesta dostupno
                </p>
                <span className={`status-badge status-${dorm.status}`}>
                  {dorm.status === 'available' ? 'Dostupno' :
                   dorm.status === 'limited' ? 'Ograničeno' :
                   dorm.status === 'full' ? 'Popunjeno' :
                   dorm.status === 'empty' ? 'Prazno' : dorm.status}
                </span>
              </div>
            ))}
        </div>
      </div>
    );
  };

  const getOccupancyLevel = (rate) => {
    if (rate >= 80) return 'high';
    if (rate >= 50) return 'medium';
    return 'low';
  };

  return (
    <div className="open-data-dashboard">
      <header className="dashboard-header">
        <div className="header-content">
          <div className="header-text">
            <h1>📊 Otvoreni Podaci</h1>
            <p>Javni pristup podacima i statistici studentskih domova</p>
          </div>
          <button onClick={() => navigate('/')} className="back-button">
            ← Nazad
          </button>
        </div>
      </header>

      <nav className="dashboard-tabs">
        <button
          className={activeTab === 'statistics' ? 'active' : ''}
          onClick={() => setActiveTab('statistics')}
        >
          📈 Statistika
        </button>
        <button
          className={activeTab === 'rooms' ? 'active' : ''}
          onClick={() => setActiveTab('rooms')}
        >
          🔍 Pretraga Soba
        </button>
        <button
          className={activeTab === 'comparison' ? 'active' : ''}
          onClick={() => setActiveTab('comparison')}
        >
          ⚖️ Uporedi Domove
        </button>
        <button
          className={activeTab === 'trends' ? 'active' : ''}
          onClick={() => setActiveTab('trends')}
        >
          📊 Trendovi
        </button>
        <button
          className={activeTab === 'heatmap' ? 'active' : ''}
          onClick={() => setActiveTab('heatmap')}
        >
          🗺️ Mapa Popunjenosti
        </button>
      </nav>

      <div className="dashboard-content">
        {loading && <div className="loading-spinner">Učitavanje...</div>}
        {error && <div className="error-message">{error}</div>}

        {!loading && !error && (
          <>
            {activeTab === 'statistics' && renderStatistics()}
            {activeTab === 'rooms' && renderRoomSearch()}
            {activeTab === 'comparison' && renderDormComparison()}
            {activeTab === 'trends' && renderTrends()}
            {activeTab === 'heatmap' && renderHeatmap()}
          </>
        )}
      </div>
    </div>
  );
};

export default OpenDataDashboard;

