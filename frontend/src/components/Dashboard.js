import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import AddStDomModal from './AddStDomModal';
import { openDataService } from '../services/openDataService';
import './Dashboard.css';

const Dashboard = () => {
  const { user, logout, deleteAccount } = useAuth();
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deleteLoading, setDeleteLoading] = useState(false);
  const [showAddStDomModal, setShowAddStDomModal] = useState(false);
  
  // Search states
  const [imeSearch, setImeSearch] = useState('');
  const [addressSearch, setAddressSearch] = useState('');
  const [imeResults, setImeResults] = useState([]);
  const [addressResults, setAddressResults] = useState([]);
  const [imeLoading, setImeLoading] = useState(false);
  const [addressLoading, setAddressLoading] = useState(false);
  const [imeError, setImeError] = useState('');
  const [addressError, setAddressError] = useState('');

  const handleDeleteAccount = async () => {
    setDeleteLoading(true);
    const result = await deleteAccount();
    
    if (!result.success) {
      alert('Greška pri brisanju računa: ' + result.error);
    }
    
    setDeleteLoading(false);
    setShowDeleteModal(false);
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('hr-HR', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const handleStDomSuccess = () => {
    // You can add logic here to refresh the list of dormitories
    // For now, we'll just show a success message
    alert('Studentski dom je uspješno kreiran!');
  };

  // Search functions
  const handleImeSearch = async (e) => {
    e.preventDefault();
    if (!imeSearch.trim()) {
      setImeResults([]);
      setImeError('');
      return;
    }

    setImeLoading(true);
    setImeError('');
    
    try {
      const response = await openDataService.searchStDomsByIme(imeSearch);
      setImeResults(response.data || []);
    } catch (error) {
      setImeError('Greška pri pretraživanju: ' + error.message);
      setImeResults([]);
    } finally {
      setImeLoading(false);
    }
  };

  const handleAddressSearch = async (e) => {
    e.preventDefault();
    if (!addressSearch.trim()) {
      setAddressResults([]);
      setAddressError('');
      return;
    }

    setAddressLoading(true);
    setAddressError('');
    
    try {
      const response = await openDataService.searchStDomsByAddress(addressSearch);
      setAddressResults(response.data || []);
    } catch (error) {
      setAddressError('Greška pri pretraživanju: ' + error.message);
      setAddressResults([]);
    } finally {
      setAddressLoading(false);
    }
  };

  const clearImeSearch = () => {
    setImeSearch('');
    setImeResults([]);
    setImeError('');
  };

  const clearAddressSearch = () => {
    setAddressSearch('');
    setAddressResults([]);
    setAddressError('');
  };

  return (
    <div className="dashboard-container">
      <div className="dashboard-header">
        <h1>Dobrodošli, {user?.first_name}!</h1>
        <button onClick={logout} className="logout-button">
          Odjavi se
        </button>
      </div>

      <div className="dashboard-content">
        <div className="profile-card">
          <h2>Profil korisnika</h2>
          <div className="profile-info">
            <div className="info-row">
              <span className="label">Korisničko ime:</span>
              <span className="value">{user?.username}</span>
            </div>
            <div className="info-row">
              <span className="label">Email:</span>
              <span className="value">{user?.email}</span>
            </div>
            <div className="info-row">
              <span className="label">Ime:</span>
              <span className="value">{user?.first_name}</span>
            </div>
            <div className="info-row">
              <span className="label">Prezime:</span>
              <span className="value">{user?.last_name}</span>
            </div>
            <div className="info-row">
              <span className="label">Uloga:</span>
              <span className="value">{user?.role}</span>
            </div>
            <div className="info-row">
              <span className="label">Registriran:</span>
              <span className="value">{user?.created_at ? formatDate(user.created_at) : 'N/A'}</span>
            </div>
          </div>
        </div>

        <div className="actions-card">
          <h2>Akcije</h2>
          <div className="action-buttons">
            {user?.role === 'admin' && (
              <button 
                onClick={() => setShowAddStDomModal(true)}
                className="add-st-dom-button"
              >
                Dodaj studentski dom
              </button>
            )}
            <button 
              onClick={() => setShowDeleteModal(true)}
              className="delete-account-button"
            >
              Obriši račun
            </button>
          </div>
        </div>

        {/* Search Section */}
        <div className="search-section">
          <h2>Pretraživanje studentskih domova</h2>
          
          {/* Ime Search */}
          <div className="search-card">
            <h3>Pretraži po imenu</h3>
            <form onSubmit={handleImeSearch} className="search-form">
              <div className="search-input-group">
                <input
                  type="text"
                  value={imeSearch}
                  onChange={(e) => setImeSearch(e.target.value)}
                  placeholder="Unesite ime studentskog doma..."
                  className="search-input"
                />
                <button type="submit" className="search-button" disabled={imeLoading}>
                  {imeLoading ? 'Pretraživanje...' : 'Pretraži'}
                </button>
                {imeSearch && (
                  <button type="button" onClick={clearImeSearch} className="clear-button">
                    Obriši
                  </button>
                )}
              </div>
            </form>
            
            {imeError && <div className="error-message">{imeError}</div>}
            
            {imeResults.length > 0 && (
              <div className="search-results">
                <h4>Rezultati pretraživanja ({imeResults.length})</h4>
                <div className="results-list">
                  {imeResults.map((stDom) => (
                    <div key={stDom._id} className="result-item">
                      <h5>{stDom.ime}</h5>
                      <p><strong>Adresa:</strong> {stDom.address}</p>
                      {stDom.description && <p><strong>Opis:</strong> {stDom.description}</p>}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Address Search */}
          <div className="search-card">
            <h3>Pretraži po adresi</h3>
            <form onSubmit={handleAddressSearch} className="search-form">
              <div className="search-input-group">
                <input
                  type="text"
                  value={addressSearch}
                  onChange={(e) => setAddressSearch(e.target.value)}
                  placeholder="Unesite adresu studentskog doma..."
                  className="search-input"
                />
                <button type="submit" className="search-button" disabled={addressLoading}>
                  {addressLoading ? 'Pretraživanje...' : 'Pretraži'}
                </button>
                {addressSearch && (
                  <button type="button" onClick={clearAddressSearch} className="clear-button">
                    Obriši
                  </button>
                )}
              </div>
            </form>
            
            {addressError && <div className="error-message">{addressError}</div>}
            
            {addressResults.length > 0 && (
              <div className="search-results">
                <h4>Rezultati pretraživanja ({addressResults.length})</h4>
                <div className="results-list">
                  {addressResults.map((stDom) => (
                    <div key={stDom._id} className="result-item">
                      <h5>{stDom.ime}</h5>
                      <p><strong>Adresa:</strong> {stDom.address}</p>
                      {stDom.description && <p><strong>Opis:</strong> {stDom.description}</p>}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {showDeleteModal && (
        <div className="modal-overlay">
          <div className="modal">
            <h3>Potvrdite brisanje računa</h3>
            <p>
              Jeste li sigurni da želite obrisati svoj račun? Ova akcija se ne može poništiti.
            </p>
            <div className="modal-buttons">
              <button 
                onClick={() => setShowDeleteModal(false)}
                className="cancel-button"
                disabled={deleteLoading}
              >
                Odustani
              </button>
              <button 
                onClick={handleDeleteAccount}
                className="confirm-delete-button"
                disabled={deleteLoading}
              >
                {deleteLoading ? 'Brisanje...' : 'Obriši račun'}
              </button>
            </div>
          </div>
        </div>
      )}

      <AddStDomModal
        isOpen={showAddStDomModal}
        onClose={() => setShowAddStDomModal(false)}
        onSuccess={handleStDomSuccess}
      />
    </div>
  );
};

export default Dashboard;
