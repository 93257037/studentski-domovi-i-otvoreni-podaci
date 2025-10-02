import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { stDomService } from '../services/stDomService';
import StDomRoomSearch from './StDomRoomSearch';
import './StDomDetail.css';

const StDomDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [stDom, setStDom] = useState(null);
  const [rooms, setRooms] = useState([]);
  const [filteredRooms, setFilteredRooms] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (id) {
      fetchStDomDetails();
    }
  }, [id]);

  const fetchStDomDetails = async () => {
    setLoading(true);
    setError('');
    
    try {
      // Fetch dormitory details
      const stDomResponse = await stDomService.getStDom(id);
      setStDom(stDomResponse.st_dom);
      
      // Fetch rooms for this dormitory
      const roomsResponse = await stDomService.getStDomRooms(id);
      const roomsData = roomsResponse.sobas || [];
      setRooms(roomsData);
      setFilteredRooms(roomsData); // Initially show all rooms
    } catch (error) {
      setError('Greška pri učitavanju podataka: ' + error.message);
    } finally {
      setLoading(false);
    }
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

  const formatLuksuzi = (luksuzi) => {
    if (!luksuzi || luksuzi.length === 0) return 'Nema luksuznih sadržaja';
    
    const luksuziMap = {
      'klima': 'Klima',
      'terasa': 'Terasa',
      'sopstveno kupatilo': 'Sopstveno kupatilo',
      'áram': 'Stram',
      'ablak': 'Ablak',
      'neisvrljan zid': 'Neisvrljan zid'
    };
    
    return luksuzi.map(l => luksuziMap[l] || l).join(', ');
  };

  const handleBack = () => {
    navigate('/dashboard');
  };

  // Handle search results from the search component
  const handleSearchResults = (searchResults) => {
    if (searchResults === null) {
      // Reset to show all rooms
      setFilteredRooms(rooms);
    } else {
      // Show filtered results - extract just the room data since the search returns rooms with st_dom info
      const roomsOnly = searchResults.map(result => ({
        id: result.id,
        krevetnost: result.krevetnost,
        luksuzi: result.luksuzi,
        created_at: result.created_at
      }));
      setFilteredRooms(roomsOnly);
    }
  };

  if (loading) {
    return (
      <div className="st-dom-detail-container">
        <div className="loading-message">Učitavanje podataka...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="st-dom-detail-container">
        <div className="error-message">{error}</div>
        <button onClick={handleBack} className="back-button">
          Povratak na dashboard
        </button>
      </div>
    );
  }

  if (!stDom) {
    return (
      <div className="st-dom-detail-container">
        <div className="error-message">Studentski dom nije pronađen</div>
        <button onClick={handleBack} className="back-button">
          Povratak na dashboard
        </button>
      </div>
    );
  }

  return (
    <div className="st-dom-detail-container">
      <div className="st-dom-detail-header">
        <button onClick={handleBack} className="back-button">
          ← Povratak na dashboard
        </button>
        <h1>{stDom.ime}</h1>
      </div>

      <div className="st-dom-detail-content">
        {/* Basic Information */}
        <div className="detail-section">
          <h2>Osnovne informacije</h2>
          <div className="info-grid">
            <div className="info-item">
              <span className="label">Ime:</span>
              <span className="value">{stDom.ime}</span>
            </div>
            <div className="info-item">
              <span className="label">Adresa:</span>
              <span className="value">{stDom.address}</span>
            </div>
            <div className="info-item">
              <span className="label">Telefon:</span>
              <span className="value">{stDom.telephone_number}</span>
            </div>
            <div className="info-item">
              <span className="label">Email:</span>
              <span className="value">
                <a href={`mailto:${stDom.email}`}>{stDom.email}</a>
              </span>
            </div>
            <div className="info-item">
              <span className="label">Kreiran:</span>
              <span className="value">{formatDate(stDom.created_at)}</span>
            </div>
            <div className="info-item">
              <span className="label">Ažuriran:</span>
              <span className="value">{formatDate(stDom.updated_at)}</span>
            </div>
          </div>
        </div>

        {/* Room Search */}
        {stDom && (
          <StDomRoomSearch 
            stDomId={stDom.id || stDom._id} 
            onSearchResults={handleSearchResults}
          />
        )}

        {/* Rooms Information */}
        <div className="detail-section">
          <h2>
            {filteredRooms.length === rooms.length 
              ? `Dostupne sobe (${rooms.length})`
              : `Filtrirane sobe (${filteredRooms.length} od ${rooms.length})`
            }
          </h2>
          {filteredRooms.length > 0 ? (
            <div className="rooms-grid">
              {filteredRooms.map((room) => (
                <div key={room.id} className="room-card">
                  <div className="room-header">
                    <h3>Soba {room.id}</h3>
                    <span className="room-capacity">{room.krevetnost} kreveta</span>
                  </div>
                  <div className="room-details">
                    <div className="room-info">
                      <span className="label">Kapacitet:</span>
                      <span className="value">{room.krevetnost} osoba</span>
                    </div>
                    <div className="room-info">
                      <span className="label">Luksuzni sadržaji:</span>
                      <span className="value">{formatLuksuzi(room.luksuzi)}</span>
                    </div>
                    <div className="room-info">
                      <span className="label">Kreirana:</span>
                      <span className="value">{formatDate(room.created_at)}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="no-rooms">
              {rooms.length === 0 ? (
                <p>Trenutno nema dostupnih soba u ovom studentskom domu.</p>
              ) : (
                <p>Nema soba koje odgovaraju vašim kriterijumima pretrage.</p>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default StDomDetail;
