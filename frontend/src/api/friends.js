const API_URL = "http://localhost:8080/friends";

export const getFriends = async (token) => {
  const response = await fetch(`${API_URL}/get-friends`, {
    method: "GET",
    headers: { Authorization: `Bearer ${token}` },
  });
  return response.json();
};

export const getFriendRequests = async (token) => {
  const response = await fetch(`${API_URL}/get-pending-friend-requests`, {
    method: "GET",
    headers: { Authorization: `Bearer ${token}` },
  });
  return response.json();
};

export const sendFriendRequest = async (friendId, token) => {
  const response = await fetch(`${API_URL}/send-friend-request`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ friend_id: friendId }),
  });
  return response.json();
};

export const acceptFriendRequest = async (friendId, token) => {
  const response = await fetch(`${API_URL}/accept-friend-request`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ friend_id: friendId }),
  });
  return response.json();
};

export const declineFriendRequest = async (friendId, token) => {
  const response = await fetch(`${API_URL}/decline-friend-request`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ friend_id: friendId }),
  });
  return response.json();
};
