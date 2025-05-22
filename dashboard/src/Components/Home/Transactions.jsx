import React from 'react';

import './Transactions.css'

const Transactions = ({ theme, setTheme }) => {

    const toggle_mode = () => {
        theme === 'light' ? setTheme('dark') : setTheme('light');
    };

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
          description: 'Plată cu cardul pe 22 mai la 03:03.',
          amount: -1578.31,
        },
        {
          title: 'ING BANK',
          date: '21/05/2025',
          description: 'Depunere numerar.',
          amount: 3250.00,
        },
        {
          title: 'EMAG',
          date: '20/05/2025',
          description: 'Achiziție produs online.',
          amount: -280.49,
        },
        {
            title: 'SC OPTIMUS DIGITAL',
            date: '22/05/2025',
            description: 'Plată cu cardul pe 22 mai la 03:03.',
            amount: -1578.31,
            },
            {
            title: 'ING BANK',
            date: '21/05/2025',
            description: 'Depunere numerar.',
            amount: 3250.00,
            },
            {
            title: 'EMAG',
            date: '20/05/2025',
            description: 'Achiziție produs online.',
            amount: -280.49,
            },
            {
                title: 'SC OPTIMUS DIGITAL',
                date: '22/05/2025',
                description: 'Plată cu cardul pe 22 mai la 03:03.',
                amount: -1578.31,
              },
              {
                title: 'ING BANK',
                date: '21/05/2025',
                description: 'Depunere numerar.',
                amount: 3250.00,
              },
              {
                title: 'EMAG',
                date: '20/05/2025',
                description: 'Achiziție produs online.',
                amount: -280.49,
              },
    ];
  
    return (
        <div className="account-transactions">
        <div className="transaction-scroll-area">
            {transactionData.slice(0, 10).map((tx, index) => (
            <div key={index} className="transaction-wrapper">
                <div className="transaction-name">
                <span className="transaction-title">{tx.title}</span>
                <span className="transaction-date">{tx.date}</span>
                <span className="transaction-description">{tx.description}</span>
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