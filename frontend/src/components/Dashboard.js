import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import AddStDomModal from './AddStDomModal';
import AddRoomModal from './AddRoomModal';
import ApplyForRoomModal from './ApplyForRoomModal';
import ManageApplicationsModal from './ManageApplicationsModal';
import ManageRoomsModal from './ManageRoomsModal';
import MyRoomInfo from './MyRoomInfo';
import { openDataService } from '../services/openDataService';
import './Dashboard.css';

// glavna dashboard komponenta - prikazuje profil korisnika, akcije, statistike i pretragu
const Dashboard = () => {
  const { user, logout, deleteAccount } = useAuth();
  const navigate = useNavigate();
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deleteLoading, setDeleteLoading] = useState(false);
  const [showAddStDomModal, setShowAddStDomModal] = useState(false);
  const [showAddRoomModal, setShowAddRoomModal] = useState(false);
  const [showApplyForRoomModal, setShowApplyForRoomModal] = useState(false);
  const [showManageApplicationsModal, setShowManageApplicationsModal] = useState(false);
  const [showManageRoomsModal, setShowManageRoomsModal] = useState(false);
  
  const [imeSearch, setImeSearch] = useState('');
  const [addressSearch, setAddressSearch] = useState('');
  const [imeResults, setImeResults] = useState([]);
  const [addressResults, setAddressResults] = useState([]);
  const [imeLoading, setImeLoading] = useState(false);
  const [addressLoading, setAddressLoading] = useState(false);
  const [imeError, setImeError] = useState('');
  const [addressError, setAddressError] = useState('');

  const [statistics, setStatistics] = useState({
    topFull: [],
    topEmpty: [],
    mostApplications: null,
    highestProsek: null
  });
  const [statisticsLoading, setStatisticsLoading] = useState(false);
  const [statisticsError, setStatisticsError] = useState('');

  // rukuje brisanjem naloga korisnika
  const handleDeleteAccount = async () => {
    setDeleteLoading(true);
    const result = await deleteAccount();
    
    if (!result.success) {
      alert('Greška pri brisanju računa: ' + result.error);
    }
    
    setDeleteLoading(false);
    setShowDeleteModal(false);
  };

  // formatira datum u citljiv oblik
  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('hr-HR', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  // callback funkcija kada je dom uspesno kreiran
  const handleStDomSuccess = () => {
    alert('Studentski dom je uspješno kreiran!');
  };

  // callback funkcija kada je soba uspesno kreirana
  const handleRoomSuccess = () => {
    alert('Soba je uspješno kreirana!');
  };

  // callback funkcija kada je aplikacija uspesno poslata
  const handleApplicationSuccess = () => {
    alert('Aplikacija je uspješno poslata!');
    window.location.reload();
  };

  // callback funkcija kada se korisnik odjavi iz sobe
  const handleCheckout = () => {
    window.location.reload();
  };

  // callback funkcija nakon upravljanja aplikacijama
  const handleManageApplicationsSuccess = () => {
    // logika se moze dodati ovde po potrebi
  };

  // callback funkcija nakon upravljanja sobama
  const handleManageRoomsSuccess = () => {
    // logika se moze dodati ovde po potrebi
  };

  // pretrazuje domove po imenu
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

  // pretrazuje domove po adresi
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

  // brise pretragu po imenu
  const clearImeSearch = () => {
    setImeSearch('');
    setImeResults([]);
    setImeError('');
  };

  // brise pretragu po adresi
  const clearAddressSearch = () => {
    setAddressSearch('');
    setAddressResults([]);
    setAddressError('');
  };

  // navigira na stranicu studentskog doma
  const handleStDomClick = (stDomId) => {
    if (stDomId) {
      navigate(`/st-dom/${stDomId}`);
    }
  };

  // ucitava statistike domova sa servera
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

  // ucitava statistike kada se komponenta mount-uje
  useEffect(() => {
    fetchStatistics();
  }, []);

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
            {user?.role === 'admin' ? (
              <>
                <button 
                  onClick={() => setShowAddStDomModal(true)}
                  className="add-st-dom-button"
                >
                  Dodaj studentski dom
                </button>
                <button 
                  onClick={() => setShowAddRoomModal(true)}
                  className="add-room-button"
                >
                  Dodaj sobu
                </button>
                <button 
                  onClick={() => setShowManageApplicationsModal(true)}
                  className="manage-applications-button"
                >
                  Upravljaj aplikacijama
                </button>
                <button 
                  onClick={() => setShowManageRoomsModal(true)}
                  className="manage-rooms-button"
                >
                  Upravljaj sobama
                </button>
              </>
            ) : (
              <button 
                onClick={() => setShowApplyForRoomModal(true)}
                className="apply-for-room-button"
              >
                Apliciraj za sobu
              </button>
            )}
            <button 
              onClick={() => navigate('/academic-year-applications')}
              className="academic-year-button"
            >
              Prihvaćene aplikacije po godini
            </button>
            <button 
              onClick={() => navigate('/open-data')}
              className="open-data-button"
            >
              📊 Open Data Dashboard
            </button>
            <button 
              onClick={() => setShowDeleteModal(true)}
              className="delete-account-button"
            >
              Obriši račun
            </button>
          </div>
        </div>

        {/* My Room Info Section (for regular users) */}
        {user?.role === 'user' && (
          <MyRoomInfo onCheckout={handleCheckout} />
        )}
      </div>


      {/* Modals */}
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

      <AddRoomModal
        isOpen={showAddRoomModal}
        onClose={() => setShowAddRoomModal(false)}
        onSuccess={handleRoomSuccess}
      />

      <ApplyForRoomModal
        isOpen={showApplyForRoomModal}
        onClose={() => setShowApplyForRoomModal(false)}
        onSuccess={handleApplicationSuccess}
      />

      <ManageApplicationsModal
        isOpen={showManageApplicationsModal}
        onClose={() => setShowManageApplicationsModal(false)}
        onSuccess={handleManageApplicationsSuccess}
      />

      <ManageRoomsModal
        isOpen={showManageRoomsModal}
        onClose={() => setShowManageRoomsModal(false)}
        onSuccess={handleManageRoomsSuccess}
      />
    </div>
  );
};

export default Dashboard;
