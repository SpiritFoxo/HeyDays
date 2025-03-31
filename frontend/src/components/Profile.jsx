import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { fetchProfileData, fetchCurrentUser } from '../api/profile';
import { sendFriendRequest } from '../api/friends';
import './Profile.css';

const Profile = () => {
  const [profileData, setProfileData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isOwnProfile, setIsOwnProfile] = useState(false);
  const [requestSent, setRequestSent] = useState(false);
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
          setIsOwnProfile(currentUser.id === Number(userId));
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
    if (!userId || requestSent || requestLoading) return;

    try {
      setRequestLoading(true);
      const response = await sendFriendRequest(Number(userId), token);

      if (response.success) {
        setRequestSent(true);
      } else {
        console.error("Failed to send friend request:", response.message);
      }
    } catch (error) {
      console.error("Error sending friend request:", error);
    } finally {
      setRequestLoading(false);
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!profileData) return <div>No profile data found</div>;

  const profilePhoto = profileData?.profile_photo || "/pfp.png";
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
                    disabled={requestSent || requestLoading}
                  >
                    {requestLoading ? 'Отправка...' : 
                     requestSent ? 'Запрос отправлен' : 'Запрос в друзья'}
                  </button>
                  <button className='social-button'>Сообщение</button>
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
