import React from 'react';
import './DeleteModal.css';

const DeleteModal = ({ onClose, onConfirm }) => {
  return (
    <div className="delete-modal-overlay" onClick={onClose}>
      <div className="delete-modal-container" onClick={(e) => e.stopPropagation()}>
        <h2 className="delete-modal-title">Confirm Delete?</h2>
        <p className="delete-modal-body">This action will permanently delete your account.</p>
        <div className="delete-modal-buttons">
          <button className="delete-modal-confirm" onClick={onConfirm}>Delete</button>
          <button className="delete-modal-cancel" onClick={onClose}>Cancel</button>
        </div>
      </div>
    </div>
  );
};

export default DeleteModal;