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

  // Statistics states
  const [statistics, setStatistics] = useState({
    topFull: [],
    topEmpty: [],
    mostApplications: null,
    highestProsek: null
  });
  const [statisticsLoading, setStatisticsLoading] = useState(false);
  const [statisticsError, setStatisticsError] = useState('');

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

  // Statistics functions
  const fetchStatistics = async () => {
    setStatisticsLoading(true);
    setStatisticsError('');
    
    try {
      const [topFullResponse, topEmptyResponse, mostApplicationsResponse, highestProsekResponse] = await Promise.all([
        openDataService.getTopFullStDoms(),
        openDataService.getTopEmptyStDoms(),
        openDataService.getStDomWithMostApplications(),
        openDataService.getStDomWithHighestAverageProsek()
      ]);

      setStatistics({
        topFull: topFullResponse.data || [],
        topEmpty: topEmptyResponse.data || [],
        mostApplications: mostApplicationsResponse.data || null,
        highestProsek: highestProsekResponse.data || null
      });
    } catch (error) {
      setStatisticsError('Greška pri učitavanju statistika: ' + error.message);
    } finally {
      setStatisticsLoading(false);
    }
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

        {/* Statistics Section */}
        <div className="statistics-section">
          <div className="statistics-header">
            <h2>Statistike studentskih domova</h2>
            <button 
              onClick={fetchStatistics} 
              className="refresh-statistics-button"
              disabled={statisticsLoading}
            >
              {statisticsLoading ? 'Učitavanje...' : 'Osvježi statistike'}
            </button>
          </div>
          
          {statisticsError && <div className="error-message">{statisticsError}</div>}
          
          {statisticsLoading ? (
            <div className="loading-message">Učitavanje statistika...</div>
          ) : (
            <div className="statistics-grid">
              {/* Most Populated */}
              <div className="statistics-card">
                <h3>Najnaseljeniji domovi</h3>
                {statistics.topFull.length > 0 ? (
                  <div className="statistics-list">
                    {statistics.topFull.map((stDom, index) => (
                      <div key={stDom._id || index} className="statistics-item">
                        <div className="rank">#{index + 1}</div>
                        <div className="info">
                          <h4>{stDom.ime}</h4>
                          <p>Broj stanara: {stDom.broj_stanara || 'N/A'}</p>
                          <p>Adresa: {stDom.address}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="no-data">Nema dostupnih podataka</p>
                )}
              </div>

              {/* Most Unpopulated */}
              <div className="statistics-card">
                <h3>Najmanje naseljeni domovi</h3>
                {statistics.topEmpty.length > 0 ? (
                  <div className="statistics-list">
                    {statistics.topEmpty.map((stDom, index) => (
                      <div key={stDom._id || index} className="statistics-item">
                        <div className="rank">#{index + 1}</div>
                        <div className="info">
                          <h4>{stDom.ime}</h4>
                          <p>Broj stanara: {stDom.broj_stanara || 'N/A'}</p>
                          <p>Adresa: {stDom.address}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="no-data">Nema dostupnih podataka</p>
                )}
              </div>

              {/* Most Applied For */}
              <div className="statistics-card">
                <h3>Dom s najviše prijava</h3>
                {statistics.mostApplications ? (
                  <div className="statistics-single">
                    <h4>{statistics.mostApplications.ime}</h4>
                    <p>Broj prijava: {statistics.mostApplications.broj_prijava || 'N/A'}</p>
                    <p>Adresa: {statistics.mostApplications.address}</p>
                  </div>
                ) : (
                  <p className="no-data">Nema dostupnih podataka</p>
                )}
              </div>

              {/* Highest Prosek */}
              <div className="statistics-card">
                <h3>Dom s najvišim prosjekom</h3>
                {statistics.highestProsek ? (
                  <div className="statistics-single">
                    <h4>{statistics.highestProsek.ime}</h4>
                    <p>Prosječni prosek: {statistics.highestProsek.prosječni_prosek ? statistics.highestProsek.prosječni_prosek.toFixed(2) : 'N/A'}</p>
                    <p>Adresa: {statistics.highestProsek.address}</p>
                  </div>
                ) : (
                  <p className="no-data">Nema dostupnih podataka</p>
                )}
              </div>
            </div>
          )}
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
