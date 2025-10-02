import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { stDomService } from '../services/stDomService';
import { openDataService } from '../services/openDataService';
import { useAuth } from '../contexts/AuthContext';
import ApplyToRoomButton from './ApplyToRoomButton';

import './RoomDetail.css';

const RoomDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { token } = useAuth();
  const [room, setRoom] = useState(null);
  const [stDom, setStDom] = useState(null);
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [applicationsLoading, setApplicationsLoading] = useState(false);
  const [error, setError] = useState('');
  const [applicationsError, setApplicationsError] = useState('');

  useEffect(() => {
    if (id) {
      fetchRoomDetails();
      fetchRoomApplications();
    }
  }, [id, token]);

  const fetchRoomDetails = async () => {
    setLoading(true);
    setError('');
    
    try {
      // Fetch room details
      const roomResponse = await stDomService.getRoom(id);
      const roomData = roomResponse.soba;
      setRoom(roomData);
      
      // Fetch dormitory details for the room
      if (roomData.st_dom_id) {
        const stDomResponse = await stDomService.getStDom(roomData.st_dom_id);
        setStDom(stDomResponse.st_dom);
      }
    } catch (error) {
      setError('Greška pri učitavanju podataka: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  const fetchRoomApplications = async () => {
    if (!token) return; // Skip if no authentication token
    
    setApplicationsLoading(true);
    setApplicationsError('');
    
    try {
      const response = await openDataService.getPrihvaceneAplikacijeForRoom(id, token);
      setApplications(response.data || []);
    } catch (error) {
      setApplicationsError('Greška pri učitavanju aplikacija: ' + error.message);
      console.error('Error fetching room applications:', error);
    } finally {
      setApplicationsLoading(false);
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
    if (!luksuzi || luksuzi.length === 0) return 'Nema';
    
    const luksuziMap = {
      'klima': 'Klima uređaj',
      'terasa': 'Terasa',
      'sopstveno kupatilo': 'Sopstveno kupatilo',
      'áram': 'Struja',
      'ablak': 'Prozor',
      'neisvrljan zid': 'Neisvrljan zid'
    };
    
    return luksuzi.map(l => luksuziMap[l] || l).join(', ');
  };

  if (loading) {
    return (
      <div className="room-detail-container">
        <div className="loading-text">Učitavanje podataka o sobi...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="room-detail-container">
        <div className="error-message">{error}</div>
        <button onClick={() => navigate(-1)} className="back-button">
          Nazad
        </button>
      </div>
    );
  }

  if (!room) {
    return (
      <div className="room-detail-container">
        <div className="error-message">Soba nije pronađena</div>
        <button onClick={() => navigate(-1)} className="back-button">
          Nazad
        </button>
      </div>
    );
  }

  return (
    <div className="room-detail-container">
      <div className="room-detail-header">
        <button onClick={() => navigate(-1)} className="back-button">
          ← Nazad
        </button>
        <h1>Detalji sobe</h1>
      </div>

      <div className="room-detail-content">
        <div className="room-info-section">
          <h2>Informacije o sobi</h2>
          <div className="room-info-grid">
            <div className="info-item">
              <span className="info-label">ID sobe:</span>
              <span className="info-value">{room.id}</span>
            </div>
            <div className="info-item">
              <span className="info-label">Kapacitet:</span>
              <span className="info-value">{room.krevetnost} {room.krevetnost === 1 ? 'krevet' : 'kreveta'}</span>
            </div>
            <div className="info-item">
              <span className="info-label">Luksuzni sadržaji:</span>
              <span className="info-value">{formatLuksuzi(room.luksuzi)}</span>
            </div>
            <div className="info-item">
              <span className="info-label">Kreirana:</span>
              <span className="info-value">{formatDate(room.created_at)}</span>
            </div>
            {room.updated_at && (
              <div className="info-item">
                <span className="info-label">Poslednja izmena:</span>
                <span className="info-value">{formatDate(room.updated_at)}</span>
              </div>
            )}
          </div>
        </div>

        {stDom && (
          <div className="dormitory-info-section">
            <h2>Informacije o domu</h2>
            <div className="dormitory-info-grid">
              <div className="info-item">
                <span className="info-label">Naziv doma:</span>
                <span className="info-value">{stDom.ime}</span>
              </div>
              <div className="info-item">
                <span className="info-label">Adresa:</span>
                <span className="info-value">{stDom.address}</span>
              </div>
              <div className="info-item">
                <span className="info-label">Telefon:</span>
                <span className="info-value">{stDom.telephone_number}</span>
              </div>
              <div className="info-item">
                <span className="info-label">Email:</span>
                <span className="info-value">{stDom.email}</span>
              </div>
            </div>
          </div>
        )}

        {/* Room Applications Section */}
        <div className="applications-section">
          <h2>Prihvaćene aplikacije za ovu sobu</h2>
          {applicationsLoading ? (
            <div className="applications-loading">
              <div className="loading-spinner"></div>
              <span>Učitavanje aplikacija...</span>
            </div>
          ) : applicationsError ? (
            <div className="applications-error">
              <p>{applicationsError}</p>
              <button 
                onClick={fetchRoomApplications} 
                className="retry-button"
              >
                Pokušaj ponovo
              </button>
            </div>
          ) : applications.length > 0 ? (
            <div className="applications-list">
              <p className="applications-count">
                Ukupno prihvaćenih aplikacija: <strong>{applications.length}</strong>
              </p>
              <div className="applications-grid">
                {applications.map((application, index) => (
                  <div key={application.id || index} className="application-card">
                    <div className="application-header">
                      <h4>Aplikacija #{application.id}</h4>
                      <span className="academic-year">{application.academic_year}</span>
                    </div>
                    <div className="application-details">
                      <div className="detail-row">
                        <span className="label">Broj indeksa:</span>
                        <span className="value">{application.broj_indexa}</span>
                      </div>
                      <div className="detail-row">
                        <span className="label">Prosek:</span>
                        <span className="value">{application.prosek}</span>
                      </div>
                      <div className="detail-row">
                        <span className="label">Kreirana:</span>
                        <span className="value">{formatDate(application.created_at)}</span>
                      </div>
                      {application.updated_at && (
                        <div className="detail-row">
                          <span className="label">Ažurirana:</span>
                          <span className="value">{formatDate(application.updated_at)}</span>
                        </div>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ) : (
            <div className="no-applications">
              <p>Trenutno nema prihvaćenih aplikacija za ovu sobu.</p>
            </div>
          )}
        </div>

        <div className="room-actions">
          <ApplyToRoomButton 
            room={room} 
            stDom={stDom}
            onSuccess={() => {
              console.log('Application submitted successfully');
              // Refresh applications after successful application
              fetchRoomApplications();
            }}
          />
        </div>
      </div>
    </div>
  );
};

export default RoomDetail;
