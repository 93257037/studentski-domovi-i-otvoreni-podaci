import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { stDomService } from '../services/stDomService';
import './ApplyToRoomButton.css';

const ApplyToRoomButton = ({ room, stDom, onSuccess }) => {
  const { isAuthenticated } = useAuth();
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({
    broj_indexa: '',
    prosek: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // Don't render if user is not authenticated
  if (!isAuthenticated) {
    return null;
  }

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
        soba_id: room.id
      };

      await stDomService.createAplikacija(aplikacijaData);
      
      // Reset form and close
      setFormData({ broj_indexa: '', prosek: '' });
      setShowForm(false);
      
      // Call success callback if provided
      if (onSuccess) {
        onSuccess();
      }
      
      // Show success message
      alert('Aplikacija je uspješno poslana!');
    } catch (err) {
      setError(err.message || 'Greška pri kreiranju aplikacije');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    setShowForm(false);
    setFormData({ broj_indexa: '', prosek: '' });
    setError('');
  };

  if (showForm) {
    return (
      <div className="apply-form-container" onClick={(e) => e.stopPropagation()}>
        <div className="apply-form-header">
          <h4>Apliciraj za sobu</h4>
          {stDom && <p className="room-info">Dom: {stDom.ime}</p>}
          <p className="room-info">Kapacitet: {room.krevetnost} kreveta</p>
        </div>

        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="apply-form">
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

          <div className="form-buttons">
            <button 
              type="button" 
              onClick={handleCancel}
              className="cancel-button"
              disabled={loading}
            >
              Otkaži
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
    );
  }

  return (
    <button 
      className="apply-to-room-button"
      onClick={(e) => {
        e.stopPropagation();
        setShowForm(true);
      }}
    >
      Apliciraj za sobu
    </button>
  );
};

export default ApplyToRoomButton;
