import moneyBagDay from '../../../assets/home_page/moneyBagDay.png'
import moneyBagNight from '../../../assets/home_page/moneyBagNight.png'

import eyeDay from '../../../assets/home_page/eyeDay.png'
import eyeNight from '../../../assets/home_page/eyeNight.png'

import hiddenEyeDay from '../../../assets/home_page/hiddenEyeDay.png'
import hiddenEyeNight from '../../../assets/home_page/hiddenEyeNight.png'

import React, { useState } from 'react';

import './Balance.css'

const Balance = ({ theme, lang }) => {

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

    const [isVisible, setIsVisible] = useState(false);
  
    return (
      <div className={`account-balance ${theme === 'dark' ? 'dark' : ''}`}>

        <p className={`account-balance-title ${lang === 'RO' ? 'title-ro' : 'title-en'} ${theme === 'dark' ? 'dark' : ''}`}>
          {lang === 'RO' ? 'Fonduri curente' : 'Current funds'}
        </p>

        <div className={`account-balance-wrapper ${!isVisible ? 'blurred' : ''} ${theme === 'dark' ? 'dark' : ''}`}>
          <span className="account-balance-funds">{intPart},</span>
          <span className="account-balance-decimals">{decimalPart}</span>
          <span className="account-balance-currency">RON</span>
        </div>

        <img
          src={(theme === 'dark' ? moneyBagNight : moneyBagDay)}
          alt=""
          className="money-bag"
        />

        <img
          src={isVisible ? (theme === 'dark' ? hiddenEyeNight : hiddenEyeDay) : (theme === 'dark' ? eyeNight : eyeDay)}
          alt="toggle visibility"
          className="eye"
          onClick={() => setIsVisible(!isVisible)}
        />
      </div>
    );
  };

export default Balance