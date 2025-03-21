import React, {useState, useEffect} from 'react';
import { useParams} from 'react-router-dom';
import './Profile.css';

const Profile = () => {

  const [profileData, setProfileData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const { userId } = useParams();

  useEffect(() => {
    const fetchProfileData = async () => {
      try {
        setLoading(true);
        let response;
        
        if (userId) {
          response = await fetch(`http://localhost:8080/api/profile/${userId}`, {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
            }
          });
        } else {
          response = await fetch('http://localhost:8080/auth/profile', {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
          });
        }

        if (!response.ok) {
          throw new Error('Failed to fetch profile data');
        }

        const data = await response.json();
        setProfileData(data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchProfileData();
  }, [userId]);

  if (error) return <div>Error: {error}</div>;
  if (!profileData) return <div>No profile data found</div>;

  const profilePhoto = profileData?.profile_photo || '/pfp.png';
  const name = profileData?.name || '';
  const surname = profileData?.surname || '';

  return (
    <div className="profile-container">

      <div className="parent">
        <div className='profile-banner'>
          <img src='/banner.png'></img>
        </div>
        <div className="main-info">
          <div className='main-info-container'>
            <img src={profilePhoto}></img>
            <p>{name} {surname}</p>
          </div>
        </div>
        <div className="social-menu">
          <div className='menu-container'>
            <div className='social-buttons'>
              <button className='social-button'>Запрос в друзья</button>
              <button className='social-button'>Сообщение</button>
            </div>

            <div className='photo-gallery'></div>
          </div>
        </div>
        <div className="user-posts">
          <div className='posts-container'>

          </div>
        </div>
      </div>
    
    </div>
  );
};

export default Profile;