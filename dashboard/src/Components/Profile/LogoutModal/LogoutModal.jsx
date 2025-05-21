import React from 'react';
import './LogoutModal.css';

const LogoutModal = ({ onClose, onConfirm }) => {
  return (
    <div className="logout-modal-overlay" onClick={onClose}>
      <div className="logout-modal-container" onClick={(e) => e.stopPropagation()}>
        <h2 className="logout-modal-title">Confirm Logout</h2>
        <p className="logout-modal-body">Are you sure you want to log out?</p>
        <div className="logout-modal-buttons">
          <button className="logout-modal-confirm" onClick={onConfirm}>Yes, log out</button>
          <button className="logout-modal-cancel" onClick={onClose}>Cancel</button>
        </div>
      </div>
    </div>
  );
};

export default LogoutModal;