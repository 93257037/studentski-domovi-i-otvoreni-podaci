import React, { useState, useEffect } from 'react';
import { stDomService } from '../services/stDomService';
import './MyRoomInfo.css';

const MyRoomInfo = ({ onCheckout }) => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [roomInfo, setRoomInfo] = useState(null);
  const [checkoutLoading, setCheckoutLoading] = useState(false);

  useEffect(() => {
    fetchMyRoom();
  }, []);

  const fetchMyRoom = async () => {
    setLoading(true);
    setError('');
    try {
      // Get my accepted applications
      const response = await stDomService.getMyPrihvaceneAplikacije();
      
      if (!response.prihvacene_aplikacije || response.prihvacene_aplikacije.length === 0) {
        setRoomInfo(null);
        return;
      }

      const myApp = response.prihvacene_aplikacije[0]; // User should only have one

      // Fetch room details
      const roomResponse = await stDomService.getRoom(myApp.soba_id);
      const room = roomResponse.soba;

      // Fetch student dorm details
      const stDomResponse = await stDomService.getStDom(room.st_dom_id);
      const stDom = stDomResponse.st_dom;

      setRoomInfo({
        application: myApp,
        room: room,
        stDom: stDom
      });
    } catch (err) {
      console.error('Error fetching room info:', err);
      setError('Greška pri učitavanju informacija o sobi: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleCheckout = async () => {
    if (!window.confirm('Da li ste sigurni da želite da se odjavite iz sobe? Ova akcija je nepovratna.')) {
      return;
    }

    setCheckoutLoading(true);
    setError('');
    try {
      await stDomService.checkoutFromRoom();
      alert('Uspješno ste se odjavili iz sobe!');
      setRoomInfo(null);
      if (onCheckout) onCheckout();
    } catch (err) {
      setError('Greška pri odjavi iz sobe: ' + err.message);
      alert('Greška: ' + err.message);
    } finally {
      setCheckoutLoading(false);
    }
  };

  const formatLuksuz = (luksuz) => {
    const labels = {
      'klima': 'Klima',
      'terasa': 'Terasa',
      'sopstveno kupatilo': 'Sopstveno kupatilo',
      'áram': 'Áram',
      'ablak': 'Ablak',
      'neisvrljan zid': 'Neisvrljan zid'
    };
    return labels[luksuz] || luksuz;
  };


  if (loading) {
    return (
      <div className="my-room-info loading">
        <p>Učitavanje informacija o sobi...</p>
      </div>
    );
  }

  if (!roomInfo) {
    return (
      <div className="my-room-info no-room">
        <div className="no-room-icon">🏠</div>
        <h3>Nemate dodijeljenu sobu</h3>
        <p>Aplicirajte za sobu kako biste dobili smještaj u studentskom domu.</p>
      </div>
    );
  }

  return (
    <div className="my-room-info">
      {error && <div className="error-message">{error}</div>}
      
      <div className="room-header">
        <h3>🏠 Moja Soba</h3>
        <button 
          onClick={handleCheckout}
          className="checkout-button"
          disabled={checkoutLoading}
        >
          {checkoutLoading ? 'Odjava...' : '🚪 Odjavi se iz sobe'}
        </button>
      </div>

      <div className="room-details-container">
        {/* Student Dorm Info */}
        <div className="info-section">
          <h4>📍 Studentski Dom</h4>
          <div className="info-grid">
            <div className="info-item">
              <span className="label">Ime:</span>
              <span className="value">{roomInfo.stDom.ime}</span>
            </div>
            <div className="info-item">
              <span className="label">Adresa:</span>
              <span className="value">{roomInfo.stDom.address}</span>
            </div>
            <div className="info-item">
              <span className="label">Telefon:</span>
              <span className="value">{roomInfo.stDom.telephone_number}</span>
            </div>
            <div className="info-item">
              <span className="label">Email:</span>
              <span className="value">{roomInfo.stDom.email}</span>
            </div>
          </div>
        </div>

        {/* Room Info */}
        <div className="info-section">
          <h4>🛏️ Informacije o Sobi</h4>
          <div className="info-grid">
            <div className="info-item">
              <span className="label">Broj kreveta:</span>
              <span className="value">{roomInfo.room.krevetnost}</span>
            </div>
            <div className="info-item">
              <span className="label">Akademska godina:</span>
              <span className="value">{roomInfo.application.academic_year}</span>
            </div>
          </div>
          
          {roomInfo.room.luksuzi && roomInfo.room.luksuzi.length > 0 && (
            <div className="amenities-section">
              <span className="label">Luksuzi:</span>
              <div className="amenities-list">
                {roomInfo.room.luksuzi.map((luksuz, index) => (
                  <span key={index} className="amenity-badge">
                    ✓ {formatLuksuz(luksuz)}
                  </span>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default MyRoomInfo;

