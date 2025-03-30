import React, { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import './Chat.css';

const Chat = () => {
  const { chatId } = useParams();
  const navigate = useNavigate();
  const [messages, setMessages] = useState([]);
  const [chatInfo, setChatInfo] = useState(null);
  const [messageText, setMessageText] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const messagesEndRef = useRef(null);
  const token = localStorage.getItem("token");
  const currentUserId = parseInt(localStorage.getItem("userId"), 10);
  const wsUrl = `ws://localhost:8080/ws?token=${token}`;

  useEffect(() => {
    const fetchChatInfo = async () => {
      try {
        const response = await fetch(`http://localhost:8080/chat/${chatId}`, {
          method: "GET",
          headers: { Authorization: `Bearer ${token}` },
        });
        if (!response.ok) throw new Error("Failed to fetch chat information");
        const data = await response.json();
        setChatInfo(data);
      } catch (err) {
        setError(err.message);
      }
    };
    fetchChatInfo();
  }, [chatId, token]);

  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const response = await fetch(`http://localhost:8080/chat/${chatId}/messages`, {
          method: "GET",
          headers: { Authorization: `Bearer ${token}` },
        });
        if (!response.ok) throw new Error("Failed to fetch messages");
        const data = await response.json();
        setMessages((data.messages || []).reverse());
        setLoading(false);
      } catch (err) {
        setError(err.message);
        setLoading(false);
      }
    };
    fetchMessages();
  }, [chatId, token]);

  useEffect(() => {
    const ws = new WebSocket(wsUrl);
    ws.onmessage = (event) => {
      try {
        const newMessage = JSON.parse(event.data);
        if (newMessage.ChatID === parseInt(chatId)) {
          setMessages((prevMessages) => [...prevMessages, newMessage]);
        }
      } catch (err) {
        console.error("WebSocket message error:", err);
      }
    };
    ws.onclose = () => console.log("WebSocket closed");
    return () => ws.close();
  }, [chatId, wsUrl]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const handleSendMessage = async (e) => {
    e.preventDefault();
    if (!messageText.trim()) return;
    try {
      const response = await fetch(`http://localhost:8080/chat/send`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
        body: JSON.stringify({ chat_id: parseInt(chatId), content: messageText }),
      });
      if (!response.ok) throw new Error("Failed to send message");
      setMessageText('');
    } catch (err) {
      setError(err.message);
    }
  };

  const handleBackClick = () => navigate('/chats');

  if (loading) return <div className="loading">Загрузка сообщений...</div>;
  if (error) return <div className="error">Ошибка: {error}</div>;

  return (
    <div className="chat-room-container">
      <div className="chat-header">
        <button className="back-button" onClick={handleBackClick}>←</button>
        <div className="chat-info">
          <div className="chat-title">{chatInfo?.title || "Чат"}</div>
          <div className="chat-participants">{chatInfo?.participants.length || 0} участников</div>
        </div>
      </div>
      <div className="messages-container">
        {messages.length > 0 ? (
          <div className="messages-list">
            {messages.map((message) => (
              <div key={message.ID} className={`message ${message.SenderID === currentUserId ? 'sent' : 'received'}`}>
                <div className="message-content">{message.Content}</div>
              </div>
            ))}
            <div ref={messagesEndRef} />
          </div>
        ) : (
          <p>No messages</p>
        )}
      </div>
      <form className="message-form" onSubmit={handleSendMessage}>
        <input
          type="text"
          value={messageText}
          onChange={(e) => setMessageText(e.target.value)}
          placeholder="Введите сообщение..."
          className="message-input"
        />
        <button type="submit" className="send-button">Отправить</button>
      </form>
    </div>
  );
};

export default Chat;
