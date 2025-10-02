import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { openDataService } from '../services/openDataService';
import ApplyToRoomButton from './ApplyToRoomButton';
import './AdvancedRoomSearch.css';

const AdvancedRoomSearch = () => {
  const navigate = useNavigate();
  
  // Search form state
  const [searchParams, setSearchParams] = useState({
    luksuzi: [],
    stDomId: '',
    address: '',
    exact: '',
    min: '',
    max: ''
  });
  
  // Available luxury amenities
  const availableLuksuzi = [
    { value: 'klima', label: 'Klima' },
    { value: 'terasa', label: 'Terasa' },
    { value: 'sopstveno kupatilo', label: 'Sopstveno kupatilo' },
    { value: 'áram', label: 'Áram' },
    { value: 'ablak', label: 'Ablak' },
    { value: 'neisvrljan zid', label: 'Neisvrljan zid' }
  ];
  
  // State for student dormitories dropdown
  const [stDoms, setStDoms] = useState([]);
  const [stDomsLoading, setStDomsLoading] = useState(false);
  
  // Search results state
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [hasSearched, setHasSearched] = useState(false);

  // Load student dormitories for dropdown
  useEffect(() => {
    const loadStDoms = async () => {
      setStDomsLoading(true);
      try {
        const response = await openDataService.getAllStDoms();
        setStDoms(response.data || []);
      } catch (error) {
        console.error('Error loading student dormitories:', error);
      } finally {
        setStDomsLoading(false);
      }
    };
    
    loadStDoms();
  }, []);

  // Handle luxury amenities checkbox changes
  const handleLuksuziChange = (luksuzValue) => {
    setSearchParams(prev => ({
      ...prev,
      luksuzi: prev.luksuzi.includes(luksuzValue)
        ? prev.luksuzi.filter(l => l !== luksuzValue)
        : [...prev.luksuzi, luksuzValue]
    }));
  };

  // Handle other form field changes
  const handleInputChange = (field, value) => {
    setSearchParams(prev => ({
      ...prev,
      [field]: value
    }));
  };

  // Validate form inputs
  const validateForm = () => {
    const { exact, min, max } = searchParams;
    
    // Check if exact is provided with min/max
    if (exact && (min || max)) {
      return 'Ne možete koristiti tačan broj kreveta sa min/max vrijednostima istovremeno.';
    }
    
    // Check if min > max
    if (min && max && parseInt(min) > parseInt(max)) {
      return 'Minimalni broj kreveta ne može biti veći od maksimalnog.';
    }
    
    // Check for negative values
    if ((exact && parseInt(exact) < 1) || (min && parseInt(min) < 1) || (max && parseInt(max) < 1)) {
      return 'Broj kreveta mora biti pozitivan broj.';
    }
    
    return null;
  };

  // Handle form submission
  const handleSearch = async (e) => {
    e.preventDefault();
    
    const validationError = validateForm();
    if (validationError) {
      setError(validationError);
      return;
    }
    
    setLoading(true);
    setError('');
    setHasSearched(true);
    
    try {
      const { luksuzi, stDomId, address, exact, min, max } = searchParams;
      
      const response = await openDataService.advancedFilterRooms(
        luksuzi.length > 0 ? luksuzi : null,
        stDomId || null,
        address || null,
        exact ? parseInt(exact) : null,
        min ? parseInt(min) : null,
        max ? parseInt(max) : null
      );
      
      setResults(response.data || []);
    } catch (error) {
      setError('Greška pri pretraživanju: ' + error.message);
      setResults([]);
    } finally {
      setLoading(false);
    }
  };

  // Clear all search parameters
  const clearSearch = () => {
    setSearchParams({
      luksuzi: [],
      stDomId: '',
      address: '',
      exact: '',
      min: '',
      max: ''
    });
    setResults([]);
    setError('');
    setHasSearched(false);
  };

  // Navigate to student dormitory detail page
  const handleStDomClick = (stDomId) => {
    if (stDomId) {
      navigate(`/st-dom/${stDomId}`);
    }
  };

  // Format luxury amenities for display
  const formatLuksuzi = (luksuzi) => {
    if (!luksuzi || luksuzi.length === 0) return 'Nema';
    return luksuzi.map(l => {
      const found = availableLuksuzi.find(al => al.value === l);
      return found ? found.label : l;
    }).join(', ');
  };

  return (
    <div className="advanced-room-search">
      <div className="navigation-header">
        <button 
          onClick={() => navigate('/dashboard')}
          className="back-button"
        >
          ← Nazad na Dashboard
        </button>
      </div>
      
      <div className="search-header">
        <h2>Napredna pretraga soba</h2>
        <p>Pretražite sobe prema različitim kriterijumima</p>
      </div>

      <form onSubmit={handleSearch} className="search-form">
        {/* Luxury Amenities Section */}
        <div className="form-section">
          <h3>Luksuzni sadržaji</h3>
          <div className="checkbox-group">
            {availableLuksuzi.map(luksuz => (
              <label key={luksuz.value} className="checkbox-label">
                <input
                  type="checkbox"
                  checked={searchParams.luksuzi.includes(luksuz.value)}
                  onChange={() => handleLuksuziChange(luksuz.value)}
                />
                <span className="checkbox-text">{luksuz.label}</span>
              </label>
            ))}
          </div>
        </div>

        {/* Student Dormitory Selection */}
        <div className="form-section">
          <h3>Studentski dom</h3>
          <select
            value={searchParams.stDomId}
            onChange={(e) => handleInputChange('stDomId', e.target.value)}
            className="form-select"
            disabled={stDomsLoading}
          >
            <option value="">Svi domovi</option>
            {stDoms.map(stDom => (
              <option key={stDom.id || stDom._id} value={stDom.id || stDom._id}>
                {stDom.ime} - {stDom.address}
              </option>
            ))}
          </select>
          {stDomsLoading && <p className="loading-text">Učitavanje domova...</p>}
        </div>

        {/* Address Search */}
        <div className="form-section">
          <h3>Adresa</h3>
          <input
            type="text"
            value={searchParams.address}
            onChange={(e) => handleInputChange('address', e.target.value)}
            placeholder="Unesite dio adrese za pretraživanje..."
            className="form-input"
          />
          <p className="help-text">Možete koristiti djelomičnu adresu (npr. "Sarajevo" ili "Zmaja")</p>
        </div>

        {/* Bed Capacity Section */}
        <div className="form-section">
          <h3>Broj kreveta (krevetnost)</h3>
          <div className="capacity-options">
            <div className="capacity-option">
              <label>Tačan broj:</label>
              <input
                type="number"
                min="1"
                value={searchParams.exact}
                onChange={(e) => handleInputChange('exact', e.target.value)}
                placeholder="npr. 2"
                className="form-input number-input"
              />
            </div>
            
            <div className="capacity-range">
              <p className="range-label">ILI raspon:</p>
              <div className="range-inputs">
                <div className="capacity-option">
                  <label>Minimum:</label>
                  <input
                    type="number"
                    min="1"
                    value={searchParams.min}
                    onChange={(e) => handleInputChange('min', e.target.value)}
                    placeholder="npr. 1"
                    className="form-input number-input"
                  />
                </div>
                <div className="capacity-option">
                  <label>Maksimum:</label>
                  <input
                    type="number"
                    min="1"
                    value={searchParams.max}
                    onChange={(e) => handleInputChange('max', e.target.value)}
                    placeholder="npr. 4"
                    className="form-input number-input"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Error Display */}
        {error && <div className="error-message">{error}</div>}

        {/* Action Buttons */}
        <div className="form-actions">
          <button type="submit" className="search-button" disabled={loading}>
            {loading ? 'Pretraživanje...' : 'Pretraži sobe'}
          </button>
          <button type="button" onClick={clearSearch} className="clear-button">
            Obriši sve
          </button>
        </div>
      </form>

      {/* Search Results */}
      {hasSearched && (
        <div className="search-results">
          <div className="results-header">
            <h3>Rezultati pretrage</h3>
            {!loading && (
              <p className="results-count">
                Pronađeno {results.length} soba{results.length !== 1 ? 'a' : ''}
              </p>
            )}
          </div>

          {loading ? (
            <div className="loading-message">Pretraživanje soba...</div>
          ) : results.length > 0 ? (
            <div className="results-grid">
              {results.map((room, index) => (
                <div key={room.id || room._id || index} className="room-card">
                  <div className="room-header">
                    <h4>Soba #{index + 1}</h4>
                    <span className="room-capacity">{room.krevetnost} krevet{room.krevetnost !== 1 ? 'a' : ''}</span>
                  </div>
                  
                  <div className="room-details">
                    <div className="detail-row">
                      <span className="detail-label">Luksuzni sadržaji:</span>
                      <span className="detail-value">{formatLuksuzi(room.luksuzi)}</span>
                    </div>
                    
                    {room.st_dom && (
                      <>
                        <div className="detail-row">
                          <span className="detail-label">Studentski dom:</span>
                          <span 
                            className="detail-value clickable-link"
                            onClick={() => handleStDomClick(room.st_dom.id || room.st_dom._id)}
                          >
                            {room.st_dom.ime}
                          </span>
                        </div>
                        <div className="detail-row">
                          <span className="detail-label">Adresa:</span>
                          <span className="detail-value">{room.st_dom.address}</span>
                        </div>
                        {room.st_dom.telephone_number && (
                          <div className="detail-row">
                            <span className="detail-label">Telefon:</span>
                            <span className="detail-value">{room.st_dom.telephone_number}</span>
                          </div>
                        )}
                        {room.st_dom.email && (
                          <div className="detail-row">
                            <span className="detail-label">Email:</span>
                            <span className="detail-value">{room.st_dom.email}</span>
                          </div>
                        )}
                      </>
                    )}
                  </div>
                  <ApplyToRoomButton 
                    room={room} 
                    stDom={room.st_dom}
                    onSuccess={() => {
                      // Optionally refresh data or show success message
                      console.log('Application submitted successfully');
                    }}
                  />
                </div>
              ))}
            </div>
          ) : (
            <div className="no-results">
              <p>Nema soba koje odgovaraju vašim kriterijumima pretrage.</p>
              <p>Pokušajte sa drugačijim parametrima.</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default AdvancedRoomSearch;
