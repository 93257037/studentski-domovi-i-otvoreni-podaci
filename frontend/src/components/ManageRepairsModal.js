import React, { useState, useEffect } from 'react';
import { stDomService } from '../services/stDomService';
import './ManageRepairsModal.css';

const ManageRepairsModal = ({ show, onClose }) => {
  const [repairs, setRepairs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [filter, setFilter] = useState('all'); // all, scheduled, in_progress, completed, cancelled

  useEffect(() => {
    if (show) {
      loadRepairs();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [show, filter]);

  const loadRepairs = async () => {
    setLoading(true);
    setError('');
    try {
      let response;
      if (filter === 'all') {
        response = await stDomService.getAllRepairs();
      } else {
        response = await stDomService.getRepairsByStatus(filter);
      }
      // The API returns { repairs: [...], count: number }
      setRepairs(response.repairs || []);
    } catch (error) {
      setError('Greška pri učitavanju popravki: ' + (error.error || error.message));
      setRepairs([]);
    } finally {
      setLoading(false);
    }
  };

  const updateRepairStatus = async (repairId, newStatus) => {
    try {
      await stDomService.updateRepair(repairId, { status: newStatus });
      await loadRepairs(); // Reload the list
    } catch (error) {
      setError('Greška pri ažuriranju statusa: ' + (error.error || error.message));
    }
  };

  const deleteRepair = async (repairId) => {
    if (!window.confirm('Da li ste sigurni da želite da obrišete ovu popravku?')) {
      return;
    }
    try {
      await stDomService.deleteRepair(repairId);
      await loadRepairs(); // Reload the list
    } catch (error) {
      setError('Greška pri brisanju popravke: ' + (error.error || error.message));
    }
  };

  const getStatusBadgeClass = (status) => {
    switch (status) {
      case 'scheduled':
        return 'status-scheduled';
      case 'in_progress':
        return 'status-in-progress';
      case 'completed':
        return 'status-completed';
      case 'cancelled':
        return 'status-cancelled';
      default:
        return '';
    }
  };

  const getStatusLabel = (status) => {
    switch (status) {
      case 'scheduled':
        return 'Zakazano';
      case 'in_progress':
        return 'U Toku';
      case 'completed':
        return 'Završeno';
      case 'cancelled':
        return 'Otkazano';
      default:
        return status;
    }
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    return date.toLocaleDateString('sr-RS', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (!show) return null;

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content manage-repairs-modal" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>Upravljanje Popravkama</h2>
          <button className="close-button" onClick={onClose}>×</button>
        </div>

        {error && <div className="error-message">{error}</div>}

        <div className="filter-section">
          <label>Filter:</label>
          <select value={filter} onChange={(e) => setFilter(e.target.value)}>
            <option value="all">Sve</option>
            <option value="scheduled">Zakazano</option>
            <option value="in_progress">U Toku</option>
            <option value="completed">Završeno</option>
            <option value="cancelled">Otkazano</option>
          </select>
        </div>

        {loading ? (
          <div className="loading-message">Učitavanje...</div>
        ) : repairs.length === 0 ? (
          <div className="no-repairs-message">Nema popravki za prikaz.</div>
        ) : (
          <div className="repairs-table-container">
            <table className="repairs-table">
              <thead>
                <tr>
                  <th>ID Sobe</th>
                  <th>Opis</th>
                  <th>Predviđeni Završetak</th>
                  <th>Status</th>
                  <th>Kreirano</th>
                  <th>Akcije</th>
                </tr>
              </thead>
              <tbody>
                {repairs.map((repair) => (
                  <tr key={repair.id}>
                    <td><code>{repair.room_id}</code></td>
                    <td className="description-cell">{repair.description}</td>
                    <td>{formatDate(repair.estimated_completion_date)}</td>
                    <td>
                      <span className={`status-badge ${getStatusBadgeClass(repair.status)}`}>
                        {getStatusLabel(repair.status)}
                      </span>
                    </td>
                    <td>{formatDate(repair.created_at)}</td>
                    <td>
                      <div className="action-buttons">
                        {repair.status === 'scheduled' && (
                          <button
                            className="btn-action btn-start"
                            onClick={() => updateRepairStatus(repair.id, 'in_progress')}
                            title="Započni popravku"
                          >
                            ▶
                          </button>
                        )}
                        {(repair.status === 'scheduled' || repair.status === 'in_progress') && (
                          <button
                            className="btn-action btn-complete"
                            onClick={() => updateRepairStatus(repair.id, 'completed')}
                            title="Označi kao završeno"
                          >
                            ✓
                          </button>
                        )}
                        {(repair.status === 'scheduled' || repair.status === 'in_progress') && (
                          <button
                            className="btn-action btn-cancel"
                            onClick={() => updateRepairStatus(repair.id, 'cancelled')}
                            title="Otkaži popravku"
                          >
                            ✕
                          </button>
                        )}
                        <button
                          className="btn-action btn-delete"
                          onClick={() => deleteRepair(repair.id)}
                          title="Obriši popravku"
                        >
                          🗑
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};

export default ManageRepairsModal;

