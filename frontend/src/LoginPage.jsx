import React, { useState } from 'react';
import './LoginPage.css';

const LoginForm = () => {
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  });

  const [errors, setErrors] = useState({});
  const [isLoading, setIsLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
  };

  const validate = () => {
    const newErrors = {};
    
    if (!formData.email) {
      newErrors.email = 'Email обязателен';
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      newErrors.email = 'Неверный формат email';
    }
    
    if (!formData.password) {
      newErrors.password = 'Пароль обязателен';
    }
    
    return newErrors;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    const validationErrors = validate();
    
    if (Object.keys(validationErrors).length === 0) {
        setIsLoading(true);
        try {
          const response = await fetch('http://localhost:8080/api/login', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData),
          });
          
          if (!response.ok) {
            throw new Error('Login failed');
          }
          
          const data = await response.json();

          localStorage.setItem('token', data.token);

          window.location.href = '/';
        } catch (error) {
          console.error('Login error:', error);
          setErrors({ 
            form: 'Неверный email или пароль' 
          });
        } finally {
          setIsLoading(false);
        }
      } else {
        setErrors(validationErrors);
      }
  };

  return (
    <div className="login-container">
      <div className="login-form-container">
        <h2 className="login-title">Вход</h2>
        
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <input
              type="email"
              name="email"
              placeholder="E-mail"
              value={formData.email}
              onChange={handleChange}
              className={errors.email ? 'input-error' : ''}
            />
            {errors.email && <span className="error-message">{errors.email}</span>}
          </div>
          
          <div className="form-group">
            <input
              type="password"
              name="password"
              placeholder="Пароль"
              value={formData.password}
              onChange={handleChange}
              className={errors.password ? 'input-error' : ''}
            />
            {errors.password && <span className="error-message">{errors.password}</span>}
          </div>
          
          <div className="register-link">
            <a href="/register">Нет аккаунта? Зарегистрироваться</a>
          </div>
          
          <button type="submit" className="submit-button">Войти</button>
        </form>
      </div>
      
      <div className="welcome-text">
        <h1>С возвращением!</h1>
      </div>
    </div>
  );
};

export default LoginForm;