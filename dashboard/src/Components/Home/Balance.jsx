import React, { useState } from 'react';
import moneyBagDay from '../../assets/home_page/moneyBagDay.png'
import moneyBagNight from '../../assets/home_page/moneyBagNight.png'
import eyeDay from '../../assets/home_page/eyeDay.png'
import eyeNight from '../../assets/home_page/eyeNight.png'
import hiddenEyeDay from '../../assets/home_page/hiddenEyeDay.png'
import hiddenEyeNight from '../../assets/home_page/hiddenEyeNight.png'

import './Balance.css'

const Balance = ({ theme, setTheme }) => {

    const toggle_mode = () => {
        theme === 'light' ? setTheme('dark') : setTheme('light');
    };

    const formatBalance = (amount) => {
        const [int, decimal] = amount
        .toFixed(2)
        .toString()
        .replace('.', ',')
        .replace(/\B(?=(\d{3})+(?!\d))/g, '.')
        .split(',');
        return { intPart: int, decimalPart: decimal };
    };

    const { intPart, decimalPart } = formatBalance(2578.31);

    const [isVisible, setIsVisible] = useState(true);
  
    return (
      <div className="account-balance">

        <p className="account-balance-title">Fonduri curente</p>

        <div className={`account-balance-wrapper ${!isVisible ? 'blurred' : ''}`}>
          <span className="account-balance-funds">{intPart},</span>
          <span className="account-balance-decimals">{decimalPart}</span>
          <span className="account-balance-currency">RON</span>
        </div>

        <img
          src={moneyBagDay}
          alt=""
          className="money-bag"
        />

        <img
          src={isVisible ? hiddenEyeDay : eyeDay}
          alt="toggle visibility"
          className="eye"
          onClick={() => setIsVisible(!isVisible)}
        />
      </div>
    );
  };

export default Balance
