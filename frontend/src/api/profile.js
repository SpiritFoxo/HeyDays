const API_URL = "http://localhost:8080";
export const fetchProfileData = async (userId, token) => {
    try {
      const url = userId
        ? `${API_URL}/user/profile/${userId}`
        : `${API_URL}/user/profile`;
  
      const response = await fetch(url, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          ...(token && { Authorization: `Bearer ${token}` }),
        },
      });
  
      if (!response.ok) throw new Error("Failed to fetch profile data");
      return await response.json();
    } catch (err) {
      throw new Error(err.message);
    }
  };
  
  export const fetchCurrentUser = async (token) => {
    try {
      const response = await fetch(`${API_URL}/user/profile`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });
  
      if (!response.ok) throw new Error("Failed to fetch current user");
      return await response.json();
    } catch (err) {
      throw new Error(err.message);
    }
  };

  export const openChat = async (token, participantId) => {
    try{
      const response = await fetch(`${API_URL}/chat/direct`,{
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({"user_id": parseInt(participantId)})
      });

      if(!response.ok) throw new Error("Failed to open chat");
      return await response.json();
    } catch (err){
      throw new Error(err.message);
    }
  };
  