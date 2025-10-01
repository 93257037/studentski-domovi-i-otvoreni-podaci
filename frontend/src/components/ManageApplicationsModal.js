import React, { useState, useEffect } from 'react';
import { stDomService } from '../services/stDomService';
import './ManageApplicationsModal.css';

const ManageApplicationsModal = ({ isOpen, onClose, onSuccess }) => {
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [actionLoading, setActionLoading] = useState(null); // Tracks which application is being processed
  const [academicYear, setAcademicYear] = useState('2024/2025'); // Default academic year

  useEffect(() => {
    if (isOpen) {
      fetchApplications();
    }
  }, [isOpen]);

  const fetchApplications = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await stDomService.getAllAplikacije();
      console.log('=== ALL APPLICATIONS ===');
      console.log('Response:', response);
      console.log('========================');
      setApplications(response.aplikacije || []);
    } catch (err) {
      setError('Greška pri učitavanju aplikacija: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async (aplikacija) => {
    if (!academicYear.trim()) {
      alert('Molimo unesite akademsku godinu!');
      return;
    }

    const confirmed = window.confirm(
      `Jeste li sigurni da želite odobriti aplikaciju za:\n` +
      `Student: ${aplikacija.broj_indexa}\n` +
      `Prosek: ${aplikacija.prosek}\n` +
      `Akademska godina: ${academicYear}`
    );

    if (!confirmed) return;

    setActionLoading(aplikacija.id);
    setError('');

    try {
      await stDomService.approveAplikacija(aplikacija.id, academicYear);
      alert('Aplikacija uspješno odobrena!');
      // Refresh the list
      await fetchApplications();
      onSuccess();
    } catch (err) {
      setError('Greška pri odobravanju aplikacije: ' + err.message);
      alert('Greška pri odobravanju aplikacije: ' + err.message);
    } finally {
      setActionLoading(null);
    }
  };

  const handleReject = async (aplikacija) => {
    const confirmed = window.confirm(
      `Jeste li sigurni da želite odbiti aplikaciju za:\n` +
      `Student: ${aplikacija.broj_indexa}\n` +
      `Prosek: ${aplikacija.prosek}\n\n` +
      `Ova akcija će trajno obrisati aplikaciju.`
    );

    if (!confirmed) return;

    setActionLoading(aplikacija.id);
    setError('');

    try {
      await stDomService.rejectAplikacija(aplikacija.id);
      alert('Aplikacija uspješno odbijena!');
      // Refresh the list
      await fetchApplications();
      onSuccess();
    } catch (err) {
      setError('Greška pri odbijanju aplikacije: ' + err.message);
      alert('Greška pri odbijanju aplikacije: ' + err.message);
    } finally {
      setActionLoading(null);
    }
  };

  const handleClose = () => {
    if (!actionLoading) {
      onClose();
      setError('');
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal-overlay">
      <div className="modal manage-applications-modal">
        <div className="modal-header">
          <h3>Upravljanje aplikacijama</h3>
          <button 
            className="close-button" 
            onClick={handleClose}
            disabled={actionLoading !== null}
          >
            ×
          </button>
        </div>

        {error && (
          <div className="error-message">
            {error}
          </div>
        )}

        <div className="academic-year-input">
          <label htmlFor="academic_year">Akademska godina za odobrenje:</label>
          <input
            type="text"
            id="academic_year"
            value={academicYear}
            onChange={(e) => setAcademicYear(e.target.value)}
            placeholder="2024/2025"
            disabled={actionLoading !== null}
          />
        </div>

        {loading ? (
          <div className="loading-text">Učitavanje aplikacija...</div>
        ) : (
          <div className="applications-list">
            {applications.length === 0 ? (
              <div className="no-data">Nema aplikacija za prikaz</div>
            ) : (
              applications.map(app => (
                <div key={app.id} className={`application-card ${!app.is_active ? 'inactive' : ''}`}>
                  <div className="application-info">
                    <div className="info-row">
                      <span className="label">Broj indeksa:</span>
                      <span className="value">{app.broj_indexa}</span>
                    </div>
                    <div className="info-row">
                      <span className="label">Prosek:</span>
                      <span className="value">{app.prosek}</span>
                    </div>
                    <div className="info-row">
                      <span className="label">ID Sobe:</span>
                      <span className="value">{app.soba_id}</span>
                    </div>
                    <div className="info-row">
                      <span className="label">Status:</span>
                      <span className={`status ${app.is_active ? 'active' : 'inactive'}`}>
                        {app.is_active ? 'Aktivna' : 'Neaktivna'}
                      </span>
                    </div>
                    {app.created_at && (
                      <div className="info-row">
                        <span className="label">Kreirana:</span>
                        <span className="value">
                          {new Date(app.created_at).toLocaleDateString('hr-HR')}
                        </span>
                      </div>
                    )}
                  </div>
                  <div className="application-actions">
                    <button
                      onClick={() => handleApprove(app)}
                      className="approve-button"
                      disabled={actionLoading !== null || !app.is_active}
                      title={!app.is_active ? 'Aplikacija nije aktivna' : 'Odobri aplikaciju'}
                    >
                      {actionLoading === app.id ? 'Obrađujem...' : 'Odobri'}
                    </button>
                    <button
                      onClick={() => handleReject(app)}
                      className="reject-button"
                      disabled={actionLoading !== null}
                    >
                      {actionLoading === app.id ? 'Obrađujem...' : 'Odbij'}
                    </button>
                  </div>
                </div>
              ))
            )}
          </div>
        )}

        <div className="modal-buttons">
          <button 
            type="button" 
            onClick={handleClose}
            className="close-modal-button"
            disabled={actionLoading !== null}
          >
            Zatvori
          </button>
        </div>
      </div>
    </div>
  );
};

export default ManageApplicationsModal;

