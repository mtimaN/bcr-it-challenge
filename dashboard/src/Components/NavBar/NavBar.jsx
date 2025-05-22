import React from 'react';
import './NavBar.css';
import bcrLogo from '../../assets/nav_bar/bcrLogo.png';
import sun from '../../assets/nav_bar/dayLogo.png';
import moon from '../../assets/nav_bar/nightLogo.png';
import { NavLink } from 'react-router-dom';
import roD from '../../assets/nav_bar/roDark.png';
import enD from '../../assets/nav_bar/enDark.png';
import roL from '../../assets/nav_bar/roLight.png';
import enL from '../../assets/nav_bar/enLight.png';

const NavBar = ({ theme, setTheme, lang, setLang }) => {
  const toggleMode = () => {
    theme === 'light' ? setTheme('dark') : setTheme('light');
  };

  const toggleLang = () => {
    setLang(lang === 'RO' ? 'EN' : 'RO');
  };

  const getLangIcon = () => {
    if (lang === 'RO') {
      return theme === 'light' ? roD : roL;
    }

    return theme === 'light' ? enD : enL;
  };

  return (
    <div className='navbar'>
      <img src={bcrLogo} alt="BCR Logo" className='logo' />

      <ul>
        <li><NavLink to="/home" className={({ isActive }) => isActive ? 'active-link' : ''}>{lang === 'RO' ? 'Acasă': 'Home'}</NavLink></li>
        <li><NavLink to="/discover" className={({ isActive }) => isActive ? 'active-link' : ''}>{lang === 'RO' ? 'Descoperă': 'Discover'}</NavLink></li>
        <li><NavLink to="/profile" className={({ isActive }) => isActive ? 'active-link' : ''}>{lang === 'RO' ? 'Profilul tău': 'Your profile'}</NavLink></li>
      </ul>

      <div className='right-icons'>
        <img
          onClick={toggleMode}
          src={theme === 'light' ? moon : sun}
          alt="Theme Toggle"
          className='toggle-icon'
        />
        <img
          onClick={toggleLang}
          src={getLangIcon()}
          alt="Language Toggle"
          className='lang-icon toggle-icon'
          style={{ marginLeft: '20px' }}
        />
      </div>
    </div>
  );
};

export default NavBar;
