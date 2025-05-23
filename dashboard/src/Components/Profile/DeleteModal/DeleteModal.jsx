import React from 'react';
import './DeleteModal.css';

const DeleteModal = ({ onClose, onConfirm }) => {
  return (
    <div className="delete-modal-overlay" onClick={onClose}>
      <div className="delete-modal-container" onClick={(e) => e.stopPropagation()}>
        <h2 className="delete-modal-title">Confirmare ștergere?</h2>
        <p className="delete-modal-body">Această acțiune va șterge contul permanent.</p>
        <div className="delete-modal-buttons">
          <button className="delete-modal-confirm" onClick={onConfirm}>Ștergere</button>
          <button className="delete-modal-cancel" onClick={onClose}>Anulare</button>
        </div>
      </div>
    </div>
  );
};

export default DeleteModal;
