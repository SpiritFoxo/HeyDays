import './App.css';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import RegistrationForm from './pages/RegistrationPage';
import MainPage from './pages/MainPage';
import LoginForm from './pages/LoginPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path='/register' element={<RegistrationForm></RegistrationForm>}></Route>
        <Route path='/*' element={<MainPage></MainPage>}></Route>
        <Route path='/login' element={<LoginForm></LoginForm>}></Route>
      </Routes>
    </Router>
  );
}

export default App;
