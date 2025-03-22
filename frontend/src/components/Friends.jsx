import React, { useState, useEffect } from 'react';
import "./Friends.css";
import { getFriends, getFriendRequests, acceptFriendRequest, declineFriendRequest } from '../api/friends';

const Friends = () => {
  const [activeTab, setActiveTab] = useState('friends');
  const [friends, setFriends] = useState([]);
  const [requests, setRequests] = useState([]);
  const token = localStorage.getItem("token"); 

  useEffect(() => {
    if (activeTab === "friends") {
      getFriends(token).then((data) => setFriends(data.friends || []));
    } else {
      getFriendRequests(token).then((data) => setRequests(data.friends || []));
    }
  }, [activeTab]);

  const handleAccept = (friendId) => {
    acceptFriendRequest(friendId, token).then(() => {
      setRequests((prev) => prev.filter((id) => id !== friendId));
      setFriends((prev) => [...prev, friendId]);
    });
  };

  const handleDecline = (friendId) => {
    declineFriendRequest(friendId, token).then(() => {
      setRequests((prev) => prev.filter((id) => id !== friendId));
    });
  };

  return (
    <div className="friends-container">
      <div className="navigation-tabs">
        <button className={`tab-button ${activeTab === 'friends' ? 'active' : ''}`} onClick={() => setActiveTab('friends')}>
          Друзья
        </button>
        <button className={`tab-button ${activeTab === 'requests' ? 'active' : ''}`} onClick={() => setActiveTab('requests')}>
          Запросы в друзья
        </button>
      </div>

      <div className="friends-content">
        {activeTab === 'friends' && (
          <div className="friends-list">
            {friends.length > 0 ? (
              friends.map((friendId) => (
                <div className="friend-item" key={friendId}>
                  <div className="profile-image"><div className="avatar"></div></div>
                  <div className="friend-name">User ID: {friendId}</div>
                </div>
              ))
            ) : <p>У вас пока нет друзей</p>}
          </div>
        )}

        {activeTab === 'requests' && (
          <div className="requests-list">
            {requests.length > 0 ? (
              requests.map((friendId) => (
                <div className="friend-item" key={friendId}>
                  <div className="profile-image"><div className="avatar"></div></div>
                  <div className="friend-name">User ID: {friendId}</div>
                  <button onClick={() => handleAccept(friendId)}>Принять</button>
                  <button onClick={() => handleDecline(friendId)}>Отклонить</button>
                </div>
              ))
            ) : <p>Нет входящих заявок</p>}
          </div>
        )}
      </div>
    </div>
  );
};

export default Friends;