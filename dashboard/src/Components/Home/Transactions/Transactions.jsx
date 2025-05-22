import React from 'react';

import './Transactions.css'

const Transactions = ({ theme, lang }) => {

    const formatBalance = (amount) => {
        return amount
          .toFixed(2)
          .toString()
          .replace('.', ',')
          .replace(/\B(?=(\d{3})+(?!\d))/g, '.');
    };

    const transactionData = [
        {
          title: 'SC OPTIMUS DIGITAL',
          date: '22/05/2025',
          description_ro: 'Plată cu cardul pe 22 mai la 03:03.',
          description_eng: 'Payment with card on May 22, 03:03.',
          amount: -1578.31,
        },
        {
          title: 'ING BANK',
          date: '21/05/2025',
          description_ro: 'Depunere numerar.',
          description_eng: 'Cash deposit.',
          amount: 3250.00,
        },
        {
          title: 'EMAG',
          date: '20/05/2025',
          description_ro: 'Achiziție produs online.',
          description_eng: 'Online product purchase.',
          amount: -280.49,
        },
        {
          title: 'SC OPTIMUS DIGITAL',
          date: '22/05/2025',
          description_ro: 'Plată cu cardul pe 22 mai la 03:03.',
          description_eng: 'Payment with card on May 22, 03:03.',
          amount: -1578.31,
        },
        {
          title: 'ING BANK',
          date: '21/05/2025',
          description_ro: 'Depunere numerar.',
          description_eng: 'Cash deposit.',
          amount: 3250.00,
        },
        {
          title: 'EMAG',
          date: '20/05/2025',
          description_ro: 'Achiziție produs online.',
          description_eng: 'Online product purchase.',
          amount: -280.49,
        },
        {
          title: 'SC OPTIMUS DIGITAL',
          date: '22/05/2025',
          description_ro: 'Plată cu cardul pe 22 mai la 03:03.',
          description_eng: 'Payment with card on May 22, 03:03.',
          amount: -1578.31,
        },
        {
          title: 'ING BANK',
          date: '21/05/2025',
          description_ro: 'Depunere numerar.',
          description_eng: 'Cash deposit.',
          amount: 3250.00,
        },
        {
          title: 'EMAG',
          date: '20/05/2025',
          description_ro: 'Achiziție produs online.',
          description_eng: 'Online product purchase.',
          amount: -280.49,
        },
    ];
  
    return (
      <div className={`account-transactions ${theme === 'dark' ? 'dark' : ''}`}>
        <div className="transaction-scroll-area">
            {transactionData.slice(0, 10).map((tx, index) => (
            <div key={index} className="transaction-wrapper">
                <div className="transaction-name">
                <span className="transaction-title">{tx.title}</span>
                <span className="transaction-date">{tx.date}</span>
                <span className="transaction-description">{lang === 'RO' ? tx.description_ro : tx.description_eng}</span>
                </div>

                <span className={`transaction-funds ${tx.amount < 0 ? 'negative' : 'positive'}`}>
                {formatBalance(tx.amount)}
                </span>
                <span className={`transaction-currency ${tx.amount < 0 ? 'negative' : 'positive'}`}>
                RON
                </span>
            </div>
            ))}
        </div>
      </div>
    );
  };

export default Transactions