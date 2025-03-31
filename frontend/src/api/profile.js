
export const fetchProfileData = async (userId, token) => {
    try {
      const url = userId
        ? `http://localhost:8080/openapi/profile/${userId}`
        : "http://localhost:8080/user/profile";
  
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
      const response = await fetch("http://localhost:8080/user/profile", {
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
  