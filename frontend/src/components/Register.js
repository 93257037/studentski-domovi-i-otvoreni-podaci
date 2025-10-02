import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { Link, useNavigate } from 'react-router-dom';
import './Auth.css';

// komponenta za registraciju novog korisnika - forma sa svim potrebnim podacima
const Register = () => {
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    first_name: '',
    last_name: '',
  });
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);
  
  const { register } = useAuth();
  const navigate = useNavigate();

  // rukuje promenama u input poljima
  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  // rukuje slanjem forme za registraciju sa validacijom podataka
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');

    if (formData.password !== formData.confirmPassword) {
      setError('Lozinke se ne poklapaju');
      return;
    }

    if (formData.password.length < 6) {
      setError('Lozinka mora imati najmanje 6 karaktera');
      return;
    }

    if (formData.username.length < 3 || formData.username.length > 20) {
      setError('Korisničko ime mora imati između 3 i 20 karaktera');
      return;
    }

    setLoading(true);

    const { confirmPassword, ...registerData } = formData;
    const result = await register(registerData);
    
    if (result.success) {
      setSuccess('Registracija uspješna! Možete se sada prijaviti.');
      setTimeout(() => {
        navigate('/login');
      }, 2000);
    } else {
      setError(result.error);
    }
    
    setLoading(false);
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h2>Registracija</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="username">Korisničko ime:</label>
            <input
              type="text"
              id="username"
              name="username"
              value={formData.username}
              onChange={handleChange}
              required
              minLength="3"
              maxLength="20"
              placeholder="Unesite korisničko ime"
            />
          </div>

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
            <label htmlFor="first_name">Ime:</label>
            <input
              type="text"
              id="first_name"
              name="first_name"
              value={formData.first_name}
              onChange={handleChange}
              required
              placeholder="Unesite vaše ime"
            />
          </div>

          <div className="form-group">
            <label htmlFor="last_name">Prezime:</label>
            <input
              type="text"
              id="last_name"
              name="last_name"
              value={formData.last_name}
              onChange={handleChange}
              required
              placeholder="Unesite vaše prezime"
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
              minLength="6"
              placeholder="Unesite lozinku (min. 6 karaktera)"
            />
          </div>

          <div className="form-group">
            <label htmlFor="confirmPassword">Potvrdite lozinku:</label>
            <input
              type="password"
              id="confirmPassword"
              name="confirmPassword"
              value={formData.confirmPassword}
              onChange={handleChange}
              required
              placeholder="Potvrdite lozinku"
            />
          </div>

          {error && <div className="error-message">{error}</div>}
          {success && <div className="success-message">{success}</div>}

          <button type="submit" disabled={loading} className="auth-button">
            {loading ? 'Registriranje...' : 'Registriraj se'}
          </button>
        </form>

        <div className="auth-links">
          <p>
            Već imate račun? <Link to="/login">Prijavite se</Link>
          </p>
        </div>
      </div>
    </div>
  );
};

export default Register;
