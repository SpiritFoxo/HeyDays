import './App.css';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import RegistrationForm from './RegistrationPage';
import MainPage from './MainPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path='/register' element={<RegistrationForm></RegistrationForm>}></Route>
        <Route path='/*' element={<MainPage></MainPage>}></Route>
      </Routes>
    </Router>
  );
}

export default App;
