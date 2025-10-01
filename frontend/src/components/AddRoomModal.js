import React, { useState, useEffect } from 'react';
import { stDomService } from '../services/stDomService';
import './AddRoomModal.css';

const AddRoomModal = ({ isOpen, onClose, onSuccess }) => {
  const [formData, setFormData] = useState({
    st_dom_id: '',
    krevetnost: '',
    luksuzi: []
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [stDoms, setStDoms] = useState([]);
  const [loadingStDoms, setLoadingStDoms] = useState(false);

  const availableLuksuzi = [
    { value: 'klima', label: 'Klima' },
    { value: 'terasa', label: 'Terasa' },
    { value: 'sopstveno kupatilo', label: 'Sopstveno kupatilo' },
    { value: 'áram', label: 'Áram' },
    { value: 'ablak', label: 'Ablak' },
    { value: 'neisvrljan zid', label: 'Neisvrljan zid' }
  ];

  useEffect(() => {
    if (isOpen) {
      fetchStDoms();
    }
  }, [isOpen]);

  const fetchStDoms = async () => {
    setLoadingStDoms(true);
    try {
      const response = await stDomService.getAllStDoms();
      setStDoms(response.st_doms || []);
    } catch (err) {
      setError('Greška pri učitavanju studentskih domova');
    } finally {
      setLoadingStDoms(false);
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleLuksuziChange = (e) => {
    const { value, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      luksuzi: checked
        ? [...prev.luksuzi, value]
        : prev.luksuzi.filter(l => l !== value)
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const roomData = {
        st_dom_id: formData.st_dom_id,
        krevetnost: parseInt(formData.krevetnost),
        luksuzi: formData.luksuzi
      };
      
      console.log('=== SENDING ROOM DATA TO BACKEND ===');
      console.log('JSON Payload:', JSON.stringify(roomData, null, 2));
      console.log('Luksuzi array:', roomData.luksuzi);
      console.log('=====================================');
      
      await stDomService.createRoom(roomData);
      onSuccess();
      onClose();
      // Reset form
      setFormData({
        st_dom_id: '',
        krevetnost: '',
        luksuzi: []
      });
    } catch (err) {
      console.error('=== ROOM CREATION ERROR ===');
      console.error('Error message:', err.message);
      console.error('Full error:', err);
      console.error('===========================');
      setError(err.message || 'Greška pri kreiranju sobe');
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    if (!loading) {
      onClose();
      setError('');
      // Reset form
      setFormData({
        st_dom_id: '',
        krevetnost: '',
        luksuzi: []
      });
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal-overlay">
      <div className="modal add-room-modal">
        <div className="modal-header">
          <h3>Dodaj novu sobu</h3>
          <button 
            className="close-button" 
            onClick={handleClose}
            disabled={loading}
          >
            ×
          </button>
        </div>

        <form onSubmit={handleSubmit} className="room-form">
          {error && (
            <div className="error-message">
              {error}
            </div>
          )}

          <div className="form-group">
            <label htmlFor="st_dom_id">Studentski dom *</label>
            {loadingStDoms ? (
              <div className="loading-text">Učitavanje domova...</div>
            ) : (
              <select
                id="st_dom_id"
                name="st_dom_id"
                value={formData.st_dom_id}
                onChange={handleInputChange}
                required
                disabled={loading}
              >
                <option value="">Odaberite studentski dom</option>
                {stDoms.map(stDom => (
                  <option key={stDom.id} value={stDom.id}>
                    {stDom.ime} - {stDom.address}
                  </option>
                ))}
              </select>
            )}
          </div>

          <div className="form-group">
            <label htmlFor="krevetnost">Broj kreveta *</label>
            <input
              type="number"
              id="krevetnost"
              name="krevetnost"
              value={formData.krevetnost}
              onChange={handleInputChange}
              required
              min="1"
              disabled={loading}
              placeholder="Unesite broj kreveta"
            />
          </div>

          <div className="form-group">
            <label>Luksuzi (opciono)</label>
            <div className="checkbox-group">
              {availableLuksuzi.map(luksuz => (
                <label key={luksuz.value} className="checkbox-label">
                  <input
                    type="checkbox"
                    value={luksuz.value}
                    checked={formData.luksuzi.includes(luksuz.value)}
                    onChange={handleLuksuziChange}
                    disabled={loading}
                  />
                  <span>{luksuz.label}</span>
                </label>
              ))}
            </div>
          </div>

          <div className="modal-buttons">
            <button 
              type="button" 
              onClick={handleClose}
              className="cancel-button"
              disabled={loading}
            >
              Odustani
            </button>
            <button 
              type="submit" 
              className="submit-button"
              disabled={loading || loadingStDoms}
            >
              {loading ? 'Kreiranje...' : 'Kreiraj sobu'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AddRoomModal;

