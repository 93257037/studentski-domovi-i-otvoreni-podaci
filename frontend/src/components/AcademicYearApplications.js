import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { openDataService } from '../services/openDataService';
import './AcademicYearApplications.css';

const AcademicYearApplications = () => {
  const { user, token, logout } = useAuth();
  const navigate = useNavigate();
  const [academicYear, setAcademicYear] = useState('');
  const [applications, setApplications] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [hasSearched, setHasSearched] = useState(false);

  const handleSearch = async (e) => {
    e.preventDefault();
    if (!academicYear.trim()) {
      setError('Molimo unesite akademsku godinu');
      return;
    }

    // Validate academic year format (should be like 2024/2025)
    const academicYearPattern = /^\d{4}\/\d{4}$/;
    if (!academicYearPattern.test(academicYear.trim())) {
      setError('Akademska godina mora biti u formatu YYYY/YYYY (npr. 2024/2025)');
      return;
    }

    setLoading(true);
    setError('');
    setHasSearched(true);
    
    try {
      const response = await openDataService.getPrihvaceneAplikacijeForAcademicYear(academicYear.trim(), token);
      setApplications(response.data || []);
    } catch (error) {
      setError('Greška pri dohvaćanju podataka: ' + error.message);
      setApplications([]);
    } finally {
      setLoading(false);
    }
  };

  const clearSearch = () => {
    setAcademicYear('');
    setApplications([]);
    setError('');
    setHasSearched(false);
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


  return (
    <div className="academic-year-container">
      <div className="academic-year-header">
        <div className="header-content">
          <h1>Prihvaćene aplikacije po akademskoj godini</h1>
          <div className="header-actions">
            <button onClick={() => navigate('/dashboard')} className="back-button">
              Povratak na početnu
            </button>
            <button onClick={logout} className="logout-button">
              Odjavi se
            </button>
          </div>
        </div>
      </div>

      <div className="academic-year-content">
        <div className="search-card">
          <h2>Pretraži prihvaćene aplikacije</h2>
          <form onSubmit={handleSearch} className="search-form">
            <div className="search-input-group">
              <label htmlFor="academicYear">Akademska godina:</label>
              <input
                id="academicYear"
                type="text"
                value={academicYear}
                onChange={(e) => setAcademicYear(e.target.value)}
                placeholder="Unesite akademsku godinu (npr. 2024/2025)"
                className="search-input"
                disabled={loading}
              />
              <div className="button-group">
                <button type="submit" className="search-button" disabled={loading}>
                  {loading ? 'Pretraživanje...' : 'Pretraži'}
                </button>
                {academicYear && (
                  <button type="button" onClick={clearSearch} className="clear-button" disabled={loading}>
                    Obriši
                  </button>
                )}
              </div>
            </div>
          </form>
          
          <div className="format-hint">
            <p><strong>Format:</strong> Akademska godina mora biti u formatu YYYY/YYYY (npr. 2024/2025)</p>
          </div>
        </div>

        {error && <div className="error-message">{error}</div>}

        {hasSearched && !loading && (
          <div className="results-section">
            <div className="results-header">
              <h2>Rezultati pretrage</h2>
              <div className="results-info">
                <span className="academic-year-label">Akademska godina: <strong>{academicYear}</strong></span>
                <span className="count-label">Ukupno aplikacija: <strong>{applications.length}</strong></span>
              </div>
            </div>

            {applications.length > 0 ? (
              <div className="applications-grid">
                {applications.map((app, index) => (
                  <div key={app.id || index} className="application-card">
                    <div className="application-header">
                      <h3>Aplikacija #{index + 1}</h3>
                      <span className="application-id">ID: {app.id}</span>
                    </div>
                    
                    <div className="application-details">
                      <div className="detail-row">
                        <span className="label">Korisnik ID:</span>
                        <span className="value">{app.user_id}</span>
                      </div>
                      <div className="detail-row">
                        <span className="label">Broj indeksa:</span>
                        <span className="value">{app.broj_indexa}</span>
                      </div>
                      <div className="detail-row">
                        <span className="label">Prosek:</span>
                        <span className="value prosek">{app.prosek}</span>
                      </div>
                      <div className="detail-row">
                        <span className="label">Soba ID:</span>
                        <span className="value">{app.soba_id}</span>
                      </div>
                      <div className="detail-row">
                        <span className="label">Akademska godina:</span>
                        <span className="value">{app.academic_year}</span>
                      </div>
                      <div className="detail-row">
                        <span className="label">Kreirana:</span>
                        <span className="value">{app.created_at ? formatDate(app.created_at) : 'N/A'}</span>
                      </div>
                      <div className="detail-row">
                        <span className="label">Ažurirana:</span>
                        <span className="value">{app.updated_at ? formatDate(app.updated_at) : 'N/A'}</span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="no-results">
                <h3>Nema rezultata</h3>
                <p>Nisu pronađene prihvaćene aplikacije za akademsku godinu "{academicYear}".</p>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default AcademicYearApplications;
