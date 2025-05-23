import React from 'react';
import './DeleteModal.css';

const DeleteModal = ({ onClose, onConfirm, lang, theme }) => {
  return (
    <div className="delete-modal-overlay" onClick={onClose}>
      <div className="delete-modal-container" onClick={(e) => e.stopPropagation()}>
        <h2 className="delete-modal-title">{lang === 'RO' ? 'Confirmare ștergere?' : 'Confirm delete?'}</h2>
        <p className="delete-modal-body">{lang === 'RO' ? 'Această acțiune va șterge contul permanent.' : 'This action will erase your account permanently.'}</p>
        <div className="delete-modal-buttons">
          <button className="delete-modal-confirm" onClick={onConfirm}>{lang === 'RO' ? 'Ștergere' : 'Delete'}</button>
          <button className="delete-modal-cancel" onClick={onClose}>{lang === 'RO' ? 'Anulare' : 'Cancel'}</button>
        </div>
      </div>
    </div>
  );
};

export default DeleteModal;
