import React from 'react'
import './NavBar.css'
import bcrLogo from '../../assets/bcrLogo.png'
import sun from '../../assets/dayLogo.png'
import moon from '../../assets/nightLogo.png'
import magGlassDark from '../../assets/magGlassD.png'
import magGlassLight from '../../assets/magGlassL.png'
import profileIconDark from '../../assets/profileIconD.png'
import profileIconLight from '../../assets/profileIconL.png'

/* for navigation */
import { Link } from 'react-router-dom';


const NavBar = ({theme, setTheme}) => {

  /* change theme logic */
  const toggle_mode = ()=>{
    theme == 'light' ? setTheme('dark') : setTheme('light');
  }

  return (
    <div className='navbar'>
        <img src={bcrLogo} alt="" className='logo'/>
        <ul>
            <li>
              <NavLink to="/" className={({ isActive }) => isActive ? 'active-link' : undefined}>Home</NavLink>
            </li>
            <li>Transactions</li>
            <li>Personalized Advice</li>
            <li><Link to="/profile">User Profile</Link></li>
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
