import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import AddStDomModal from './AddStDomModal';
import './Dashboard.css';

const Dashboard = () => {
  const { user, logout, deleteAccount } = useAuth();
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deleteLoading, setDeleteLoading] = useState(false);
  const [showAddStDomModal, setShowAddStDomModal] = useState(false);

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
