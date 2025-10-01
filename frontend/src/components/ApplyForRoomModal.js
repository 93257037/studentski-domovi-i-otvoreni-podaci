import React, { useState, useEffect } from 'react';
import { stDomService } from '../services/stDomService';
import './ApplyForRoomModal.css';

const ApplyForRoomModal = ({ isOpen, onClose, onSuccess }) => {
  const [step, setStep] = useState(1); // 1: Select Dorm, 2: Select Room, 3: Fill Details
  const [stDoms, setStDoms] = useState([]);
  const [rooms, setRooms] = useState([]);
  const [selectedStDom, setSelectedStDom] = useState(null);
  const [selectedRoom, setSelectedRoom] = useState(null);
  const [formData, setFormData] = useState({
    broj_indexa: '',
    prosek: ''
  });
  const [loading, setLoading] = useState(false);
  const [loadingStDoms, setLoadingStDoms] = useState(false);
  const [loadingRooms, setLoadingRooms] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (isOpen) {
      fetchStDoms();
      // Reset state when modal opens
      setStep(1);
      setSelectedStDom(null);
      setSelectedRoom(null);
      setRooms([]);
      setFormData({ broj_indexa: '', prosek: '' });
      setError('');
    }
  }, [isOpen]);

  const fetchStDoms = async () => {
    setLoadingStDoms(true);
    setError('');
    try {
      const response = await stDomService.getAllStDoms();
      setStDoms(response.st_doms || []);
    } catch (err) {
      setError('Greška pri učitavanju studentskih domova: ' + err.message);
    } finally {
      setLoadingStDoms(false);
    }
  };

  const handleSelectStDom = async (stDom) => {
    setSelectedStDom(stDom);
    setLoadingRooms(true);
    setError('');
    try {
      const response = await stDomService.getStDomRooms(stDom.id);
      console.log('=== ROOMS RESPONSE ===');
      console.log('Response:', response);
      console.log('======================');
      setRooms(response.sobas || []);
      setStep(2);
    } catch (err) {
      setError('Greška pri učitavanju soba: ' + err.message);
    } finally {
      setLoadingRooms(false);
    }
  };

  const handleSelectRoom = (room) => {
    setSelectedRoom(room);
    setStep(3);
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const aplikacijaData = {
        broj_indexa: formData.broj_indexa,
        prosek: parseInt(formData.prosek),
        soba_id: selectedRoom.id
      };

      console.log('=== SUBMITTING APPLICATION ===');
      console.log('JSON Payload:', JSON.stringify(aplikacijaData, null, 2));
      console.log('==============================');

      await stDomService.createAplikacija(aplikacijaData);
      onSuccess();
      onClose();
      // Reset form
      setStep(1);
      setSelectedStDom(null);
      setSelectedRoom(null);
      setRooms([]);
      setFormData({ broj_indexa: '', prosek: '' });
    } catch (err) {
      console.error('=== APPLICATION ERROR ===');
      console.error('Error:', err);
      console.error('=========================');
      setError(err.message || 'Greška pri kreiranju aplikacije');
    } finally {
      setLoading(false);
    }
  };

  const handleBack = () => {
    if (step === 2) {
      setStep(1);
      setSelectedStDom(null);
      setRooms([]);
    } else if (step === 3) {
      setStep(2);
      setSelectedRoom(null);
      setFormData({ broj_indexa: '', prosek: '' });
    }
  };

  const handleClose = () => {
    if (!loading) {
      onClose();
      setStep(1);
      setSelectedStDom(null);
      setSelectedRoom(null);
      setRooms([]);
      setFormData({ broj_indexa: '', prosek: '' });
      setError('');
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

  if (!isOpen) return null;

  return (
    <div className="modal-overlay">
      <div className="modal apply-room-modal">
        <div className="modal-header">
          <h3>
            {step === 1 && 'Odaberite studentski dom'}
            {step === 2 && 'Odaberite sobu'}
            {step === 3 && 'Popunite detalje aplikacije'}
          </h3>
          <button 
            className="close-button" 
            onClick={handleClose}
            disabled={loading}
          >
            ×
          </button>
        </div>

        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        {/* Step 1: Select Student Dorm */}
        {step === 1 && (
          <div className="apply-step">
            {loadingStDoms ? (
              <div className="loading-text">Učitavanje studentskih domova...</div>
            ) : (
              <div className="stdom-list">
                {stDoms.length === 0 ? (
                  <div className="no-data">Nema dostupnih studentskih domova</div>
                ) : (
                  stDoms.map(stDom => (
                    <div 
                      key={stDom.id} 
                      className="stdom-card"
                      onClick={() => handleSelectStDom(stDom)}
                    >
                      <h4>{stDom.ime}</h4>
                      <p><strong>Adresa:</strong> {stDom.address}</p>
                      <p><strong>Telefon:</strong> {stDom.telephone_number}</p>
                      <p><strong>Email:</strong> {stDom.email}</p>
                    </div>
                  ))
                )}
              </div>
            )}
          </div>
        )}

        {/* Step 2: Select Room */}
        {step === 2 && (
          <div className="apply-step">
            <div className="selected-info">
              <strong>Odabrani dom:</strong> {selectedStDom?.ime}
            </div>
            {loadingRooms ? (
              <div className="loading-text">Učitavanje soba...</div>
            ) : (
              <div className="room-list">
                {rooms.length === 0 ? (
                  <div className="no-data">Nema dostupnih soba u ovom domu</div>
                ) : (
                  rooms.map(room => (
                    <div 
                      key={room.id} 
                      className="room-card"
                      onClick={() => handleSelectRoom(room)}
                    >
                      <h4>Soba (ID: {room.id.slice(-6)})</h4>
                      <p><strong>Broj kreveta:</strong> {room.krevetnost}</p>
                      {room.luksuzi && room.luksuzi.length > 0 && (
                        <div className="amenities">
                          <strong>Luksuzi:</strong>
                          <ul>
                            {room.luksuzi.map((luksuz, index) => (
                              <li key={index}>{formatLuksuz(luksuz)}</li>
                            ))}
                          </ul>
                        </div>
                      )}
                    </div>
                  ))
                )}
              </div>
            )}
            <div className="modal-buttons">
              <button 
                type="button" 
                onClick={handleBack}
                className="back-button"
              >
                Nazad
              </button>
            </div>
          </div>
        )}

        {/* Step 3: Fill Application Details */}
        {step === 3 && (
          <div className="apply-step">
            <div className="selected-info">
              <p><strong>Dom:</strong> {selectedStDom?.ime}</p>
              <p><strong>Soba:</strong> {selectedRoom?.krevetnost} kreveta</p>
            </div>
            <form onSubmit={handleSubmit} className="application-form">
              <div className="form-group">
                <label htmlFor="broj_indexa">Broj indeksa *</label>
                <input
                  type="text"
                  id="broj_indexa"
                  name="broj_indexa"
                  value={formData.broj_indexa}
                  onChange={handleInputChange}
                  required
                  disabled={loading}
                  placeholder="Unesite broj indeksa (npr. RI12/2021)"
                />
              </div>

              <div className="form-group">
                <label htmlFor="prosek">Prosječna ocjena (6-10) *</label>
                <input
                  type="number"
                  id="prosek"
                  name="prosek"
                  value={formData.prosek}
                  onChange={handleInputChange}
                  required
                  min="6"
                  max="10"
                  step="1"
                  disabled={loading}
                  placeholder="Unesite prosječnu ocjenu"
                />
              </div>

              <div className="modal-buttons">
                <button 
                  type="button" 
                  onClick={handleBack}
                  className="back-button"
                  disabled={loading}
                >
                  Nazad
                </button>
                <button 
                  type="submit" 
                  className="submit-button"
                  disabled={loading}
                >
                  {loading ? 'Slanje...' : 'Apliciraj'}
                </button>
              </div>
            </form>
          </div>
        )}
      </div>
    </div>
  );
};

export default ApplyForRoomModal;

