import React, { useState, useEffect } from 'react';
import { Routes, Route, Link, useLocation } from 'react-router-dom';
import './MainPage.css';
import Profile from "../components/Profile"
import Feed from '../components/Feed';
import Messages from '../components/Messages';
import Friends from '../components/Friends';
import Communities from '../components/Communities';
import Chat from '../components/Chat';

const MainPage = () => {
  const location = useLocation();
  const [activeItem, setActiveItem] = useState('feed');
  useEffect(() => {
    const path = location.pathname.split('/')[1] || 'feed';
    setActiveItem(path);
  }, [location]);

  return (
    <div className="main-container">
      <header className="app-header">
        <div className="logo">HeyDay</div>
      </header>
     
      <div className="sidebar">
        <nav className="nav-menu">
          <ul>
            <li className={activeItem === 'profile' ? 'active' : ''}>
              <Link to="/profile" onClick={() => setActiveItem('profile')}>
                <i className="icon profile-icon"><img src='/profile.png' className='icon' alt="Profile" /></i>
                <span>Профиль</span>
              </Link>
            </li>
            <li className={activeItem === 'feed' ? 'active' : ''}>
              <Link to="/feed" onClick={() => setActiveItem('feed')}>
                <i className="icon feed-icon"><img src='/news.png' className='icon' alt="Feed" /></i>
                <span>Лента</span>
              </Link>
            </li>
            <li className={activeItem === 'messages' ? 'active' : ''}>
              <Link to="/messages" onClick={() => setActiveItem('messages')}>
                <i className="icon messages-icon"><img src='/messages.png' className='icon' alt="Messages" /></i>
                <span>Сообщения</span>
              </Link>
            </li>
            <li className={activeItem === 'friends' ? 'active' : ''}>
              <Link to="/friends" onClick={() => setActiveItem('friends')}>
                <i className="icon friends-icon"><img src='/friends.png' className='icon' alt="Friends" /></i>
                <span>Друзья</span>
              </Link>
            </li>
            <li className={activeItem === 'communities' ? 'active' : ''}>
              <Link to="/communities" onClick={() => setActiveItem('communities')}>
                <i className="icon communities-icon"><img src='/communities.png' className='icon' alt="Communities" /></i>
                <span>Сообщества</span>
              </Link>
            </li>
          </ul>
        </nav>
      </div>
     
      <main className="content-area">
        <Routes>
          <Route path="/profile" element={<Profile />} />
          <Route path="/profile/:userId" element={<Profile />} />
          <Route path="/feed" element={<Feed />} />
          <Route path="/messages" element={<Messages />} />
          <Route path="/chat/:chatId" element={<Chat />} />
          <Route path="/friends" element={<Friends />} />
          <Route path="/communities" element={<Communities />} />
          <Route index element={<Feed />} />
        </Routes>
      </main>
    </div>
  );
};

export default MainPage;