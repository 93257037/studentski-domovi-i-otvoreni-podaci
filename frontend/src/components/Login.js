import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { Link, useNavigate } from 'react-router-dom';
import './Auth.css';

// komponenta za prijavu korisnika - forma sa email i lozinkom
const Login = () => {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  
  const { login } = useAuth();
  const navigate = useNavigate();

  // rukuje promenama u input poljima
  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  // rukuje slanjem forme za prijavu
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    const result = await login(formData.email, formData.password);
    
    if (result.success) {
      navigate('/dashboard');
    } else {
      setError(result.error);
    }
    
    setLoading(false);
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h2>Prijava</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="email">Email:</label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              required
              placeholder="Unesite vaš email"
            />
          </div>

          <div className="form-group">
            <label htmlFor="password">Lozinka:</label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
              placeholder="Unesite vašu lozinku"
            />
          </div>

          {error && <div className="error-message">{error}</div>}

          <button type="submit" disabled={loading} className="auth-button">
            {loading ? 'Prijavljivanje...' : 'Prijavi se'}
          </button>
        </form>

        <div className="auth-links">
          <p>
            Nemate račun? <Link to="/register">Registrirajte se</Link>
          </p>
        </div>

        <div className="open-data-section">
          <div className="divider">
            <span>ili</span>
          </div>
          <button 
            type="button" 
            onClick={() => navigate('/open-data')}
            className="open-data-access-button"
          >
            📊 Pristupite otvorenim podacima
          </button>
          <p className="open-data-description">
            Pregledajte javne statistike i podatke o studentskim domovima bez prijavljivanja
          </p>
        </div>
      </div>
    </div>
  );
};

export default Login;
