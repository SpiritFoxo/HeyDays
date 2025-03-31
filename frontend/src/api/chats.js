const API_URL = "http://localhost:8080/chat";

export const fetchChats = async (token) => {
    try {
      const response = await fetch(`${API_URL}/list`, {
        method: "GET",
        headers: { Authorization: `Bearer ${token}` },
      });
  
      if (!response.ok) {
        throw new Error("Failed to fetch chats");
      }
  
      const data = await response.json();
      return data.chats || [];
    } catch (err) {
      throw new Error(err.message);
    }
  };

  
  export const fetchChatInfo = async (chatId, token) => {
    try {
      const response = await fetch(`${API_URL}/${chatId}`, {
        method: "GET",
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!response.ok) throw new Error("Failed to fetch chat information");
      return await response.json();
    } catch (err) {
      throw new Error(err.message);
    }
  };
  
  export const fetchMessages = async (chatId, token) => {
    try {
      const response = await fetch(`${API_URL}/${chatId}/messages`, {
        method: "GET",
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!response.ok) throw new Error("Failed to fetch messages");
      const data = await response.json();
      return (data.messages || []).reverse();
    } catch (err) {
      throw new Error(err.message);
    }
  };
  
  export const sendMessage = async (chatId, messageText, token) => {
    try {
      const response = await fetch(`${API_URL}/send`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
        body: JSON.stringify({ chat_id: parseInt(chatId), content: messageText }),
      });
      if (!response.ok) throw new Error("Failed to send message");
    } catch (err) {
      throw new Error(err.message);
    }
  };
  