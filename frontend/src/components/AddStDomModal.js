import React, { useState } from 'react';
import { stDomService } from '../services/stDomService';
import './AddStDomModal.css';

const AddStDomModal = ({ isOpen, onClose, onSuccess }) => {
  const [formData, setFormData] = useState({
    ime: '',
    address: '',
    telephone_number: '',
    email: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

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
      await stDomService.createStDom(formData);
      onSuccess();
      onClose();
      // Reset form
      setFormData({
        ime: '',
        address: '',
        telephone_number: '',
        email: ''
      });
    } catch (err) {
      setError(err.message || 'Greška pri kreiranju studentskog doma');
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
        ime: '',
        address: '',
        telephone_number: '',
        email: ''
      });
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal-overlay">
      <div className="modal add-st-dom-modal">
        <div className="modal-header">
          <h3>Dodaj novi studentski dom</h3>
          <button 
            className="close-button" 
            onClick={handleClose}
            disabled={loading}
          >
            ×
          </button>
        </div>

        <form onSubmit={handleSubmit} className="st-dom-form">
          {error && (
            <div className="error-message">
              {error}
            </div>
          )}

          <div className="form-group">
            <label htmlFor="ime">Naziv doma *</label>
            <input
              type="text"
              id="ime"
              name="ime"
              value={formData.ime}
              onChange={handleInputChange}
              required
              disabled={loading}
              placeholder="Unesite naziv studentskog doma"
            />
          </div>

          <div className="form-group">
            <label htmlFor="address">Adresa *</label>
            <input
              type="text"
              id="address"
              name="address"
              value={formData.address}
              onChange={handleInputChange}
              required
              disabled={loading}
              placeholder="Unesite adresu doma"
            />
          </div>

          <div className="form-group">
            <label htmlFor="telephone_number">Broj telefona *</label>
            <input
              type="tel"
              id="telephone_number"
              name="telephone_number"
              value={formData.telephone_number}
              onChange={handleInputChange}
              required
              disabled={loading}
              placeholder="Unesite broj telefona"
            />
          </div>

          <div className="form-group">
            <label htmlFor="email">Email *</label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={handleInputChange}
              required
              disabled={loading}
              placeholder="Unesite email adresu"
            />
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
              disabled={loading}
            >
              {loading ? 'Kreiranje...' : 'Kreiraj dom'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AddStDomModal;
