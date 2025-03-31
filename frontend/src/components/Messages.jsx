import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import '../css/Messages.css';
import { fetchChats } from '../api/chats';

const Messages = () => {
  const [chats, setChats] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const token = localStorage.getItem("token");
  const wsUrl = `ws://localhost:8080/ws?token=${token}`;

  useEffect(() => {
    const loadChats = async () => {
      try {
        const chatsData = await fetchChats(token);
        setChats(chatsData);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    loadChats();
    const ws = new WebSocket(wsUrl);
    ws.onmessage = (event) => {
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
            src={chat?.participants?.[0]?.photo_url || "/fallback.png"} 
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
