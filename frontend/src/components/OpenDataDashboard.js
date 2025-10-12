import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { openDataService } from '../services/openDataService';
import './OpenDataDashboard.css';

// Open Data Dashboard - showcases all 6 open data functionalities
const OpenDataDashboard = () => {
  const navigate = useNavigate();
  const { token } = useAuth();
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

  // Dorm list for filters
  const [dormList, setDormList] = useState([]);

  // Application Trends
  const [trendsData, setTrendsData] = useState(null);

  // Occupancy Heatmap
  const [heatmapData, setHeatmapData] = useState(null);

  // Academic Year Applications
  const [academicYear, setAcademicYear] = useState('');
  const [yearApplications, setYearApplications] = useState([]);
  const [hasSearchedYear, setHasSearchedYear] = useState(false);

  // Load initial data
  useEffect(() => {
    loadInitialData();
  }, []);

  const loadInitialData = async () => {
    try {
      // Load dorm list for filters
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
            <button onClick={() => handleExport('dorm-statistics', 'json')} className="btn btn-secondary">
              Izvezi Statistiku Domova (JSON)
            </button>
            <button onClick={() => handleExport('dorm-statistics', 'csv')} className="btn btn-secondary">
              Izvezi Statistiku Domova (CSV)
            </button>
            <button onClick={() => handleExport('dorms', 'json')} className="btn btn-secondary">
              Izvezi Listu Domova (JSON)
            </button>
            <button onClick={() => handleExport('dorms', 'csv')} className="btn btn-secondary">
              Izvezi Listu Domova (CSV)
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

        <div className="export-section" style={{ marginTop: '20px' }}>
          <h3>Izvoz Svih Soba</h3>
          <div className="export-buttons">
            <button onClick={() => handleExport('rooms', 'json')} className="btn btn-secondary">
              Izvezi kao JSON
            </button>
            <button onClick={() => handleExport('rooms', 'csv')} className="btn btn-secondary">
              Izvezi kao CSV
            </button>
          </div>
        </div>

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
                    <tr 
                      key={room.room_id} 
                      onClick={() => navigate(`/room/${room.room_id}`)}
                      style={{ cursor: 'pointer' }}
                      className="clickable-row"
                    >
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

  const renderTrends = () => {
    if (!trendsData) return null;

    return (
      <div className="trends-container">
        <h2>Analiza Kretanja Prijava</h2>

        <div className="yearly-trends">
          <h3>Godišnja Kretanja</h3>
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
                      Nema dostupnih podataka o kretanjima
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          <div className="export-section">
            <h3>Izvoz Podataka</h3>
            <div className="export-buttons">
              <button onClick={() => handleExport('godisnja-kretanja', 'json')} className="btn btn-secondary">
                Izvezi Godišnja Kretanja (JSON)
              </button>
              <button onClick={() => handleExport('godisnja-kretanja', 'csv')} className="btn btn-secondary">
                Izvezi Godišnja Kretanja (CSV)
              </button>
            </div>
          </div>
        </div>

        <div className="dorm-trends">
          <h3>Kretanja Prijava po Domovima</h3>
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
                      Nema dostupnih podataka o kretanjima domova
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          <div className="export-section">
            <h3>Izvoz Podataka</h3>
            <div className="export-buttons">
              <button onClick={() => handleExport('application-list', 'json')} className="btn btn-secondary">
                Izvezi Listu Prijava (JSON)
              </button>
              <button onClick={() => handleExport('application-list', 'csv')} className="btn btn-secondary">
                Izvezi Listu Prijava (CSV)
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  };

  const renderHeatmap = () => {
    if (!heatmapData) return null;

    return (
      <div className="heatmap-container">
        <h2>Karta Popunjenosti u Realnom Vremenu</h2>

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

        <div className="export-section">
          <h3>Izvoz Podataka</h3>
          <div className="export-buttons">
            <button onClick={() => handleExport('room-types', 'json')} className="btn btn-secondary">
              Izvezi Tipove Soba (JSON)
            </button>
            <button onClick={() => handleExport('room-types', 'csv')} className="btn btn-secondary">
              Izvezi Tipove Soba (CSV)
            </button>
          </div>
        </div>
      </div>
    );
  };

  const getOccupancyLevel = (rate) => {
    if (rate >= 80) return 'high';
    if (rate >= 50) return 'medium';
    return 'low';
  };

  const renderAcademicYearApplications = () => {
    const handleYearSearch = async (e) => {
      e.preventDefault();
      if (!academicYear.trim()) {
        setError('Molimo unesite akademsku godinu');
        return;
      }

      // Validate academic year format (should be like 2024/2025)
      const academicYearPattern = /^\d{4}\/\d{4}$/;
      if (!academicYearPattern.test(academicYear.trim())) {
        setError('Akademska godina mora biti u formatu YYYY/YYYY (npr. 2024/2025)');
        return;
      }

      setLoading(true);
      setError('');
      setHasSearchedYear(true);
      
      try {
        const response = await openDataService.getPrihvaceneAplikacijeForAcademicYear(academicYear.trim());
        setYearApplications(response.data?.prihvacene_aplikacije || []);
      } catch (error) {
        setError('Greška pri dohvaćanju podataka: ' + error.message);
        setYearApplications([]);
      } finally {
        setLoading(false);
      }
    };

    const clearYearSearch = () => {
      setAcademicYear('');
      setYearApplications([]);
      setError('');
      setHasSearchedYear(false);
    };

    const formatDate = (dateString) => {
      if (!dateString) return 'N/A';
      return new Date(dateString).toLocaleDateString('hr-HR', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    };

    return (
      <div className="tab-content">
        <h2>Prihvaćene aplikacije po akademskoj godini</h2>
        <p className="description">Pretražite prihvaćene aplikacije za određenu akademsku godinu</p>

        <div className="search-card">
          <form onSubmit={handleYearSearch} className="search-form">
            <div className="search-input-group">
              <label htmlFor="academicYear">Akademska godina:</label>
              <input
                id="academicYear"
                type="text"
                value={academicYear}
                onChange={(e) => setAcademicYear(e.target.value)}
                placeholder="Unesite akademsku godinu (npr. 2024/2025)"
                className="search-input"
                disabled={loading}
              />
              <div className="button-group">
                <button type="submit" className="btn btn-primary" disabled={loading}>
                  {loading ? 'Pretraživanje...' : 'Pretraži'}
                </button>
                {academicYear && (
                  <button type="button" onClick={clearYearSearch} className="btn btn-secondary" disabled={loading}>
                    Obriši
                  </button>
                )}
              </div>
            </div>
          </form>
          
          <div className="format-hint">
            <p><strong>Format:</strong> Akademska godina mora biti u formatu YYYY/YYYY (npr. 2024/2025)</p>
          </div>
        </div>

        <div className="export-section" style={{ marginTop: '20px' }}>
          <h3>Izvoz Svih Prihvaćenih Aplikacija</h3>
          <div className="export-buttons">
            <button onClick={() => handleExport('accepted-applications', 'json')} className="btn btn-secondary">
              Izvezi kao JSON
            </button>
            <button onClick={() => handleExport('accepted-applications', 'csv')} className="btn btn-secondary">
              Izvezi kao CSV
            </button>
          </div>
        </div>

        {hasSearchedYear && !loading && (
          <div className="results-section">
            <div className="results-header">
              <h3>Rezultati pretrage</h3>
              <div className="results-info">
                <span className="academic-year-label">Akademska godina: <strong>{academicYear}</strong></span>
                <span className="count-label">Ukupno aplikacija: <strong>{yearApplications.length}</strong></span>
              </div>
            </div>

            {yearApplications.length > 0 ? (
              <div className="applications-grid">
                {yearApplications.map((app, index) => (
                  <div key={app.id || index} className="application-card">
                    <div className="application-header">
                      <h4>Aplikacija #{index + 1}</h4>
                      <span className="application-id">ID: {app.id}</span>
                    </div>
                    
                    <div className="application-details">
                      <div className="detail-row">
                        <span className="label">Korisnik ID:</span>
                        <span className="value">{app.user_id}</span>
                      </div>
                      <div className="detail-row">
                        <span className="label">Broj indeksa:</span>
                        <span className="value">{app.broj_indexa}</span>
                      </div>
                      <div className="detail-row">
                        <span className="label">Prosek:</span>
                        <span className="value prosek">{app.prosek}</span>
                      </div>
                      <div className="detail-row">
                        <span className="label">Soba ID:</span>
                        <span className="value">{app.soba_id}</span>
                      </div>
                      <div className="detail-row">
                        <span className="label">Akademska godina:</span>
                        <span className="value">{app.academic_year}</span>
                      </div>
                      <div className="detail-row">
                        <span className="label">Kreirana:</span>
                        <span className="value">{app.created_at ? formatDate(app.created_at) : 'N/A'}</span>
                      </div>
                      <div className="detail-row">
                        <span className="label">Ažurirana:</span>
                        <span className="value">{app.updated_at ? formatDate(app.updated_at) : 'N/A'}</span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="no-results">
                <h3>Nema rezultata</h3>
                <p>Nisu pronađene prihvaćene aplikacije za akademsku godinu "{academicYear}".</p>
              </div>
            )}
          </div>
        )}
      </div>
    );
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
          className={activeTab === 'trends' ? 'active' : ''}
          onClick={() => setActiveTab('trends')}
        >
          📊 Kretanja
        </button>
        <button
          className={activeTab === 'heatmap' ? 'active' : ''}
          onClick={() => setActiveTab('heatmap')}
        >
          🗺️ Karta Popunjenosti
        </button>
        <button
          className={activeTab === 'academic-year' ? 'active' : ''}
          onClick={() => setActiveTab('academic-year')}
        >
          📚 Prijave po godini
        </button>
      </nav>

      <div className="dashboard-content">
        {loading && <div className="loading-spinner">Učitavanje...</div>}
        {error && <div className="error-message">{error}</div>}

        {!loading && !error && (
          <>
            {activeTab === 'statistics' && renderStatistics()}
            {activeTab === 'rooms' && renderRoomSearch()}
            {activeTab === 'trends' && renderTrends()}
            {activeTab === 'heatmap' && renderHeatmap()}
            {activeTab === 'academic-year' && renderAcademicYearApplications()}
          </>
        )}
      </div>
    </div>
  );
};

export default OpenDataDashboard;

