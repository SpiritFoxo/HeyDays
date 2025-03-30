import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import './Messages.css';

const Messages = () => {
  const [chats, setChats] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const token = localStorage.getItem("token");
  const wsUrl = `ws://localhost:8080/ws?token=${token}`;

  useEffect(() => {
    const fetchChats = async () => {
      try {
        const response = await fetch("http://localhost:8080/chat/list", {
          method: "GET",
          headers: { Authorization: `Bearer ${token}` },
        });

        if (!response.ok) {
          throw new Error("Failed to fetch chats");
        }

        const data = await response.json();
        setChats(data.chats || []);
        setLoading(false);
      } catch (err) {
        setError(err.message);
        setLoading(false);
      }
    };

    fetchChats();
    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      console.log("WebSocket подключен");
    };

    ws.onmessage = (event) => {
      console.log("Новое сообщение в WebSocket:", event.data);
      try {
        const newMessage = JSON.parse(event.data);
    
        setChats((prevChats) => {
          return prevChats.map((chat) => {
            if (chat.id === newMessage.ChatID) {
              return {
                ...chat,
                last_message: newMessage.Content,
                last_sender_name: `Пользователь ${newMessage.SenderID}`,
                last_message_time: newMessage.CreatedAt,
                unread_count: chat.unread_count + 1,
              };
            }
            return chat;
          });
        });
      } catch (err) {
        console.error("Ошибка при обработке WebSocket-сообщения:", err);
      }
    };

    ws.onerror = (err) => {
      console.error("Ошибка WebSocket:", err);
    };

    ws.onclose = () => {
      console.log("WebSocket соединение закрыто");
    };

    return () => {
      ws.close();
    };
  }, [token]);

  if (loading) return <div className="loading">Загрузка чатов...</div>;
  if (error) return <div className="error">Ошибка: {error}</div>;

  return (
    <div className="chat-list-container">
      <h2>Диалоги</h2>

      {chats.length > 0 ? (
  <div className="chats-list">
    {chats.map((chat) => (
      <Link to={`/chat/${chat.id}`} key={chat.id} className="chat-item">
        <div className="chat-avatar">
          <img 
            src={chat?.participants?.[0]?.photo_url || "/default-chat.png"} 
            alt={chat.name || "Чат"} 
          />
        </div>
        <div className="chat-details">
          <div className="chat-title">
            {chat.name || "Без названия"}
          </div>
          <div className="chat-last-message">
            {chat.last_sender_name ? `${chat.last_sender_name}: ` : ""}{chat.last_message || "Нет сообщений"}
          </div>
        </div>
        <div className="chat-meta">
          {chat.last_message_time && (
            <div className="chat-time">
              {new Date(chat.last_message_time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
            </div>
          )}
          {chat.unread_count > 0 && (
            <div className="unread-count">{chat.unread_count}</div>
          )}
        </div>
      </Link>
    ))}
  </div>
      ) : (
        <div className="no-chats">
          <p>У вас пока нет диалогов</p>
          <button className="new-chat-btn">Начать новый диалог</button>
        </div>
      )}
    </div>
  );
};

export default Messages;
