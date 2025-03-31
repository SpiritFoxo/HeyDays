import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { fetchProfileData, fetchCurrentUser, openChat } from '../api/profile';
import { sendFriendRequest } from '../api/friends';
import '../css/Profile.css';

const Profile = () => {
  const [profileData, setProfileData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isOwnProfile, setIsOwnProfile] = useState(false);
  const [requestLoading, setRequestLoading] = useState(false);
  const { userId } = useParams();
  const token = localStorage.getItem("token");

  useEffect(() => {
    const loadProfile = async () => {
      try {
        setLoading(true);
        const profile = await fetchProfileData(userId, token);
        setProfileData(profile);
 
        if (userId) {
          const currentUser = await fetchCurrentUser(token);
          const isOwn = currentUser.id === Number(userId);
          setIsOwnProfile(isOwn);
        } else {
          setIsOwnProfile(true);
        }
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
 
    loadProfile();
  }, [userId, token]);
 
  const handleSendFriendRequest = async () => {
    if (!userId || 
        profileData.friendship_status === 'pending' || 
        profileData.friendship_status === 'accepted' || 
        requestLoading) {
      return;
    }
 
    try {
      setRequestLoading(true);
      const response = await sendFriendRequest(Number(userId), token);
 
      if (response.success) {
        const updatedProfile = await fetchProfileData(userId, token);
        setProfileData(updatedProfile);
      } else {
        console.error("Failed to send friend request:", response.message);
      }
    } catch (error) {
      console.error("Error sending friend request:", error);
    } finally {
      setRequestLoading(false);
    }
  };

  const navigate = useNavigate();
  const handleOpenChat = async () => {
    try {
      const response = await openChat(token, userId);
      
      if (response.chat_id) {
        navigate(`/chat/${response.chat_id}`);
      }
    } catch (error) {
      console.error("Ошибка при открытии чата:", error.message);
    }
  };
  
  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!profileData) return <div>No profile data found</div>;
  
  const profilePhoto = profileData?.profile_photo || "/fallback.png";
  const name = profileData?.name || "";
  const surname = profileData?.surname || "";
  
  return (
    <div className="profile-container">
      <div className="parent">
        <div className='profile-banner'>
          <img src='/banner.png' alt="Profile banner" />
        </div>
        <div className="main-info">
          <div className='main-info-container'>
            <img src={profilePhoto} alt="Profile" />
            <p>{name} {surname}</p>
          </div>
        </div>
        <div className="social-menu">
          <div className='menu-container'>
            <div className='social-buttons'>
            {isOwnProfile ? (
              <button className='social-button'>Редактировать профиль</button>
            ) : (
              <>
                <button 
                  className='social-button'
                  onClick={handleSendFriendRequest}
                  disabled={profileData.friendship_status === 'pending' || profileData.friendship_status === 'accepted' || requestLoading}
                >
                  {requestLoading
                    ? 'Отправка...'
                    : profileData.friendship_status === 'pending'
                    ? 'Запрос отправлен'
                    : profileData.friendship_status === 'accepted'
                    ? 'Вы друзья'
                    : 'Добавить в друзья'}
                </button>
                <button className='social-button' onClick={handleOpenChat}>Сообщение</button>
              </>
            )}
            </div>
            <div className='photo-gallery'></div>
          </div>
        </div>
        <div className="user-posts">
          <div className='posts-container'></div>
        </div>
      </div>
    </div>
  );
};

export default Profile;