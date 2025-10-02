import React, { useState } from 'react';
import { openDataService } from '../services/openDataService';
import './StDomRoomSearch.css';

const StDomRoomSearch = ({ stDomId, onSearchResults }) => {
  // Search form state
  const [searchParams, setSearchParams] = useState({
    luksuzi: [],
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
  
  // Search state
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

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
    
    try {
      const { luksuzi, exact, min, max } = searchParams;
      
      const response = await openDataService.advancedFilterRooms(
        luksuzi.length > 0 ? luksuzi : null,
        stDomId, // Always filter by the current st_dom
        null, // No address filtering as requested
        exact ? parseInt(exact) : null,
        min ? parseInt(min) : null,
        max ? parseInt(max) : null
      );
      
      // Pass results back to parent component
      onSearchResults(response.data || []);
    } catch (error) {
      setError('Greška pri pretraživanju: ' + error.message);
      onSearchResults([]);
    } finally {
      setLoading(false);
    }
  };

  // Clear all search parameters
  const clearSearch = () => {
    setSearchParams({
      luksuzi: [],
      exact: '',
      min: '',
      max: ''
    });
    setError('');
    // Reset to show all rooms
    onSearchResults(null);
  };

  return (
    <div className="st-dom-room-search">
      <div className="search-header">
        <h3>Pretraži sobe u ovom domu</h3>
        <p>Filtrirajte sobe prema vašim kriterijumima</p>
      </div>

      <form onSubmit={handleSearch} className="search-form">
        {/* Luxury Amenities Section */}
        <div className="form-section">
          <h4>Luksuzni sadržaji</h4>
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

        {/* Bed Capacity Section */}
        <div className="form-section">
          <h4>Broj kreveta (krevetnost)</h4>
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
            Prikaži sve sobe
          </button>
        </div>
      </form>
    </div>
  );
};

export default StDomRoomSearch;
