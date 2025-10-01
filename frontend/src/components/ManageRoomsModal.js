import React, { useState, useEffect } from 'react';
import { stDomService } from '../services/stDomService';
import './ManageRoomsModal.css';

const ManageRoomsModal = ({ isOpen, onClose, onSuccess }) => {
  const [acceptedApps, setAcceptedApps] = useState([]);
  const [payments, setPayments] = useState({});
  const [users, setUsers] = useState({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [actionLoading, setActionLoading] = useState(null);

  // Search filters
  const [searchIndex, setSearchIndex] = useState('');
  const [searchRoom, setSearchRoom] = useState('');
  const [searchUsername, setSearchUsername] = useState('');
  const [searchStDom, setSearchStDom] = useState('');
  const [paymentStatus, setPaymentStatus] = useState('');

  useEffect(() => {
    if (isOpen) {
      fetchData();
    }
  }, [isOpen]);

  const fetchData = async () => {
    setLoading(true);
    setError('');
    try {
      // Fetch all accepted applications
      const appsResponse = await stDomService.getAllPrihvaceneAplikacije();
      const apps = appsResponse.prihvacene_aplikacije || [];
      setAcceptedApps(apps);

      // Fetch payments and user info for each application
      const paymentsData = {};
      const usersData = {};
      
      for (const app of apps) {
        try {
          const paymentResponse = await stDomService.getPaymentsByUser(app.user_id);
          paymentsData[app.user_id] = paymentResponse.payments || [];
        } catch (err) {
          console.log(`No payments for user ${app.user_id}`);
        }

        // Store user info - we'll display the user_id for now
        // TODO: Add endpoint to SSO service to fetch user info by ID for admins
        usersData[app.user_id] = {
          username: `User-${app.user_id.slice(-8)}`, // Last 8 chars of user ID as identifier
          fullId: app.user_id,
          index: app.broj_indexa
        };
      }
      
      setPayments(paymentsData);
      setUsers(usersData);
    } catch (err) {
      setError('Greška pri učitavanju podataka: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleMarkAsPaid = async (paymentId) => {
    setActionLoading(paymentId);
    setError('');
    try {
      await stDomService.markPaymentAsPaid(paymentId);
      alert('Plaćanje je označeno kao plaćeno!');
      await fetchData();
      onSuccess();
    } catch (err) {
      setError('Greška pri označavanju plaćanja: ' + err.message);
      alert('Greška: ' + err.message);
    } finally {
      setActionLoading(null);
    }
  };

  const handleMarkAsUnpaid = async (paymentId) => {
    setActionLoading(paymentId);
    setError('');
    try {
      await stDomService.markPaymentAsUnpaid(paymentId);
      alert('Plaćanje je označeno kao neplaćeno!');
      await fetchData();
      onSuccess();
    } catch (err) {
      setError('Greška pri označavanju plaćanja: ' + err.message);
      alert('Greška: ' + err.message);
    } finally {
      setActionLoading(null);
    }
  };

  const handleEvict = async (app) => {
    const reason = window.prompt(
      `Unesite razlog za iseljenje studenta:\n` +
      `Student: ${app.broj_indexa}\n` +
      `Prosek: ${app.prosek}`
    );

    if (!reason || !reason.trim()) {
      alert('Razlog za iseljenje je obavezan!');
      return;
    }

    const confirmed = window.confirm(
      `Jeste li sigurni da želite iseliti studenta?\n` +
      `Student: ${app.broj_indexa}\n` +
      `Razlog: ${reason}`
    );

    if (!confirmed) return;

    setActionLoading(app.id);
    setError('');

    try {
      await stDomService.evictStudent(app.user_id, reason);
      alert('Student je uspješno iseljen!');
      await fetchData();
      onSuccess();
    } catch (err) {
      setError('Greška pri iseljenju studenta: ' + err.message);
      alert('Greška: ' + err.message);
    } finally {
      setActionLoading(null);
    }
  };

  const getPaymentForApp = (app) => {
    const userPayments = payments[app.user_id] || [];
    // Return latest payment or most recent one
    return userPayments.length > 0 ? userPayments[0] : null;
  };

  const getFilteredApps = () => {
    return acceptedApps.filter(app => {
      if (searchIndex && !app.broj_indexa.toLowerCase().includes(searchIndex.toLowerCase())) {
        return false;
      }
      if (searchRoom && !app.soba_id.toLowerCase().includes(searchRoom.toLowerCase())) {
        return false;
      }
      if (searchUsername) {
        const user = users[app.user_id];
        const searchLower = searchUsername.toLowerCase();
        if (!user || (!user.fullId.toLowerCase().includes(searchLower) && 
                      !user.index.toLowerCase().includes(searchLower))) {
          return false;
        }
      }
      if (searchStDom) {
        // This would require fetching room details to get st_dom_id
        // For now, we'll skip this filter or you'd need to enhance the data structure
      }
      if (paymentStatus) {
        const payment = getPaymentForApp(app);
        if (!payment) return paymentStatus === 'none';
        if (payment.status !== paymentStatus) return false;
      }
      return true;
    });
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleDateString('hr-HR');
  };

  const getPaymentStatus = (app) => {
    const payment = getPaymentForApp(app);
    if (!payment) return { status: 'none', label: 'Nema plaćanja' };
    
    if (payment.status === 'paid') return { status: 'paid', label: 'Plaćeno' };
    if (payment.status === 'overdue') return { status: 'overdue', label: 'Prekoračeno' };
    return { status: 'pending', label: 'Neplaćeno' };
  };

  const handleClose = () => {
    if (!actionLoading) {
      onClose();
      setError('');
      // Reset filters
      setSearchIndex('');
      setSearchRoom('');
      setSearchUsername('');
      setSearchStDom('');
      setPaymentStatus('');
    }
  };

  if (!isOpen) return null;

  const filteredApps = getFilteredApps();

  return (
    <div className="modal-overlay">
      <div className="modal manage-rooms-modal">
        <div className="modal-header">
          <h3>Upravljanje sobama i plaćanjima</h3>
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

        {/* Search Filters */}
        <div className="search-filters">
          <h4>Pretraživanje</h4>
          <div className="filters-grid">
            <input
              type="text"
              placeholder="Pretraži po indeksu..."
              value={searchIndex}
              onChange={(e) => setSearchIndex(e.target.value)}
              disabled={actionLoading !== null}
            />
            <input
              type="text"
              placeholder="Pretraži po ID korisnika..."
              value={searchUsername}
              onChange={(e) => setSearchUsername(e.target.value)}
              disabled={actionLoading !== null}
            />
            <input
              type="text"
              placeholder="Pretraži po ID sobe..."
              value={searchRoom}
              onChange={(e) => setSearchRoom(e.target.value)}
              disabled={actionLoading !== null}
            />
            <select
              value={paymentStatus}
              onChange={(e) => setPaymentStatus(e.target.value)}
              disabled={actionLoading !== null}
            >
              <option value="">Svi statusi plaćanja</option>
              <option value="paid">Plaćeno</option>
              <option value="pending">Na čekanju</option>
              <option value="overdue">Prekoračeno</option>
              <option value="none">Bez plaćanja</option>
            </select>
          </div>
        </div>

        {loading ? (
          <div className="loading-text">Učitavanje podataka...</div>
        ) : (
          <div className="rooms-list">
            {filteredApps.length === 0 ? (
              <div className="no-data">Nema rezultata za prikaz</div>
            ) : (
              filteredApps.map(app => {
                const payment = getPaymentForApp(app);
                const user = users[app.user_id] || { username: 'N/A' };
                const paymentStatusInfo = getPaymentStatus(app);
                
                return (
                  <div key={app.id} className="room-card">
                    <div className="room-info">
                      <div className="info-section">
                        <h4>Informacije o studentu</h4>
                        <div className="info-row">
                          <span className="label">Broj indeksa:</span>
                          <span className="value">{app.broj_indexa}</span>
                        </div>
                        <div className="info-row">
                          <span className="label">ID Korisnika:</span>
                          <span className="value small">{user.fullId || app.user_id}</span>
                        </div>
                        <div className="info-row">
                          <span className="label">Prosek:</span>
                          <span className="value">{app.prosek}</span>
                        </div>
                        <div className="info-row">
                          <span className="label">Akademska godina:</span>
                          <span className="value">{app.academic_year}</span>
                        </div>
                      </div>

                      <div className="info-section">
                        <h4>Informacije o sobi</h4>
                        <div className="info-row">
                          <span className="label">ID Sobe:</span>
                          <span className="value small">{app.soba_id}</span>
                        </div>
                        <div className="info-row">
                          <span className="label">Status plaćanja:</span>
                          <span className={`payment-status ${paymentStatusInfo.status}`}>
                            {paymentStatusInfo.label}
                          </span>
                        </div>
                        {payment && (
                          <div className="info-row">
                            <span className="label">Plaćeno za period:</span>
                            <span className="value">{payment.payment_period}</span>
                          </div>
                        )}
                        <div className="info-row">
                          <span className="label">Kreirana:</span>
                          <span className="value">{formatDate(app.created_at)}</span>
                        </div>
                      </div>

                      {payment && (
                        <div className="info-section payment-section">
                          <h4>Informacije o plaćanju</h4>
                          <div className="info-row">
                            <span className="label">Iznos:</span>
                            <span className="value">{payment.amount} €</span>
                          </div>
                          <div className="info-row">
                            <span className="label">Period:</span>
                            <span className="value">{payment.payment_period}</span>
                          </div>
                          <div className="info-row">
                            <span className="label">Status:</span>
                            <span className={`payment-status ${payment.status}`}>
                              {payment.status === 'paid' && 'Plaćeno'}
                              {payment.status === 'pending' && 'Na čekanju'}
                              {payment.status === 'overdue' && 'Prekoračeno'}
                            </span>
                          </div>
                          <div className="info-row">
                            <span className="label">Rok:</span>
                            <span className="value">{formatDate(payment.due_date)}</span>
                          </div>
                          {payment.paid_at && (
                            <div className="info-row">
                              <span className="label">Plaćeno:</span>
                              <span className="value">{formatDate(payment.paid_at)}</span>
                            </div>
                          )}
                        </div>
                      )}
                    </div>

                    <div className="room-actions">
                      {payment && payment.status !== 'paid' && (
                        <button
                          onClick={() => handleMarkAsPaid(payment.id)}
                          className="mark-paid-button"
                          disabled={actionLoading !== null}
                        >
                          {actionLoading === payment.id ? 'Obrađujem...' : '✓ Označi plaćeno'}
                        </button>
                      )}
                      {payment && payment.status === 'paid' && (
                        <button
                          onClick={() => handleMarkAsUnpaid(payment.id)}
                          className="mark-unpaid-button"
                          disabled={actionLoading !== null}
                        >
                          {actionLoading === payment.id ? 'Obrađujem...' : '✗ Označi neplaćeno'}
                        </button>
                      )}
                      <button
                        onClick={() => handleEvict(app)}
                        className="evict-button"
                        disabled={actionLoading !== null}
                      >
                        {actionLoading === app.id ? 'Obrađujem...' : '🚪 Iseli studenta'}
                      </button>
                    </div>
                  </div>
                );
              })
            )}
          </div>
        )}

        <div className="modal-footer">
          <div className="results-count">
            Prikazano: {filteredApps.length} od {acceptedApps.length} soba
          </div>
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

export default ManageRoomsModal;

