import React from 'react'
import './NavBar.css'
import bcrLogo from '../../assets/bcrLogo.png'
import sun from '../../assets/dayLogo.png'
import moon from '../../assets/nightLogo.png'
import magGlassDark from '../../assets/magGlassD.png'
import magGlassLight from '../../assets/magGlassL.png'

/* for navigation */
import { NavLink, Link } from 'react-router-dom';


const NavBar = ({theme, setTheme}) => {

  /* change theme logic */
  const toggle_mode = ()=>{
    theme == 'light' ? setTheme('dark') : setTheme('light');
  }

  return (
    <div className='navbar'>
        <img src={bcrLogo} alt="" className='logo'/>
        <ul>
            <li><NavLink to="/home" className={({ isActive }) => isActive ? 'active-link' : ''}>Home</NavLink></li>
            <li><NavLink to="/discover" className={({ isActive }) => isActive ? 'active-link' : ''}>Discover</NavLink></li>
            <li><NavLink to="/products" className={({ isActive }) => isActive ? 'active-link' : ''}>Products</NavLink></li>
            <li><NavLink to="/profile" className={({ isActive }) => isActive ? 'active-link' : ''}>Profile</NavLink></li>
            <li><NavLink to="/settings" className={({ isActive }) => isActive ? 'active-link' : ''}>Settings</NavLink></li>
        </ul>

        <div className='search-box'>
            <input type="text" placeholder='Search'/>
            <img src={theme =='light' ? magGlassLight : magGlassDark} alt=""/>
        </div>

        <img onClick={()=>{toggle_mode()}} src={theme =='light' ? moon : sun}
        alt="" className='toggle-icon'/>

    </div>
  )
}

export default NavBar
