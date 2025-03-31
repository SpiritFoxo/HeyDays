import React, { useState } from 'react';
import '../css/RegistrationPage.css';

const RegistrationForm = () => {
  const [formData, setFormData] = useState({
    email: '',
    name: '',
    surname: '',
    password: '',
    confirmPassword: ''
  });

  const [errors, setErrors] = useState({});
  const [isLoading, setIsLoading] = useState(false);
  const [registerSuccess, setRegisterSuccess] = useState(false);

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

    if (!formData.name) {
      newErrors.name = 'Имя обязательно';
    }
    
    if (!formData.surname) {
      newErrors.surname = 'Фамилия обязательна';
    }
    
    if (!formData.password) {
      newErrors.password = 'Пароль обязателен';
    } else if (formData.password.length < 6) {
      newErrors.password = 'Пароль должен содержать минимум 6 символов';
    }

    if (formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = 'Пароли не совпадают';
    }
    
    return newErrors;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    const validationErrors = validate();
    
    if (Object.keys(validationErrors).length === 0) {
      console.log('Form submitted successfully', formData);
      setIsLoading(true);
      
      try {
        const apiData = {
          email: formData.email,
          name: formData.name,
          surname: formData.surname,
          password: formData.password
        };
        
        const response = await fetch('http://localhost:8080/auth/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(apiData),
        });
        
        if (!response.ok) {
          const errorData = await response.json();
          throw new Error(errorData.message || 'Ошибка при регистрации');
        }
        const data = await response.json();
        console.log('Registration successful', data);

        if (data.token) {
          localStorage.setItem('token', data.token);
        }
        
        setRegisterSuccess(true);
        window.location.href = '/login';
      } catch (error) {
        console.error('Registration error:', error);
        
        if (error.message.includes('email')) {
          setErrors({
            ...validationErrors,
            email: 'Email уже используется',
            form: 'Ошибка при регистрации'
          });
        } else {
          setErrors({
            ...validationErrors,
            form: 'Ошибка при регистрации: ' + error.message
          });
        }
      } finally {
        setIsLoading(false);
      }
    } else {
      setErrors(validationErrors);
    }

  };

  return (
    <div className="registration-container">
      <div className="registration-form-container">
        <h2 className="registration-title">Регистрация</h2>
        
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
              type="text"
              name="name"
              placeholder="Имя"
              value={formData.name}
              onChange={handleChange}
              className={errors.name ? 'input-error' : ''}
            />
            {errors.name && <span className="error-message">{errors.name}</span>}
          </div>
          
          <div className="form-group">
            <input
              type="text"
              name="surname"
              placeholder="Фамилия"
              value={formData.surname}
              onChange={handleChange}
              className={errors.surname ? 'input-error' : ''}
            />
            {errors.surname && <span className="error-message">{errors.surname}</span>}
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
          
          <div className="form-group">
            <input
              type="password"
              name="confirmPassword"
              placeholder="Повторите пароль"
              value={formData.confirmPassword}
              onChange={handleChange}
              className={errors.confirmPassword ? 'input-error' : ''}
            />
            {errors.confirmPassword && <span className="error-message">{errors.confirmPassword}</span>}
          </div>
          
          <div className="login-link">
            <a href="/login">Уже есть на нашей платформе?</a>
          </div>
          
          <button type="submit" className="submit-button">Отправить</button>
        </form>
      </div>
      
      <div className="welcome-text">
        <h1>Добро пожаловать</h1>
      </div>
    </div>
  );
};

export default RegistrationForm;