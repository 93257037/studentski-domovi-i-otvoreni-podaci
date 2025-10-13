import React, { useState, useEffect } from 'react';
import { openDataService } from '../services/openDataService';
import { stDomService } from '../services/stDomService';
import './ScheduleRepairModal.css';

const ScheduleRepairModal = ({ show, onClose, token }) => {
  const [step, setStep] = useState(1); // 1: room search, 2: repair details
  const [searchFilters, setSearchFilters] = useState({
    dormId: '',
    minCapacity: '',
    maxCapacity: '',
    amenities: [],
    onlyAvailable: false,
    limit: 20
  });
  const [roomResults, setRoomResults] = useState([]);
  const [dormList, setDormList] = useState([]);
  const [availableAmenities, setAvailableAmenities] = useState([]);
  const [selectedRoom, setSelectedRoom] = useState(null);
  const [description, setDescription] = useState('');
  const [completionDate, setCompletionDate] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [searchLoading, setSearchLoading] = useState(false);

  useEffect(() => {
    if (show) {
      loadInitialData();
    }
  }, [show]);

  const loadInitialData = async () => {
    try {
      const dormsResponse = await openDataService.getDormList();
      setDormList(dormsResponse.dorms || []);

      const amenitiesResponse = await openDataService.getAvailableAmenities();
      setAvailableAmenities(amenitiesResponse.amenities || []);
    } catch (error) {
      console.error('Error loading initial data:', error);
    }
  };

  const handleSearch = async (e) => {
    e.preventDefault();
    setSearchLoading(true);
    setError('');
    
    try {
      const response = await openDataService.searchAvailableRooms(searchFilters);
      setRoomResults(response.rooms || []);
      
      if (!response.rooms || response.rooms.length === 0) {
        setError('Nema rezultata pretrage. Pokušajte sa različitim filterima.');
      }
    } catch (error) {
      setError('Greška pri pretraživanju soba: ' + (error.response?.data?.error || error.message));
      setRoomResults([]);
    } finally {
      setSearchLoading(false);
    }
  };

  const handleRoomSelect = (room) => {
    setSelectedRoom(room);
    setStep(2);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!description.trim()) {
      setError('Molimo unesite opis popravke');
      return;
    }

    if (!completionDate) {
      setError('Molimo unesite predviđeni datum završetka');
      return;
    }

    setLoading(true);
    setError('');

    try {
      // Convert date to ISO 8601 format
      const isoDate = new Date(completionDate).toISOString();
      
      await stDomService.scheduleRepair(
        selectedRoom.room_id,
        description,
        isoDate
      );

      alert('Popravka uspešno zakazana!');
      handleClose();
    } catch (error) {
      setError('Greška pri zakazivanju popravke: ' + (error.message || 'Nepoznata greška'));
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setStep(1);
    setSearchFilters({
      dormId: '',
      minCapacity: '',
      maxCapacity: '',
      amenities: [],
      onlyAvailable: false,
      limit: 20
    });
    setRoomResults([]);
    setSelectedRoom(null);
    setDescription('');
    setCompletionDate('');
    setError('');
    onClose();
  };

  const handleBack = () => {
    setStep(1);
    setSelectedRoom(null);
    setDescription('');
    setCompletionDate('');
    setError('');
  };

  const handleAmenityToggle = (amenity) => {
    setSearchFilters(prev => ({
      ...prev,
      amenities: prev.amenities.includes(amenity)
        ? prev.amenities.filter(a => a !== amenity)
        : [...prev.amenities, amenity]
    }));
  };

  if (!show) return null;

  return (
    <div className="modal-overlay" onClick={handleClose}>
      <div className="modal-content repair-modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>{step === 1 ? 'Pretraži Sobu' : 'Zakaži Popravku'}</h2>
          <button className="close-button" onClick={handleClose}>×</button>
        </div>

        {error && <div className="error-message">{error}</div>}

        {step === 1 ? (
          <div className="room-search-step">
            <form onSubmit={handleSearch} className="search-form">
              <div className="form-group">
                <label>Dom:</label>
                <select
                  value={searchFilters.dormId}
                  onChange={(e) => setSearchFilters({...searchFilters, dormId: e.target.value})}
                >
                  <option value="">Svi domovi</option>
                  {dormList.map(dorm => (
                    <option key={dorm.id} value={dorm.id}>
                      {dorm.name}
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-row">
                <div className="form-group">
                  <label>Min. Kapacitet:</label>
                  <input
                    type="number"
                    value={searchFilters.minCapacity}
                    onChange={(e) => setSearchFilters({...searchFilters, minCapacity: e.target.value})}
                    min="1"
                  />
                </div>
                <div className="form-group">
                  <label>Maks. Kapacitet:</label>
                  <input
                    type="number"
                    value={searchFilters.maxCapacity}
                    onChange={(e) => setSearchFilters({...searchFilters, maxCapacity: e.target.value})}
                    min="1"
                  />
                </div>
              </div>

              <div className="form-group">
                <label>Pogodnosti:</label>
                <div className="amenities-grid">
                  {availableAmenities.map(amenity => (
                    <label key={amenity} className="amenity-checkbox">
                      <input
                        type="checkbox"
                        checked={searchFilters.amenities.includes(amenity)}
                        onChange={() => handleAmenityToggle(amenity)}
                      />
                      {amenity}
                    </label>
                  ))}
                </div>
              </div>

              <div className="form-group">
                <label className="checkbox-label">
                  <input
                    type="checkbox"
                    checked={searchFilters.onlyAvailable}
                    onChange={(e) => setSearchFilters({...searchFilters, onlyAvailable: e.target.checked})}
                  />
                  Samo dostupne sobe
                </label>
              </div>

              <button type="submit" className="btn btn-primary" disabled={searchLoading}>
                {searchLoading ? 'Pretraživanje...' : 'Pretraži'}
              </button>
            </form>

            {roomResults.length > 0 && (
              <div className="search-results">
                <h3>Rezultati ({roomResults.length})</h3>
                <div className="results-table-container">
                  <table className="results-table">
                    <thead>
                      <tr>
                        <th>ID Sobe</th>
                        <th>Dom</th>
                        <th>Adresa</th>
                        <th>Kapacitet</th>
                        <th>Popunjeno</th>
                        <th>Pogodnosti</th>
                        <th>Akcija</th>
                      </tr>
                    </thead>
                    <tbody>
                      {roomResults.map((room, index) => (
                        <tr key={index}>
                          <td><code>{room.room_id}</code></td>
                          <td>{room.dorm_name}</td>
                          <td>{room.dorm_address}</td>
                          <td>{room.capacity}</td>
                          <td>{room.occupied}/{room.capacity}</td>
                          <td>{room.amenities ? room.amenities.join(', ') : 'Nema'}</td>
                          <td>
                            <button
                              className="btn btn-small"
                              onClick={() => handleRoomSelect(room)}
                            >
                              Izaberi
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </div>
        ) : (
          <div className="repair-details-step">
            <div className="selected-room-info">
              <h3>Izabrana soba:</h3>
              <p><strong>Dom:</strong> {selectedRoom.dorm_name}</p>
              <p><strong>Adresa:</strong> {selectedRoom.dorm_address}</p>
              <p><strong>Kapacitet:</strong> {selectedRoom.capacity}</p>
            </div>

            <form onSubmit={handleSubmit} className="repair-form">
              <div className="form-group">
                <label>Opis popravke: *</label>
                <textarea
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Unesite kratak opis potrebne popravke..."
                  rows="4"
                  required
                />
              </div>

              <div className="form-group">
                <label>Predviđeni datum završetka: *</label>
                <input
                  type="date"
                  value={completionDate}
                  onChange={(e) => setCompletionDate(e.target.value)}
                  min={new Date().toISOString().split('T')[0]}
                  required
                />
              </div>

              <div className="button-group">
                <button type="button" className="btn btn-secondary" onClick={handleBack}>
                  Nazad
                </button>
                <button type="submit" className="btn btn-primary" disabled={loading}>
                  {loading ? 'Zakazivanje...' : 'Zakaži Popravku'}
                </button>
              </div>
            </form>
          </div>
        )}
      </div>
    </div>
  );
};

export default ScheduleRepairModal;

