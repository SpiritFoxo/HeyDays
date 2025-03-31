import React, { useState, useEffect } from 'react';
import "../css/Friends.css";
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
      setRequests((prev) => prev.filter((friend) => friend.id !== friendId));
      const acceptedFriend = requests.find(friend => friend.id === friendId);
      setFriends((prev) => [...prev, acceptedFriend]);
    });
  };

  const handleDecline = (friendId) => {
    declineFriendRequest(friendId, token).then(() => {
      setRequests((prev) => prev.filter((friend) => friend.id !== friendId));
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
              friends.map(({ id, name, surname, profile_photo }) => (
                <div className="friend-item" key={id}>
                  <div className="profile-image">
                    <img className="avatar" src={profile_photo || "/fallback.png"} alt={name} />
                  </div>
                  <div className="friend-info">
                    <span className="friend-name">{name} {surname}</span>
                  </div>
                </div>
              ))
            ) : <p>У вас пока нет друзей</p>}
          </div>
        )}

        {activeTab === 'requests' && (
          <div className="requests-list">
            {requests.length > 0 ? (
              requests.map(({ id, name, surname, profile_photo }) => (
                <div className="friend-item" key={id}>
                  <div className="profile-image">
                    <img className="avatar" src={profile_photo || "/fallback.png"} alt={name} />
                  </div>
                  <div className="friend-info">
                    <span className="friend-name">{name} {surname}</span>
                  </div>
                  <button onClick={() => handleAccept(id)}>Принять</button>
                  <button onClick={() => handleDecline(id)}>Отклонить</button>
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
